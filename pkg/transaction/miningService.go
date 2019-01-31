package transaction

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionLock asks a pool to lock a transaction using the address related
	RequestTransactionLock(pool Pool, txLock Lock) error

	//RequestTransactionUnlock asks a pool to unlock a transaction using the address related
	RequestTransactionUnlock(pool Pool, txLock Lock) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx Transaction, masterValid MasterValidation, validChan chan<- MinerValidation)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, tx Transaction, ackChan chan<- bool)
}

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//MiningService handles transaction mining
type MiningService struct {
	poolR     PoolRequester
	poolFSrv  PoolFindingService
	sharedSrv shared.Service
	minerIP   string
	minerPubK string
	minerPvk  string
}

//NewMiningService creates a new transaction mining service
func NewMiningService(pR PoolRequester, pfS PoolFindingService, sS shared.Service, mIP string, mPub string, mPv string) MiningService {
	return MiningService{
		poolR:     pR,
		poolFSrv:  pfS,
		sharedSrv: sS,
		minerIP:   mIP,
		minerPubK: mPub,
		minerPvk:  mPv,
	}
}

//LeadTransactionValidation validate the transaction as a master peer and lead the mining workflow
//
//The workflow includes:
// - Locks the transaction
// - Pre-validate (master validation)
// - Executes the proof of work
// - Requests validation confirmations
// - Requests storage
// - Unlocks the transaction
func (s MiningService) LeadTransactionValidation(tx Transaction, minValids int) {
	log.Printf("transaction %s is pending\n", tx.txHash)
	//TODO: ask storage pool to store in pending

	go func() {
		lastValidationPool, validationPool, storagePool, err := s.findPools(tx)
		if err != nil {
			log.Printf("transaction pool finding failed: %s\n", err.Error())
			return
		}

		//TODO: find a solution when no last validation pool (example for the first transaction)
		if lastValidationPool != nil {
			lockTx, err := NewLock(tx.TransactionHash(), tx.address, s.minerPubK)
			if err != nil {
				return
			}
			if err := s.poolR.RequestTransactionLock(lastValidationPool, lockTx); err != nil {
				log.Printf("transaction lock failed: %s\n", err.Error())
				return
			}

			log.Printf("transaction %s is locked\n", tx.txHash)
		}

		masterValid, confirmValids, err := s.mineTransaction(tx, validationPool, lastValidationPool, minValids)
		if err != nil {
			log.Printf("transaction mining failed: %s\n", err.Error())
			return
		}
		if err := tx.AddMining(masterValid, confirmValids); err != nil {
			log.Printf("transaction mining is invalid: %s\n", err.Error())
			return
		}
		log.Printf("transaction %s is validated \n", tx.txHash)

		s.requestTransactionStorage(tx, storagePool)
		log.Printf("transaction %s is stored \n", tx.txHash)

		//TODO: find a solution when no last validation pool (example for the first transaction)
		if lastValidationPool != nil {
			unlockTx, err := NewLock(tx.txHash, tx.address, s.minerPubK)
			if err != nil {
				return
			}
			if err := s.poolR.RequestTransactionUnlock(lastValidationPool, unlockTx); err != nil {
				log.Printf("transaction lock failed: %s", err.Error())
				return
			}
		}
	}()
}

func (s MiningService) findPools(tx Transaction) (lastValidationPool, validationPool, storagePool Pool, err error) {
	lastValidationPool, err = s.poolFSrv.FindLastValidationPool(tx.Address(), tx.Type())
	if err != nil {
		return
	}

	validationPool, err = s.poolFSrv.FindValidationPool(tx.TransactionHash())
	if err != nil {
		return
	}

	storagePool, err = s.poolFSrv.FindStoragePool(tx.Address())
	if err != nil {
		return
	}

	return lastValidationPool, validationPool, storagePool, err
}

func (s MiningService) mineTransaction(tx Transaction, vPool, lastVPool Pool, minValids int) (masterValid MasterValidation, confirms []MinerValidation, err error) {
	if _, err = tx.IsValid(); err != nil {
		return
	}

	preValidation, pow, err := s.preValidateTx(tx)
	if err != nil {
		return
	}

	masterValid, err = NewMasterValidation(lastVPool, pow, preValidation)
	if err != nil {
		return
	}

	validations, err := s.requestValidations(tx, masterValid, vPool, minValids)
	if err != nil {
		return
	}

	for _, v := range validations {
		if v.Status() == ValidationKO {
			err = ErrInvalidTransaction
			return
		}
	}

	return masterValid, validations, nil
}

func (s MiningService) preValidateTx(tx Transaction) (MinerValidation, string, error) {

	pow, err := s.performPow(tx)
	if err != nil {
		return MinerValidation{}, "", err
	}

	validStatus := ValidationKO
	if pow != "" {
		validStatus = ValidationOK
	}

	v, err := s.buildMinerValidation(validStatus)
	if err != nil {
		return MinerValidation{}, "", err
	}
	return v, pow, nil
}

func (s MiningService) performPow(tx Transaction) (pow string, err error) {
	emKeys, err := s.sharedSrv.ListSharedEmitterKeyPairs()
	if err != nil {
		return
	}

	txBytes, err := tx.MarshalBeforeSignature()
	if err != nil {
		return "", err
	}

	for _, kp := range emKeys {
		err := crypto.VerifySignature(string(txBytes), kp.PublicKey(), tx.emSig)
		if err == nil {
			return kp.PublicKey(), nil
		}
	}

	return "", nil
}

func (s MiningService) requestTransactionStorage(tx Transaction, sP Pool) {
	minReplicas := s.getMinimumReplicas(tx.txHash)
	storageAck := make(chan bool, 0)

	var wg sync.WaitGroup
	wg.Add(minReplicas)

	go func() {
		for range storageAck {
			wg.Done()
		}
	}()

	//TODO: provide a context to handle timeout

	go s.poolR.RequestTransactionStorage(sP, tx, storageAck)
	wg.Wait()

}

func (s MiningService) requestValidations(tx Transaction, masterValid MasterValidation, vPool Pool, minValids int) ([]MinerValidation, error) {
	validChan := make(chan MinerValidation)
	validations := make([]MinerValidation, 0)

	defer func() {
		close(validChan)
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	//Listen validations and stop when the minimum has been reached
	go func() {
		for v := range validChan {
			validations = append(validations, v)
			if len(validations) == minValids {
				wg.Done()
				break
			}
		}
	}()

	//TODO: provide a context to handle timeout and returns error

	go s.poolR.RequestTransactionValidations(vPool, tx, masterValid, validChan)
	wg.Wait()

	return validations, nil
}

//ValidateTransaction provide a validation including:
// - Check the transaction validity
// - Check the master validation
// - Check the transaction integrity
func (s MiningService) ValidateTransaction(tx Transaction, mv MasterValidation) (MinerValidation, error) {
	if _, err := tx.IsValid(); err != nil {
		return s.buildMinerValidation(ValidationKO)
	}

	if _, err := mv.IsValid(); err != nil {
		return s.buildMinerValidation(ValidationKO)
	}

	return s.buildMinerValidation(ValidationOK)
}

func (s MiningService) buildMinerValidation(status ValidationStatus) (MinerValidation, error) {
	v := MinerValidation{
		minerPubk: s.minerPubK,
		status:    status,
		timestamp: time.Now(),
	}
	b, err := json.Marshal(v)
	if err != nil {
		return MinerValidation{}, err
	}
	sig, err := crypto.Sign(string(b), s.minerPvk)
	if err != nil {
		return MinerValidation{}, err
	}
	v.minerSig = sig
	return v, nil
}

//GetMinimumTransactionValidation returns the validation from a transaction hash
func (s MiningService) GetMinimumTransactionValidation(txHash string) int {
	return 1
}

func (s MiningService) getMinimumReplicas(txHash string) int {
	return 1
}
