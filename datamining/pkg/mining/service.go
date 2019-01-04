package mining

import (
	"errors"
	"log"
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"

	emlisting "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

//ErrUnsupportedTransaction when the transaction does not have transaction miners associated
var ErrUnsupportedTransaction = errors.New("Unsupported transaction")

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//TransactionType represents the transaction type
type TransactionType int

const (
	//KeychainTransaction represents transaction related to keychain
	KeychainTransaction TransactionType = 0

	//IDTransaction represents transaction related to ID data
	IDTransaction TransactionType = 1
)

//TransactionMiner define methods a transaction miner must define
type TransactionMiner interface {

	//GetLastTransactionHash returns the last transaction from a given address
	GetLastTransactionHash(addr string) (string, error)

	//CheckAsMaster performs checks on some data like a master peer
	CheckAsMaster(txHash string, data interface{}) error

	//CheckAsSlave performs checks on some data like a peer inside a validation pool
	CheckAsSlave(txHash string, data interface{}) error
}

type signer interface {
	PowSigVerifier
	ValidationSigner
}

//Service defines methods for global mining
type Service interface {

	//LeadMining process workflow to lead mining (like elected master peer)
	//
	//The workflow includes:
	// - Locks the transaction
	// - Checks (as master)
	// - Executes the proof of work
	// - Requests validations (as slave)
	// - Requests storage
	// - Unlocks the transaction
	//
	//It also in charge of notify the transaction status during this workflow
	LeadMining(txHash string, addr string, data interface{}, vPool datamining.Pool, txType TransactionType, emSig string) error

	//Validate performs checks like a peer in a validation pool and create a validation (successed or not)
	Validate(txHash string, data interface{}, txType TransactionType) (Validation, error)
}

type service struct {
	aiClient AIClient
	notif    Notifier
	poolF    PoolFinder
	poolR    PoolRequester
	signer   signer
	emLister emlisting.Service
	config   system.UnirisConfig
	txMiners map[TransactionType]TransactionMiner
}

//NewService creates a new global mining service
func NewService(aiCli AIClient, n Notifier, pF PoolFinder, pR PoolRequester, signer signer, emLister emlisting.Service, config system.UnirisConfig, txMiners map[TransactionType]TransactionMiner) Service {
	return service{aiCli, n, pF, pR, signer, emLister, config, txMiners}
}

func (s service) LeadMining(txHash string, addr string, data interface{}, vPool datamining.Pool, txType TransactionType, emSig string) error {

	if s.txMiners[txType] == nil {
		return s.notif.NotifyTransactionStatus(txHash, TxInvalid)
	}

	if err := s.notif.NotifyTransactionStatus(txHash, TxPending); err != nil {
		return err
	}

	lastVPool, sPool, err := s.findPools(addr)
	if err != nil {
		return err
	}

	if err := s.requestLock(txHash, addr, lastVPool); err != nil {
		return err
	}

	//Process asynchrounously the transaction mining, the validation, the storage and the unlocking
	go func() {
		if err := s.processMining(txHash, data, addr, emSig, lastVPool, vPool, sPool, txType); err != nil {
			log.Printf("Mining error: %s", err.Error())
		}
	}()

	return nil
}

func (s service) processMining(txHash string, data interface{}, addr string, emSig string, lastVPool, vPool, sPool datamining.Pool, txType TransactionType) error {
	endorsement, err := s.mine(txHash, data, addr, emSig, lastVPool, vPool, txType)
	if err != nil {
		if err == ErrInvalidTransaction {
			return s.notif.NotifyTransactionStatus(txHash, TxInvalid)
		}
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, TxApproved); err != nil {
		return err
	}

	minReplicas, err := s.aiClient.GetMininumReplications(txHash)
	if err != nil {
		return err
	}
	if err := s.poolR.RequestStorage(minReplicas, sPool, data, endorsement, txType); err != nil {
		if err := s.notif.NotifyTransactionStatus(txHash, TxInvalid); err != nil {
			return err
		}
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxReplicated); err != nil {
		return err
	}

	return s.requestUnlock(txHash, addr, lastVPool)
}

func (s service) findPools(addr string) (datamining.Pool, datamining.Pool, error) {
	lastVPool, err := s.poolF.FindLastValidationPool(addr)
	if err != nil {
		return nil, nil, err
	}

	sPool, err := s.poolF.FindStoragePool(addr)
	if err != nil {
		return nil, nil, err
	}

	return lastVPool, sPool, nil
}

func (s service) requestLock(txHash string, addr string, lastVPool datamining.Pool) error {
	//Build lock transaction
	lock := lock.TransactionLock{
		TxHash:         txHash,
		MasterRobotKey: s.config.SharedKeys.Robot.PublicKey,
		Address:        addr,
	}

	if err := s.poolR.RequestLock(lastVPool, lock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxLocked); err != nil {
		return err
	}

	return nil
}

func (s service) requestUnlock(txHash string, addr string, lastVPool datamining.Pool) error {
	//Build unlock transaction
	lock := lock.TransactionLock{
		TxHash:         txHash,
		MasterRobotKey: s.config.SharedKeys.Robot.PublicKey,
		Address:        addr,
	}
	if err := s.poolR.RequestUnlock(lastVPool, lock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxUnlocked); err != nil {
		return err
	}

	return nil
}

func (s service) mine(txHash string, data interface{}, addr string, txEmSig string, lastVPool, vPool datamining.Pool, txType TransactionType) (Endorsement, error) {
	//Execute transaction specific master checks
	if err := s.txMiners[txType].CheckAsMaster(txHash, data); err != nil {
		return nil, err
	}

	//Execute the Proof of Work
	pow := pow{
		lastVPool:   lastVPool,
		emLister:    s.emLister,
		robotPubKey: s.config.PublicKey,
		robotPvKey:  s.config.PrivateKey,
		signer:      s.signer,
		txEmSig:     txEmSig,
		txData:      data,
		txType:      txType,
	}
	masterValid, err := pow.execute()
	if err != nil {
		return nil, err
	}

	minValid, err := s.aiClient.GetMininumValidations(txHash)
	if err != nil {
		return nil, err
	}

	//Ask a pool of peers to validate the transaction
	valids, err := s.poolR.RequestValidations(minValid, vPool, txHash, data, txType)
	if err != nil {
		return nil, err
	}

	//Check if the validations passed
	for _, v := range valids {
		if v.Status() == ValidationKO {
			return nil, ErrInvalidTransaction
		}
	}

	lastTxHash, err := s.txMiners[txType].GetLastTransactionHash(addr)
	if err != nil {
		return nil, err
	}
	return NewEndorsement(lastTxHash, txHash, masterValid, valids), nil
}

func (s service) Validate(txHash string, data interface{}, txType TransactionType) (Validation, error) {
	if s.txMiners[txType] == nil {
		return nil, ErrUnsupportedTransaction
	}

	if err := s.txMiners[txType].CheckAsSlave(txHash, data); err != nil {
		if err == ErrInvalidTransaction || err == ErrUnsupportedTransaction {
			return s.buildValidation(ValidationKO)
		}
		return nil, err
	}
	return s.buildValidation(ValidationOK)
}

func (s service) buildValidation(status ValidationStatus) (Validation, error) {
	v := validation{
		pubk:      s.config.PublicKey,
		status:    status,
		timestamp: time.Now(),
	}
	sValid, err := s.signer.SignValidation(v, s.config.PrivateKey)
	if err != nil {
		return nil, err
	}
	return sValid, nil
}
