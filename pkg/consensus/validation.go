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
// - Unlocks the transaction
func LeadMining(tx chain.Transaction, minValids int, poolR PoolRequester, pub, pv string, emR shared.EmitterDatabaseReader) error {
	log.Printf("transaction %s is in progress\n", tx.TransactionHash())

	lastVPool, err := findLastValidationPool(tx.Address(), tx.TransactionType(), poolR)
	if err != nil {
		return err
	}

	sPool, err := FindStoragePool(tx.Address())
	if err != nil {
		return err
	}

	//TODO: ask storage pool to store in in progress

	if err := lockTransaction(tx, lastVPool, poolR, pub); err != nil {
		return err
	}

	go func() {
		minedTx, err := mineTransaction(tx, lastVPool, minValids, pub, pv, emR, poolR)
		if err != nil {
			fmt.Printf("transaction mining failed: %s", err.Error())
			return
		}
		if err := storeTransaction(minedTx, sPool, poolR); err != nil {
			fmt.Printf("transaction storage failed: %s", err.Error())
			return
		}
		if err := unlockTransaction(tx, lastVPool, poolR); err != nil {
			fmt.Printf("transaction unlock failed: %s", err.Error())
			return
		}
	}()

	return nil
}

func lockTransaction(tx chain.Transaction, lastVPool Pool, poolR PoolRequester, masterPublicKey string) error {
	//TODO: find a solution when no last validation pool (example for the first transaction)
	if lastVPool != nil {
		if err := poolR.RequestTransactionLock(lastVPool, tx.TransactionHash(), tx.Address(), masterPublicKey); err != nil {
			return fmt.Errorf("transaction lock failed: %s", err.Error())
		}

		log.Printf("transaction %s is locked\n", tx.TransactionHash())
	}
	return nil
}

func mineTransaction(tx chain.Transaction, lastVPool Pool, minValids int, minerPub, minerPv string, emR shared.EmitterDatabaseReader, poolR PoolRequester) (chain.Transaction, error) {

	vPool, err := FindValidationPool(tx)
	if err != nil {
		return tx, fmt.Errorf("transaction find validation pool failed: %s", err.Error())
	}

	masterValid, err := preValidateTransaction(tx, lastVPool, minValids, minerPub, minerPv, emR)
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
	log.Printf("transaction %s is validated \n", tx.TransactionHash())
	return tx, nil
}

func storeTransaction(tx chain.Transaction, sPool Pool, poolR PoolRequester) error {
	minReplicas := GetMinimumReplicas(tx.TransactionHash())
	if err := poolR.RequestTransactionStorage(sPool, minReplicas, tx); err != nil {
		return fmt.Errorf("transaction storage failed: %s", err.Error())
	}
	log.Printf("transaction %s is stored \n", tx.TransactionHash())
	return nil
}

func unlockTransaction(tx chain.Transaction, lastVPool Pool, poolR PoolRequester) error {
	//TODO: find a solution when no last validation pool (example for the first transaction)
	if lastVPool != nil {
		if err := poolR.RequestTransactionUnlock(lastVPool, tx.TransactionHash(), tx.Address()); err != nil {
			return fmt.Errorf("transaction unlock failed: %s", err.Error())
		}
	}
	return nil
}

//findPools retrieve the needed pools for the transaction mining process (last validation pool, new validation pool and storage pool)
func findPools(tx chain.Transaction, poolR PoolRequester) (lastValidationPool, validationPool, storagePool Pool, err error) {
	lastValidationPool, err = findLastValidationPool(tx.Address(), tx.TransactionType(), poolR)
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
func preValidateTransaction(tx chain.Transaction, lastVPool Pool, minValids int, pub, pv string, emR shared.EmitterDatabaseReader) (chain.MasterValidation, error) {
	if _, err := tx.IsValid(); err != nil {
		return chain.MasterValidation{}, err
	}

	pow, err := proofOfWork(tx, emR)
	if err != nil {
		return chain.MasterValidation{}, err
	}
	validStatus := chain.ValidationKO
	if pow != "" {
		validStatus = chain.ValidationOK
	}
	preValid, err := buildMinerValidation(validStatus, pub, pv)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	lastMinersKeys := make([]string, 0)
	for _, pm := range lastVPool {
		lastMinersKeys = append(lastMinersKeys, pm.PublicKey())
	}
	masterValid, err := chain.NewMasterValidation(lastMinersKeys, pow, preValid)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	return masterValid, nil
}

func proofOfWork(tx chain.Transaction, emR shared.EmitterDatabaseReader) (pow string, err error) {
	emKeys, err := emR.EmitterKeys()
	if err != nil {
		return
	}

	txBytes, err := tx.MarshalBeforeSignature()
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

func findLastValidationPool(txAddr string, txType chain.TransactionType, req PoolRequester) (Pool, error) {

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

	pm := make([]PoolMember, 0)
	for _, key := range tx.MasterValidation().PreviousTransactionMiners() {
		//TODO: find ip address and port
		pm = append(pm, PoolMember{
			pubK: key,
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		})

	}

	return Pool(pm), nil
}

func requestValidations(tx chain.Transaction, masterValid chain.MasterValidation, vPool Pool, minValids int, poolR PoolRequester) ([]chain.MinerValidation, error) {
	validations, err := poolR.RequestTransactionValidations(vPool, tx, minValids, masterValid)
	if err != nil {
		return nil, err
	}

	if !IsValidationConsensuReach(validations) {
		return nil, errors.New("invalid transaction")
	}

	return validations, nil
}

//ConfirmTransactionValidation approve the transaction validation by master and ensure its integrity
func ConfirmTransactionValidation(tx chain.Transaction, masterV chain.MasterValidation, pub, pv string) (chain.MinerValidation, error) {

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

	return buildMinerValidation(status, pub, pv)
}

func buildMinerValidation(s chain.ValidationStatus, pub, pv string) (chain.MinerValidation, error) {
	b, err := json.Marshal(map[string]interface{}{
		"status":     s,
		"timestamp":  time.Now().Unix(),
		"public_key": pub,
	})
	if err != nil {
		return chain.MinerValidation{}, err
	}
	sig, err := crypto.Sign(string(b), pv)
	if err != nil {
		return chain.MinerValidation{}, err
	}
	return chain.NewMinerValidation(s, time.Now(), pub, sig)
}

//GetMinimumValidation returns the validation from a transaction hash
func GetMinimumValidation(txHash string) int {
	return 1
}

// IsValidationConsensuReach determinates if for the miner validations the consensus is reached
func IsValidationConsensuReach(valids []chain.MinerValidation) bool {
	//TODO: maybe to improve
	for _, v := range valids {
		if v.Status() == chain.ValidationKO {
			return false
		}
	}
	return true
}
