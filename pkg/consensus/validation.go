package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

//LeadMining lead the mining workflow
//
//The workflow includes:
// - Locks the transaction
// - Pre-validate (master validation)
// - Executes the proof of work
// - Requests validation confirmations
// - Requests storage
func LeadMining(tx chain.Transaction, minValids int, poolR PoolRequester, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, emR shared.EmitterDatabaseReader) error {
	log.Printf("transaction %x is in progress\n", tx.TransactionHash())

	if !tx.Address().IsValid() {
		return errors.New("invalid transaction address")
	}

	lastVPool, err := findLastValidationPool(tx.Address().Digest(), tx.TransactionType(), poolR)
	if err != nil {
		return err
	}

	sPool, err := FindStoragePool(tx.Address().Digest())
	if err != nil {
		return err
	}

<<<<<<< HEAD
	if err := poolR.RequestTransactionTimeLock(sPool, tx.TransactionHash(), tx.Address(), nodePub); err != nil {
=======
	//TODO: ask storage pool to store in in progress

	if err := poolR.RequestTransactionLock(sPool, tx.TransactionHash(), tx.Address(), nodePub); err != nil {
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash
		return fmt.Errorf("transaction lock failed: %s", err.Error())
	}

	log.Printf("transaction %x is locked\n", tx.TransactionHash())

	go func() {
		minedTx, err := mineTransaction(tx, lastVPool, minValids, nodePub, nodePv, emR, poolR)
		if err != nil {
			fmt.Printf("transaction mining failed: %s\n", err.Error())
			return
		}
		if err := storeTransaction(minedTx, sPool, poolR); err != nil {
			fmt.Printf("transaction storage failed: %s\n", err.Error())
			return
		}
	}()

	return nil
}

func mineTransaction(tx chain.Transaction, lastVPool Pool, minValids int, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, emR shared.EmitterDatabaseReader, poolR PoolRequester) (chain.Transaction, error) {

	vPool, err := FindValidationPool(tx)
	if err != nil {
		return tx, fmt.Errorf("transaction find validation pool failed: %s", err.Error())
	}

	masterValid, err := preValidateTransaction(tx, lastVPool, minValids, nodePub, nodePv, emR)
	if err != nil {
		return tx, fmt.Errorf("transaction pre-validation failed: %s", err.Error())
	}
	confirmValids, err := requestValidations(tx, masterValid, vPool, minValids, poolR)
	if err != nil {
		return tx, fmt.Errorf("transaction validation confirmations failed: %s", err.Error())
	}
	if err := tx.Mined(masterValid, confirmValids); err != nil {
		return tx, fmt.Errorf("transaction mining is invalid: %s", err.Error())
	}
	log.Printf("transaction %x is validated \n", tx.TransactionHash())
	return tx, nil
}

func storeTransaction(tx chain.Transaction, sPool Pool, poolR PoolRequester) error {
	minReplicas := GetMinimumReplicas(tx.TransactionHash().Digest())
	if err := poolR.RequestTransactionStorage(sPool, minReplicas, tx); err != nil {
		return fmt.Errorf("transaction storage failed: %s", err.Error())
	}
	log.Printf("transaction %x is stored \n", tx.TransactionHash())
	return nil
}

//findPools retrieve the needed pools for the transaction mining process (last validation pool, new validation pool and storage pool)
func findPools(tx chain.Transaction, poolR PoolRequester) (lastValidationPool, validationPool, storagePool Pool, err error) {
	lastValidationPool, err = findLastValidationPool(tx.Address().Digest(), tx.TransactionType(), poolR)
	if err != nil {
		return
	}

	validationPool, err = FindValidationPool(tx)
	if err != nil {
		return
	}

	storagePool, err = FindStoragePool(tx.Address())
	if err != nil {
		return
	}

	return
}

//preValidateTransaction checks the incoming transaction as master node by ensure the transaction integrity and perform the proof of work. A valiation will result from this action
func preValidateTransaction(tx chain.Transaction, lastVPool Pool, minValids int, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, emR shared.EmitterDatabaseReader) (chain.MasterValidation, error) {
	if _, err := tx.IsValid(); err != nil {
		return chain.MasterValidation{}, err
	}

	pow, err := proofOfWork(tx, emR)
	if err != nil {
		return chain.MasterValidation{}, err
	}
	validStatus := chain.ValidationKO
	if pow != nil {
		validStatus = chain.ValidationOK
	}
	preValid, err := buildValidation(validStatus, nodePub, nodePv)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	lastsKeys := make([]crypto.PublicKey, 0)
	for _, pm := range lastVPool {
		lastsKeys = append(lastsKeys, pm.PublicKey())
	}
	masterValid, err := chain.NewMasterValidation(lastsKeys, pow, preValid)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	return masterValid, nil
}

func proofOfWork(tx chain.Transaction, emR shared.EmitterDatabaseReader) (pow crypto.PublicKey, err error) {
	emKeys, err := emR.EmitterKeys()
	if err != nil {
		return
	}

	txBytes, err := tx.MarshalBeforeEmitterSignature()
	if err != nil {
		return nil, err
	}

	for _, kp := range emKeys {
		if ok := kp.PublicKey().Verify(txBytes, tx.EmitterSignature()); ok {
			return kp.PublicKey(), nil
		}
	}

	return nil, nil
}

func findLastValidationPool(txAddr []byte, txType chain.TransactionType, req PoolRequester) (Pool, error) {

	sPool, err := FindStoragePool(txAddr)
	if err != nil {
		return nil, err
	}

	tx, err := req.RequestLastTransaction(sPool, txAddr, txType)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	pm := make([]Node, 0)
	for _, key := range tx.MasterValidation().PreviousTransactionNodes() {
		//TODO: find ip address and port
		pm = append(pm, Node{
			publicKey: key,
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
		})

	}

	return Pool(pm), nil
}

func requestValidations(tx chain.Transaction, masterValid chain.MasterValidation, vPool Pool, minValids int, poolR PoolRequester) ([]chain.Validation, error) {
	validations, err := poolR.RequestTransactionValidations(vPool, tx, minValids, masterValid)
	if err != nil {
		return nil, err
	}

	if !IsValidationConsensusReach(validations) {
		return nil, errors.New("invalid transaction")
	}

	return validations, nil
}

//ConfirmTransactionValidation approve the transaction validation by master and ensure its integrity
func ConfirmTransactionValidation(tx chain.Transaction, masterV chain.MasterValidation, pub crypto.PublicKey, pv crypto.PrivateKey) (chain.Validation, error) {

	var status chain.ValidationStatus

	if _, err := tx.IsValid(); err != nil {
		log.Printf("Transaction validation confirmation failed: %s\n", err.Error())
		status = chain.ValidationKO
	} else if _, err := masterV.IsValid(); err != nil {
		log.Printf("Transaction master validation confirmation failed: %s\n", err.Error())
		status = chain.ValidationKO
	} else {
		status = chain.ValidationOK
	}

	return buildValidation(status, pub, pv)
}

func buildValidation(s chain.ValidationStatus, pub crypto.PublicKey, pv crypto.PrivateKey) (chain.Validation, error) {
	pubBytes, err := pub.Marshal()
	if err != nil {
		return chain.Validation{}, err
	}

	vBytes, err := json.Marshal(map[string]interface{}{
		"status":     s,
		"timestamp":  time.Now().Unix(),
		"public_key": pubBytes,
	})
	if err != nil {
		return chain.Validation{}, err
	}

	vSig, err := pv.Sign(vBytes)
	if err != nil {
		return chain.Validation{}, err
	}
	return chain.NewValidation(s, time.Now(), pub, vSig)
}

//GetMinimumValidation returns the validation from a transaction hash
func GetMinimumValidation(txHash []byte) int {
	return 1
}

// IsValidationConsensusReach determinates if for the node validations the consensus is reached
func IsValidationConsensusReach(valids []chain.Validation) bool {
	//TODO: maybe to improve
	for _, v := range valids {
		if v.Status() == chain.ValidationKO {
			return false
		}
	}
	return true
}
