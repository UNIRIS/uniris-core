package mining

import (
	"errors"
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"

	biodlisting "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//Checker define methods for a checker to implement
type Checker interface {
	CheckAsMaster(txHash string, data interface{}) error
	CheckAsSlave(txHash string, data interface{}) error
}

//Signer defines methods to handle lead mining signing
type Signer interface {
	lock.Signer
	PowSigner
}

//Service defines methods for global mining
type Service interface {
	LeadMining(txHash string, addr string, data interface{}, vPool Pool, txType TransactionType, biodSig string) error
	Validate(txHash string, data interface{}, txType TransactionType) (datamining.Validation, error)
}

type service struct {
	notif      Notifier
	poolF      PoolFinder
	poolR      PoolRequester
	signer     Signer
	biodLister biodlisting.Service
	robotKey   string
	robotPvKey string
	checks     map[TransactionType]Checker
}

//NewService creates a new global mining service
func NewService(n Notifier, pF PoolFinder, pR PoolRequester, sig Signer, biodLister biodlisting.Service, robotKey, robotPvKey string, checks map[TransactionType]Checker) Service {
	return service{n, pF, pR, sig, biodLister, robotKey, robotPvKey, checks}
}

func (s service) LeadMining(txHash string, addr string, data interface{}, vPool Pool, txType TransactionType, biodSig string) error {
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

	masterValid, valids, err := s.mine(txHash, data, biodSig, lastVPool, vPool, txType)
	if err != nil {
		if err == ErrInvalidTransaction {
			if err := s.notif.NotifyTransactionStatus(txHash, TxInvalid); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := s.notif.NotifyTransactionStatus(txHash, TxApproved); err != nil {
		return err
	}

	endorsement := datamining.NewEndorsement(time.Now(), txHash, masterValid, valids)
	if err := s.poolR.RequestStorage(sPool, data, endorsement, txType); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxReplicated); err != nil {
		return err
	}

	return s.requestUnlock(txHash, addr, lastVPool)
}

func (s service) findPools(addr string) (Pool, Pool, error) {
	lastVPool, err := s.poolF.FindLastValidationPool(addr)
	if err != nil {
		return nil, nil, err
	}

	sPool, err := s.poolF.FindStoragePool()
	if err != nil {
		return nil, nil, err
	}

	return lastVPool, sPool, nil
}

func (s service) requestLock(txHash string, addr string, lastVPool Pool) error {
	//Build lock transaction
	lock := lock.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey, Address: addr}
	sigLock, err := s.signer.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}

	if err := s.poolR.RequestLock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxLocked); err != nil {
		return err
	}

	return nil
}

func (s service) requestUnlock(txHash string, addr string, lastVPool Pool) error {
	//Build unlock transaction
	lock := lock.TransactionLock{TxHash: txHash, MasterRobotKey: s.robotKey, Address: addr}
	sigLock, err := s.signer.SignLock(lock, s.robotPvKey)
	if err != nil {
		return err
	}

	if err := s.poolR.RequestUnlock(lastVPool, lock, sigLock); err != nil {
		return err
	}
	if err := s.notif.NotifyTransactionStatus(txHash, TxUnlocked); err != nil {
		return err
	}

	return nil
}

func (s service) mine(txHash string, data interface{}, biodSig string, lastVPool, vPool Pool, txType TransactionType) (datamining.MasterValidation, []datamining.Validation, error) {
	//Execute transaction specific master checks
	if err := s.checks[txType].CheckAsMaster(txHash, data); err != nil {
		return nil, nil, err
	}

	//Execute the Proof of Work
	masterValid, err := NewPOW(s.biodLister, s.signer, s.robotKey, s.robotPvKey).Execute(txHash, biodSig, lastVPool)
	if err != nil {
		return nil, nil, err
	}

	//Ask a pool of peers to validate the transaction
	valids, err := s.poolR.RequestValidations(vPool, data, txType)
	if err != nil {
		return nil, nil, err
	}

	//Check if the validations passed
	for _, v := range valids {
		if v.Status() == datamining.ValidationKO {
			return nil, nil, ErrInvalidTransaction
		}
	}

	return masterValid, valids, nil
}

func (s service) Validate(txHash string, data interface{}, txType TransactionType) (datamining.Validation, error) {
	if err := s.checks[txType].CheckAsSlave(txHash, data); err != nil {
		if err == ErrInvalidTransaction {
			return s.buildValidation(datamining.ValidationKO)
		}
		return nil, err
	}
	return s.buildValidation(datamining.ValidationOK)
}

func (s service) buildValidation(status datamining.ValidationStatus) (datamining.Validation, error) {
	v := UnsignedValidation{
		PublicKey: s.robotKey,
		Status:    status,
		Timestamp: time.Now(),
	}
	signature, err := s.signer.SignValidation(v, s.robotPvKey)
	if err != nil {
		return nil, err
	}
	return datamining.NewValidation(v.Status, v.Timestamp, v.PublicKey, signature), nil
}
