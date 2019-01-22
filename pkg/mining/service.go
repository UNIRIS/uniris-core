package mining

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/crypto"
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
	pooler       pooling.Service
	poolR        pooling.PoolRequester
	lister       listing.Service
	minerPubKey  string
	minerPrivKey string
}

//NewService creates a new mining service
func NewService(pool pooling.Service, pR pooling.PoolRequester, l listing.Service, minerPubK string, minerPvKey string) Service {
	return Service{
		pooler:       pool,
		poolR:        pR,
		lister:       l,
		minerPubKey:  minerPubK,
		minerPrivKey: minerPvKey,
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

	errChan := make(chan error)
	go func() {
		for err := range errChan {
			log.Print(err)
		}
	}()

	go s.leadMining(tx, minValids, errChan)
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
		return s.buildMinerValidation(uniris.ValidationKO, time.Now(), s.minerPubKey)
	}
	return s.buildMinerValidation(uniris.ValidationOK, time.Now(), s.minerPubKey)
}

func (s Service) leadMining(tx uniris.Transaction, minValids int, errChan chan error) {
	lastValidationPool, validationPool, storagePool, err := s.findPools(tx)
	if err != nil {
		errChan <- err
		return
	}

	if err := s.poolR.RequestTransactionLock(lastValidationPool, tx.TransactionHash(), tx.Address()); err != nil {
		errChan <- err
		return
	}

	log.Printf("Transaction %s is locked\n", tx.TransactionHash())

	masterValid, confirmValids, err := s.mineTransaction(tx, validationPool, lastValidationPool, minValids)
	if err != nil {
		errChan <- err
		return
	}
	minedTx := uniris.NewMinedTransaction(tx, masterValid, confirmValids)
	log.Printf("Transaction %s is validated \n", tx.TransactionHash())

	if err := s.requestTransactionStorage(minedTx, storagePool, lastValidationPool); err != nil {
		errChan <- err
		return
	}
}

func (s Service) findPools(tx uniris.Transaction) (lastValidationPool, validationPool, storagePool pooling.Pool, err error) {
	lastValidationPool, err = s.pooler.FindLastValidationPool(tx.Address())
	if err != nil {
		return
	}

	validationPool, err = s.pooler.FindValidationPool(tx.TransactionHash())
	if err != nil {
		return
	}

	storagePool, err = s.pooler.FindStoragePool(tx.Address())
	if err != nil {
		return
	}

	return lastValidationPool, validationPool, storagePool, err
}

func (s Service) mineTransaction(tx uniris.Transaction, vPool, lastVPool pooling.Pool, minValids int) (mv uniris.MasterValidation, confirms []uniris.MinerValidation, err error) {
	if err = tx.CheckTransactionIntegrity(); err != nil {
		return
	}

	preValidation, pow, err := s.preValidation(tx)
	if err != nil {
		return
	}

	confirmations, err := s.requestConfirmations(tx, vPool, minValids)
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

	return uniris.NewMasterValidation(lastVPool.Peers(), pow, preValidation), confirms, nil
}

func (s Service) requestConfirmations(tx uniris.Transaction, vPool pooling.Pool, minValids int) ([]uniris.MinerValidation, error) {
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
			if replies == len(vPool.Peers()) {
				wg.Done()
				break
			}
		}
	}()
	go s.poolR.RequestTransactionValidations(vPool, tx, validChan, replyChan)
	wg.Wait()

	return validations, nil

}

func (s Service) requestTransactionStorage(tx uniris.Transaction, sP, lastP pooling.Pool) error {
	//Get minimum replicas of the transaction hash
	minReplicas := 1 //TODO:

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
	if err := s.poolR.RequestTransactionUnlock(lastP, tx.TransactionHash(), tx.Address()); err != nil {
		return err
	}
	return nil
}

func (s Service) preValidation(tx uniris.Transaction) (v uniris.MinerValidation, pow string, err error) {
	pow, err = s.performPow(tx)
	if err != nil {
		return
	}

	validStatus := uniris.ValidationKO
	if pow != "" {
		validStatus = uniris.ValidationOK
	}

	v, err = s.buildMinerValidation(validStatus, time.Now(), s.minerPubKey)
	if err != nil {
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

func (s Service) buildMinerValidation(status uniris.ValidationStatus, ts time.Time, pubK string) (v uniris.MinerValidation, err error) {
	vBytes, err := json.Marshal(uniris.NewMinerValidation(status, ts, s.minerPubKey, ""))
	if err != nil {
		return
	}

	signature, err := crypto.Sign(string(vBytes), s.minerPrivKey)
	if err != nil {
		return
	}
	return uniris.NewMinerValidation(status, ts, pubK, signature), nil
}
