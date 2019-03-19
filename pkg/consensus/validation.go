package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
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
func LeadMining(tx chain.Transaction, minValids int, wHeaders []chain.NodeHeader, poolR PoolRequester, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, sharedKeyReader shared.KeyReader, nodeReader NodeReader) error {
	fmt.Printf("transaction %x is in progress\n", tx.TransactionHash())

	if !tx.Address().IsValid() {
		return errors.New("invalid transaction address")
	}

	lastVPool, err := findLastValidationPool(tx.Address(), tx.TransactionType(), poolR, nodeReader)
	if err != nil {
		return err
	}

	sPool, err := FindStoragePool(tx.Address(), nodeReader)
	if err != nil {
		return err
	}

	if err := poolR.RequestTransactionTimeLock(sPool, tx.TransactionHash(), tx.Address(), nodePub); err != nil {
		return fmt.Errorf("transaction lock failed: %s", err.Error())
	}

	fmt.Printf("transaction %x is locked\n", tx.TransactionHash())

	go func() {
		vPool, err := FindValidationPool(tx)
		if err != nil {
			fmt.Printf("transaction find validation pool failed: %s\n", err.Error())
			return
		}

		masterValid, err := preValidateTransaction(tx, wHeaders, sPool, vPool, lastVPool, minValids, nodePub, nodePv, sharedKeyReader)
		if err != nil {
			fmt.Printf("transaction pre-validation failed: %s\n", err.Error())
			return
		}
		confirmValids, err := requestValidations(tx, masterValid, vPool, minValids, poolR)
		if err != nil {
			fmt.Printf("transaction validation confirmations failed: %s\n", err.Error())
		}
		if err := tx.Mined(masterValid, confirmValids); err != nil {
			fmt.Printf("transaction mining is invalid: %s\n", err.Error())
		}
		fmt.Printf("transaction %x is validated \n", tx.TransactionHash())
		if err := storeTransaction(tx, sPool, poolR); err != nil {
			fmt.Printf("transaction storage failed: %s\n", err.Error())
			return
		}
	}()

	return nil
}

func storeTransaction(tx chain.Transaction, sPool Pool, poolR PoolRequester) error {
	minReplicas := GetMinimumReplicas(tx.TransactionHash().Digest())
	if err := poolR.RequestTransactionStorage(sPool, minReplicas, tx); err != nil {
		return fmt.Errorf("transaction storage failed: %s", err.Error())
	}
	fmt.Printf("transaction %x is stored \n", tx.TransactionHash())
	return nil
}

//preValidateTransaction checks the incoming transaction as master node by ensure the transaction integrity and perform the proof of work. A valiation will result from this action
func preValidateTransaction(tx chain.Transaction, wHeaders []chain.NodeHeader, sPool Pool, vPool Pool, lastVPool Pool, minValids int, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, sharedKeyReader shared.KeyReader) (chain.MasterValidation, error) {
	if _, err := tx.IsValid(); err != nil {
		return chain.MasterValidation{}, err
	}

	pow, err := proofOfWork(tx, sharedKeyReader)
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

func proofOfWork(tx chain.Transaction, sharedKeyReader shared.KeyReader) (pow crypto.PublicKey, err error) {
	emKeys, err := sharedKeyReader.EmitterCrossKeypairs()
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

func findLastValidationPool(txAddr crypto.VersionnedHash, txType chain.TransactionType, req PoolRequester, nodeReader NodeReader) (Pool, error) {

	sPool, err := FindStoragePool(txAddr, nodeReader)
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

	if !IsValidationConsensusReach(validations) {
		return nil, errors.New("invalid transaction")
	}

	return validations, nil
}

//ConfirmTransactionValidation approve the transaction validation by master and ensure its integrity
func ConfirmTransactionValidation(tx chain.Transaction, masterV chain.MasterValidation, pub crypto.PublicKey, pv crypto.PrivateKey) (chain.Validation, error) {

	var status chain.ValidationStatus

	if _, err := tx.IsValid(); err != nil {
		fmt.Printf("Transaction validation confirmation failed: %s\n", err.Error())
		status = chain.ValidationKO
	} else if _, err := masterV.IsValid(); err != nil {
		fmt.Printf("Transaction master validation confirmation failed: %s\n", err.Error())
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
func GetMinimumValidation(txHash crypto.VersionnedHash) int {
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
