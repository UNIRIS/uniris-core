package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"time"
)

type poolRequester interface {
	RequestTransactionTimeLock(sPool interface{}, addr []byte, publicKey interface{}) error
	RequestTransactionStorage(sPool interface{}, tx interface{}, minReplicas int) error
	RequestTransactionValidations(vPool interface{}, tx interface{}, minValids int) ([]interface{}, error)
}

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Marshal() []byte
}

type privateKey interface {
	Sign(data []byte) ([]byte, error)
}

type transaction interface {
	Address() []byte
	Type() int
	Data() map[string]interface{}
	Timestamp() time.Time
	PreviousPublicKey() interface{}
	Signature() []byte
	OriginSignature() []byte
	CoordinatorStamp() interface{}
	CrossValidations() []interface{}
	MarshalBeforeOriginSignature() ([]byte, error)
	MarshalRoot() ([]byte, error)
}

//CoordinateTransactionProcessing handles the transaction processing by the coordinator node
//The workflow includes:
// - TimeLock the transaction
// - Create coordinator validation stamp
// - Executes the proof of work
// - Requests cross validation stamps
// - Requests storage
func CoordinateTransactionProcessing(tx interface{}, nbValidations int, coordList interface{}, nodePv interface{}, nodePub interface{}, nodeReader interface{}, originPublicKeys []interface{}, poolReq interface{}) error {

	t, ok := tx.(transaction)
	if !ok {
		return errors.New("mining: transaction type is invalid")
	}

	nPub, ok := nodePub.(publicKey)
	if !ok {
		return errors.New("mining: node public key is invalid")
	}
	nPv, ok := nodePv.(privateKey)
	if !ok {
		return errors.New("mining: node private key is invalid")
	}
	poolR, ok := poolReq.(poolRequester)
	if !ok {
		return errors.New("mining: pool requester is invalid")
	}

	poolPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "poolElection/plugin.so"))
	if err != nil {
		return fmt.Errorf("mining: %s", err.Error())
	}

	storPSym, err := poolPlugin.Lookup("FindStoragePool")
	if err != nil {
		return fmt.Errorf("mining: %s", err.Error())
	}
	storePoolF := storPSym.(func([]byte, interface{}, interface{}, interface{}) (interface{}, error))
	sPool, err := storePoolF(t.Address(), nodePv, nodePub, nodeReader)
	if err != nil {
		return fmt.Errorf("mining: %s", err.Error())
	}

	if err := poolR.RequestTransactionTimeLock(sPool, t.Address(), nPub); err != nil {
		return fmt.Errorf("mining: transaction lock failed: %s", err.Error())
	}

	go coordTxAsync(t, nbValidations, coordList, nPv, nPub, nodeReader, originPublicKeys, poolR, sPool)

	return nil
}

func coordTxAsync(t transaction, nbValidations int, coordList interface{}, nodePv privateKey, nodePub publicKey, nodeReader interface{}, originPublicKeys []interface{}, poolR poolRequester, sPool interface{}) {
	poolPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "poolElection/plugin.so"))
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}
	validPSym, err := poolPlugin.Lookup("FindValidationPool")
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}
	validPoolF := validPSym.(func([]byte, interface{}, interface{}, interface{}) (interface{}, error))
	vPool, err := validPoolF(t.Address(), nodePv, nodePub, nodeReader)
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	coordStamp, err := coordinatorValidation(t, nodePub, nodePv, originPublicKeys, coordList, vPool, sPool)
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	tPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}
	tSym, err := tPlugin.Lookup("NewTransaction")
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}
	fNewTx := tSym.(func([]byte, int, map[string]interface{}, time.Time, interface{}, []byte, []byte, interface{}, []interface{}) (interface{}, error))

	minedTx, err := fNewTx(t.Address(), t.Type(), t.Data(), t.Timestamp(), t.PreviousPublicKey(), t.Signature(), t.OriginSignature(), coordStamp, nil)
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	crossValidations, err := poolR.RequestTransactionValidations(vPool, minedTx, nbValidations)
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	finalTx, err := fNewTx(t.Address(), t.Type(), t.Data(), t.Timestamp(), t.PreviousPublicKey(), t.Signature(), t.OriginSignature(), coordStamp, crossValidations)
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	rPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "replication/plugin.so"))
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	minRSym, err := rPlugin.Lookup("GetMinimumReplicas")
	if err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}
	nbReplicas := minRSym.(func(interface{}) int)(t)

	if err := poolR.RequestTransactionStorage(sPool, finalTx, nbReplicas); err != nil {
		fmt.Printf("mining: %s\n", err.Error())
		return
	}

	return
}

func coordinatorValidation(t transaction, nodePub publicKey, nodePv privateKey, originKeys []interface{}, coordList interface{}, vPool interface{}, sPool interface{}) (interface{}, error) {

	var stamp interface{}

	/*
	* Checks the validity of the transaction
	* If the transaction is invalid no return an error
	* but mark the validation stamp as not valid
	 */
	tPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	tSym, err := tPlugin.Lookup("IsTransactionValid")
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	isTxValid := tSym.(func(tx interface{}) (bool, string))
	if ok, reason := isTxValid(t); !ok {
		fmt.Printf("mining: %s\n", reason)
		vStamp, err := createValidationStamp(0, time.Now(), nodePub, nodePv)
		if err != nil {
			return nil, fmt.Errorf("mining: %s", reason)
		}
		stamp = vStamp
	}

	/*
	* Performs the POW to find the key of the origin signature
	* If the POW no return an error
	* but mark the validation stamp as not valid
	 */
	pow, err := performPOW(t, originKeys)
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	if pow == nil {
		vStamp, err := createValidationStamp(0, time.Now(), nodePub, nodePv)
		if err != nil {
			return nil, fmt.Errorf("mining: %s", err.Error())
		}
		stamp = vStamp
	}

	/*
	* Create the coordinator stamp
	* by loading the Coordinator plugin
	* and the method NewCoordinatorStamp
	 */
	coordPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "coordinatorStamp/plugin.so"))
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	coordNewSym, err := coordPlugin.Lookup("NewCoordinatorStamp")
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	f := coordNewSym.(func(prevCrossV [][]byte, pow interface{}, validStamp interface{}, txHash []byte, elecCoordNodes interface{}, elecCrossVNodes interface{}, elecStorNodes interface{}) (interface{}, error))

	hPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.go"))
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}

	HSym, err := hPlugin.Lookup("Hash")
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}

	tJSON, err := t.MarshalRoot()
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	txHash := HSym.(func([]byte) []byte)(tJSON)

	//TODO: load previous validation pool
	cs, err := f(nil, pow, stamp, txHash, coordList, vPool, sPool)
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	return cs, nil
}

//CrossValidateTransaction checks the validity of the transaction and returns a validation stamp
func CrossValidateTransaction(tx interface{}, nodePub interface{}, nodePv interface{}) (interface{}, error) {

	pubKey, ok := nodePub.(publicKey)
	if !ok {
		return nil, errors.New("mining: node public key type is invalid")
	}

	pvKey, ok := nodePv.(privateKey)
	if !ok {
		return nil, errors.New("mining: node private key type is invalid")
	}

	var status = 1

	/*
	* Check the validatity of the transaction
	* by calling the plugin Transaction and the method IsTransactionValid
	 */
	tPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	tSym, err := tPlugin.Lookup("IsTransactionValid")
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	isTxValid := tSym.(func(tx interface{}) (bool, string))
	if ok, reason := isTxValid(tx); !ok {
		log.Printf("Transaction validation confirmation failed: " + reason)
		status = 0
	}

	stamp, err := createValidationStamp(status, time.Now(), pubKey, pvKey)
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	return stamp, nil
}

//performPOW retrieves the public key using for the origin signature of the transaction
func performPOW(tx interface{}, allOriginKeys []interface{}) (pow interface{}, err error) {

	t, ok := tx.(transaction)
	if !ok {
		return nil, errors.New("mining: transaction type is not valid")
	}

	txBytes, err := t.MarshalBeforeOriginSignature()
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}

	for _, k := range allOriginKeys {
		pubKey, ok := k.(publicKey)
		if !ok {
			return nil, errors.New("origin public key is invalid")
		}

		if ok, err := pubKey.Verify(txBytes, t.OriginSignature()); err != nil {
			return nil, err
		} else if ok {
			return pubKey, nil
		}
	}

	return nil, nil
}

//Creating the validation stamp (even if the validation is KO)
//using the Validation plugin and the method NewValidationStamp
//and signing it with the local node private key
func createValidationStamp(status int, timestamp time.Time, pub publicKey, pv privateKey) (interface{}, error) {

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "validationStamp/plugin.so"))
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	vSym, err := p.Lookup("NewValidationStamp")
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	f := vSym.(func(status int, t time.Time, nodePubk interface{}, nodeSig []byte) (interface{}, error))

	v, err := json.Marshal(map[string]interface{}{
		"status":     status,
		"public_key": pub.Marshal(),
		"timestamp":  time.Now().Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}
	sig, err := pv.Sign(v)
	if err != nil {
		return nil, fmt.Errorf("mining: %s", err.Error())
	}

	return f(status, time.Now(), pub, sig)
}
