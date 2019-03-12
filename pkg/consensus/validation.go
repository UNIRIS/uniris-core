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
func LeadMining(tx chain.Transaction, minValids int, wHeaders []chain.NodeHeader, poolR PoolRequester, pub, pv string, emR shared.EmitterReader) error {
	log.Printf("transaction %s is in progress\n", tx.TransactionHash())

	lastVPool, err := findLastValidationPool(tx.Address(), tx.TransactionType(), poolR)
	if err != nil {
		return err
	}

	sPool, err := FindStoragePool(tx.Address())
	if err != nil {
		return err
	}

	if err := poolR.RequestTransactionTimeLock(sPool, tx.TransactionHash(), tx.Address(), pub); err != nil {
		return fmt.Errorf("transaction lock failed: %s", err.Error())
	}

	log.Printf("transaction %s is locked\n", tx.TransactionHash())

	go func() {
		minedTx, err := mineTransaction(tx, wHeaders, sPool, lastVPool, minValids, pub, pv, emR, poolR)
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

func mineTransaction(tx chain.Transaction, wHeaders []chain.NodeHeader, sPool Pool, lastVPool Pool, minValids int, nodePub, nodePv string, emR shared.EmitterReader, poolR PoolRequester) (chain.Transaction, error) {

	vPool, err := FindValidationPool(tx)
	if err != nil {
		return tx, fmt.Errorf("transaction find validation pool failed: %s", err.Error())
	}

	masterValid, err := preValidateTransaction(tx, wHeaders, sPool, vPool, lastVPool, minValids, nodePub, nodePv, emR)
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
func preValidateTransaction(tx chain.Transaction, wHeaders []chain.NodeHeader, sPool Pool, vPool Pool, lastVPool Pool, minValids int, pub, pv string, emR shared.EmitterReader) (chain.MasterValidation, error) {
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
	preValid, err := buildValidation(validStatus, pub, pv)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	lastsKeys := make([]string, 0)
	for _, pm := range lastVPool {
		lastsKeys = append(lastsKeys, pm.PublicKey())
	}

	vHeaders, sHeaders := buildHeaders(vPool, sPool)
	masterValid, err := chain.NewMasterValidation(lastsKeys, pow, preValid, wHeaders, vHeaders, sHeaders)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	return masterValid, nil
}

func buildHeaders(vPool Pool, sPool Pool) (vHeaders []chain.NodeHeader, sHeaders []chain.NodeHeader) {
	for _, n := range vPool {
		//TODO: retrieve real value (patch, is unreachable, is OK)
		vHeaders = append(vHeaders, chain.NewNodeHeader(n.PublicKey(), true, true, 0, true))
	}

	for _, n := range sPool {
		//TODO: retrieve real value (patch, is unreachable)
		sHeaders = append(sHeaders, chain.NewNodeHeader(n.PublicKey(), true, true, 0, true))
	}

	return
}

func proofOfWork(tx chain.Transaction, emR shared.EmitterReader) (pow string, err error) {
	emKeys, err := emR.EmitterKeys()
	if err != nil {
		return
	}

	txBytes, err := tx.MarshalBeforeEmitterSignature()
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

	pm := make([]Node, 0)
	for _, key := range tx.MasterValidation().PreviousValidationNodes() {
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

	if !IsValidationConsensuReach(validations) {
		return nil, errors.New("invalid transaction")
	}

	return validations, nil
}

//ConfirmTransactionValidation approve the transaction validation by master and ensure its integrity
func ConfirmTransactionValidation(tx chain.Transaction, masterV chain.MasterValidation, pub, pv string) (chain.Validation, error) {

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

func buildValidation(s chain.ValidationStatus, pub, pv string) (chain.Validation, error) {
	b, err := json.Marshal(map[string]interface{}{
		"status":     s,
		"timestamp":  time.Now().Unix(),
		"public_key": pub,
	})
	if err != nil {
		return chain.Validation{}, err
	}
	sig, err := crypto.Sign(string(b), pv)
	if err != nil {
		return chain.Validation{}, err
	}
	return chain.NewValidation(s, time.Now(), pub, sig)
}

//GetMinimumValidation returns the validation from a transaction hash
func GetMinimumValidation(txHash string) int {
	return 1
}

// IsValidationConsensuReach determinates if for the node validations the consensus is reached
func IsValidationConsensuReach(valids []chain.Validation) bool {
	//TODO: maybe to improve
	for _, v := range valids {
		if v.Status() == chain.ValidationKO {
			return false
		}
	}
	return true
}
