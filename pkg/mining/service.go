package mining

import (
	"errors"
	"log"
	"sync"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/pooling"
)

//Signer handle signature creation for miner validation
type Signer interface {
	SignMinerValidation(status uniris.ValidationStatus, timestamp time.Time, pubKey string) (string, error)
}

//ErrInvalidTransaction is returned a transaction is invalid
var ErrInvalidTransaction = errors.New("Invalid transaction")

//Service handle transaction mining
type Service struct {
	pooler      pooling.Service
	poolR       pooling.PoolRequester
	lister      listing.Service
	minerPubKey string
	signer      Signer
	txVerifier  uniris.TransactionVerifier
	hasher      uniris.TransactionHasher
}

//LeadTransactionMining process workflow to lead mining (like elected master peer)
//
//The workflow includes:
// - Locks the transaction
// - Checks (as master)
// - Executes the proof of work
// - Requests validations (as slave)
// - Requests storage
// - Unlocks the transaction
func (s Service) LeadTransactionMining(tx uniris.Transaction, addr string, vPool pooling.Pool, minValids int) {

	log.Printf("Transaction %s is pending\n", tx.TransactionHash())
	//TODO: store the transaction in pending storage

	errorChan := make(chan error)
	go func() {
		for err := range errorChan {
			log.Print(err)
			//TODO: store the transaction in KO storage
			close(errorChan)
		}
	}()

	go func() {
		lastP, err := s.pooler.FindLastValidationPool(addr)
		if err != nil {
			errorChan <- err
			return
		}

		sP, err := s.pooler.FindStoragePool(addr)
		if err != nil {
			errorChan <- err
			return
		}

		if err := s.poolR.RequestTransactionLock(lastP, tx.TransactionHash(), addr); err != nil {
			errorChan <- err
			return
		}

		log.Printf("Transaction %s is locked\n", tx.TransactionHash())

		txValid, err := s.mineTransaction(tx, vPool, lastP, minValids)
		if err != nil {
			errorChan <- err
			return
		}
		tx.AddMining(txValid)
		log.Printf("Transaction %s is validated \n", tx.TransactionHash())

		if err := s.requestTransactionStorage(tx, addr, sP, lastP); err != nil {
			errorChan <- err
			return
		}

		//TODO: remove the transaction from the pending storage
	}()
}

func (s Service) mineTransaction(tx uniris.Transaction, vPool, lastVPool pooling.Pool, minValids int) (m uniris.TransactionMining, err error) {
	if err = s.checkTransactionIntegrity(tx); err != nil {
		return
	}
	masterValidation, pow, err := s.performPow(tx)
	if err != nil {
		return
	}

	validChan := make(chan uniris.MinerValidation)
	validations := make([]uniris.MinerValidation, 0)

	var wg sync.WaitGroup
	wg.Add(minValids)

	go func() {
		for v := range validChan {
			validations = append(validations, v)
			if len(validations) == minValids {
				wg.Done()
			}
		}
	}()
	go s.poolR.RequestTransactionValidations(vPool, tx, validChan)
	wg.Wait()

	if len(validations) < minValids {
		//TODO: improve to avoid transaction failure. maybe ask again the same pool or choose a new pool
		err = ErrInvalidTransaction
		return
	}

	//Check if the validations passed
	nbKO := 0
	for _, v := range validations {
		if v.Status() == uniris.ValidationKO {
			nbKO++
		}
	}
	//TODO: to improve to avoid transaction failure. maybe ask again the same pool or choose a new pool
	if nbKO == len(validations) {
		err = ErrInvalidTransaction
		return
	}

	return uniris.NewTransactionMining(lastVPool.Peers(), pow, masterValidation, validations), nil
}

func (s Service) requestTransactionStorage(tx uniris.Transaction, addr string, sP, lastP pooling.Pool) error {
	//Get minimum replicas of the transaction hash
	minReplicas := 1 //TODO: contact AI service

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
	if err := s.poolR.RequestTransactionUnlock(lastP, tx.TransactionHash(), addr); err != nil {
		return err
	}
	return nil
}

//ValidateTransaction performs checks create a validation (successed or not)
func (s Service) ValidateTransaction(tx uniris.Transaction) (v uniris.MinerValidation, err error) {
	if err = s.checkTransactionIntegrity(tx); err != nil {
		if err == ErrInvalidTransaction {
			return s.buildMinerValidation(uniris.ValidationKO, time.Now(), s.minerPubKey)
		}
		return
	}
	return s.buildMinerValidation(uniris.ValidationOK, time.Now(), s.minerPubKey)
}

func (s Service) checkTransactionIntegrity(tx uniris.Transaction) error {

	hash, err := s.hasher.HashTransaction(tx)
	if err != nil {
		return err
	}
	if hash != tx.TransactionHash() {
		return ErrInvalidTransaction
	}

	ok, err := s.txVerifier.VerifyTransactionSignature(tx, tx.IDPublicKey(), tx.IDSignature())
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidTransaction
	}

	return nil
}

func (s Service) performPow(tx uniris.Transaction) (v uniris.MinerValidation, pow string, err error) {

	emKeys, err := s.lister.ListSharedEmitterKeyPairs()
	if err != nil {
		return
	}

	validStatus := uniris.ValidationKO
	powKey, err := s.findTransactionEmitterPublicKey(tx, emKeys)
	if err != nil {
		return
	}
	if powKey != "" {
		validStatus = uniris.ValidationOK
	}

	v, err = s.buildMinerValidation(validStatus, time.Now(), s.minerPubKey)
	if err != nil {
		return
	}
	return v, powKey, nil
}

func (s Service) findTransactionEmitterPublicKey(tx uniris.Transaction, emKeys []uniris.SharedKeys) (string, error) {
	for _, kp := range emKeys {

		ok, err := s.txVerifier.VerifyTransactionSignature(tx, kp.PublicKey(), tx.EmitterSignature())
		if err != nil {
			return "", err
		}
		if ok {
			return kp.PublicKey(), nil
		}
	}

	return "", nil
}

func (s Service) buildMinerValidation(status uniris.ValidationStatus, ts time.Time, pubK string) (v uniris.MinerValidation, err error) {
	signature, err := s.signer.SignMinerValidation(status, ts, pubK)
	if err != nil {
		return
	}
	return uniris.NewMinerValidation(status, ts, pubK, signature), nil
}
