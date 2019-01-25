package mining

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/electing"
	"github.com/uniris/uniris-core/pkg/listing"
)

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//Service handle transaction mining
type Service struct {
	poolR        uniris.PoolRequester
	lister       listing.Service
	minerPubKey  string
	minerPrivKey string
	currentIP    string
}

//NewService creates a new mining service
func NewService(l listing.Service, poolR uniris.PoolRequester, minerPubK string, minerPvKey string, currentIP string) Service {
	return Service{
		poolR:        poolR,
		lister:       l,
		minerPubKey:  minerPubK,
		minerPrivKey: minerPvKey,
		currentIP:    currentIP,
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
func (s Service) LeadTransactionValidation(tx uniris.Transaction, minValids int) {
	log.Printf("Transaction %s is pending\n", tx.TransactionHash())
	go func() {
		lastValidationPool, validationPool, storagePool, err := s.findPools(tx)
		if err != nil {
			log.Printf("Transaction pool finding failed - Error: %s\n", err.Error())
			return
		}

		lockTx, err := uniris.NewLock(tx.TransactionHash(), tx.Address(), s.currentIP)
		if err != nil {
			log.Printf("Lock creation failed - Error: %s\n", err.Error())
			return
		}
		if err := s.poolR.RequestTransactionLock(lastValidationPool, lockTx); err != nil {
			log.Printf("Transaction locking failed - Error: %s\n", err.Error())
			return
		}

		log.Printf("Transaction %s is locked\n", tx.TransactionHash())

		masterValid, confirmValids, err := s.mineTransaction(tx, validationPool, lastValidationPool, minValids)
		if err != nil {
			log.Printf("Transaction mining failed - Error: %s\n", err.Error())
			return
		}
		if err := tx.AddMining(masterValid, confirmValids); err != nil {
			log.Printf("Transaction mining is invalid - Error: %s\n", err.Error())
			return
		}
		log.Printf("Transaction %s is validated \n", tx.TransactionHash())

		if err := s.requestTransactionStorage(tx, storagePool, lastValidationPool); err != nil {
			log.Printf("Transaction storage failed - Error: %s\n", err.Error())
			return
		}
	}()
}

//ConfirmTransactionValidation confirms the transaction by providing a validation including:
// - Check the proof of work is valid
// - Check the mining (by the master peer) is valid (signatures checks)
// - Check the transaction integrity
func (s Service) ConfirmTransactionValidation(tx uniris.Transaction) (v uniris.MinerValidation, err error) {
	if err = tx.CheckProofOfWork(); err != nil {
		return
	}
	if err = tx.MasterValidation().Validation().CheckValidation(); err != nil {
		return
	}

	if err = tx.CheckTransactionIntegrity(); err != nil {
		valid, err = uniris.NewMinerValidation(uniris.ValidationKO, time.Now(), s.minerPubKey, "")
		if err != nil {
			return
		}
	} else {
		valid, err = uniris.NewMinerValidation(uniris.ValidationOK, time.Now(), s.minerPubKey, "")
		if err != nil {
			return
		}
	}
	if err = v.Sign(s.minerPrivKey); err != nil {
		return
	}
	return valid, nil
}

func (s Service) findPools(tx uniris.Transaction) (lastValidationPool, validationPool, storagePool uniris.Pool, err error) {
	lastValidationPool, err = electing.FindLastValidationPool(tx.Address(), s.poolR)
	if err != nil {
		return
	}

	validationPool, err = electing.FindValidationPool(tx.TransactionHash())
	if err != nil {
		return
	}

	storagePool, err = electing.FindStoragePool(tx.Address())
	if err != nil {
		return
	}

	return lastValidationPool, validationPool, storagePool, err
}

func (s Service) mineTransaction(tx uniris.Transaction, vPool, lastVPool uniris.Pool, minValids int) (masterValid uniris.MasterValidation, confirms []uniris.MinerValidation, err error) {
	if err = tx.CheckTransactionIntegrity(); err != nil {
		return
	}

	preValidation, pow, err := s.preValidateTx(tx)
	if err != nil {
		return
	}

	masterValid, err = uniris.NewMasterValidation(lastVPool, pow, preValidation)
	if err != nil {
		return
	}

	confirmations, err := s.requestConfirmations(tx, masterValid, vPool, minValids)
	if err != nil {
		return
	}

	if len(confirmations) < minValids {
		//TODO: improve to avoid transaction failure. maybe ask again the same pool or choose a new pool
		err = ErrInvalidTransaction
		return
	}

	//Check if the validations passed
	nbKO := 0
	for _, v := range confirmations {
		if v.Status() == uniris.ValidationKO {
			nbKO++
		}
	}
	//TODO: to improve to avoid transaction failure. maybe ask again the same pool or choose a new pool
	if nbKO == len(confirmations) {
		err = ErrInvalidTransaction
		return
	}

	return masterValid, confirms, nil
}

func (s Service) requestConfirmations(tx uniris.Transaction, masterValid uniris.MasterValidation, vPool uniris.Pool, minValids int) ([]uniris.MinerValidation, error) {
	validChan := make(chan uniris.MinerValidation)
	replyChan := make(chan bool)
	validations := make([]uniris.MinerValidation, 0)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for v := range validChan {
			validations = append(validations, v)
		}
	}()

	//Listen replies and stop listening when the minimum validation confirmation is reached
	//or when all the peers inside the pool have been reached
	go func() {
		replies := 0
		for range replyChan {
			replies++
			if replies == minValids {
				wg.Done()
				break
			}
			if replies == len(vPool) {
				wg.Done()
				break
			}
		}
	}()
	go s.poolR.RequestTransactionValidations(vPool, tx, masterValid, validChan, replyChan)
	wg.Wait()

	return validations, nil

}

func (s Service) requestTransactionStorage(tx uniris.Transaction, sP, lastP uniris.Pool) error {
	minReplicas := getMinimumReplicas(tx.TransactionHash())
	storageAck := make(chan bool, 0)

	var wg sync.WaitGroup
	wg.Add(minReplicas)

	go func() {
		for range storageAck {
			wg.Done()
		}
	}()

	go s.poolR.RequestTransactionStorage(sP, tx, storageAck)
	wg.Wait()

	log.Printf("Transaction %s is stored \n", tx.TransactionHash())

	unlockTx, err := uniris.NewLock(tx.TransactionHash(), tx.Address(), s.currentIP)
	if err != nil {
		return fmt.Errorf("Unlock creation failed - Error: %s", err.Error())
	}
	if err := s.poolR.RequestTransactionUnlock(lastP, unlockTx); err != nil {
		return err
	}
	return nil
}

func (s Service) preValidateTx(tx uniris.Transaction) (v uniris.MinerValidation, pow string, err error) {
	pow, err = s.performPow(tx)
	if err != nil {
		return
	}

	validStatus := uniris.ValidationKO
	if pow != "" {
		validStatus = uniris.ValidationOK
	}

	v, err = uniris.NewMinerValidation(validStatus, time.Now(), s.minerPubKey, "")
	if err != nil {
		return
	}
	if err = v.Sign(s.minerPrivKey); err != nil {
		return
	}
	return v, pow, nil
}

func (s Service) performPow(tx uniris.Transaction) (pow string, err error) {
	emKeys, err := s.lister.ListSharedEmitterKeyPairs()
	if err != nil {
		return
	}

	txBytes, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}

	for _, kp := range emKeys {
		err := crypto.VerifySignature(string(txBytes), kp.PublicKey(), tx.EmitterSignature())
		if err == nil {
			return kp.PublicKey(), nil
		}
	}

	return "", nil
}
