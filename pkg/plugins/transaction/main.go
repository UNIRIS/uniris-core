package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"time"
)

var (
	//KeychainTransactionType represents a Transaction related to keychain
	KeychainTransactionType = 0

	//IDTransactionType represents a Transaction related to ID data
	IDTransactionType = 1

	//ContractTransactionType represents a Transaction related to a smart contract
	ContractTransactionType = 2

	//ContractResponseTransactionType represents a Transaction related to a smart contract response
	ContractResponseTransactionType = 3

	//SystemTransactionType represents a transaction related to the network/infrastructure
	SystemTransactionType = 4
)

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

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Marshal() []byte
}

type coordinatorStamp interface {
	ProofOfWork() interface{}
	TransactionHash() []byte
}

type tx struct {
	addr          []byte
	txType        int
	data          map[string]interface{}
	timestamp     time.Time
	pubKey        interface{}
	sig           []byte
	originSig     []byte
	coordStmp     interface{}
	confirmValids []interface{}
}

//NewTransaction creates a new transaction and checks its integrity as well as the validation stamps
func NewTransaction(addr []byte, txType int, data map[string]interface{}, timestamp time.Time, pubK interface{}, sig []byte, originSig []byte, coordS interface{}, crossV []interface{}) (interface{}, error) {

	t := tx{
		addr:          addr,
		txType:        txType,
		data:          data,
		timestamp:     timestamp,
		pubKey:        pubK,
		sig:           sig,
		originSig:     originSig,
		coordStmp:     coordS,
		confirmValids: crossV,
	}

	if ok, reason := IsTransactionValid(t); !ok {
		return nil, errors.New(reason)
	}

	return t, nil
}

func (t tx) Address() []byte {
	return t.addr
}

func (t tx) Type() int {
	return t.txType
}

func (t tx) Data() map[string]interface{} {
	return t.data
}

func (t tx) Timestamp() time.Time {
	return t.timestamp
}

func (t tx) PreviousPublicKey() interface{} {
	return t.pubKey
}

func (t tx) Signature() []byte {
	return t.sig
}

func (t tx) OriginSignature() []byte {
	return t.originSig
}

func (t tx) CoordinatorStamp() interface{} {
	return t.coordStmp
}
func (t tx) CrossValidations() []interface{} {
	vv := make([]interface{}, len(t.confirmValids))
	for _, v := range t.confirmValids {
		vv = append(vv, v)
	}
	return vv
}

func (t tx) MarshalBeforeOriginSignature() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":       t.addr,
		"data":       t.data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": t.pubKey.(publicKey).Marshal(),
		"signature":  t.sig,
	})
}

func (t tx) MarshalRoot() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":             t.addr,
		"data":             t.data,
		"timestamp":        t.timestamp.Unix(),
		"type":             t.txType,
		"public_key":       t.pubKey.(publicKey).Marshal(),
		"signature":        t.sig,
		"origin_signature": t.originSig,
	})
}

//IsTransactionValid checks if a transaction is valid
//starting by the root then if coordinator stamp is here it will be validated (calling coordinatorStamp plugin)
//following by the cross validation (calling validationStamp plugin as well)
func IsTransactionValid(tx interface{}) (bool, string) {
	t, ok := tx.(transaction)
	if !ok {
		return false, "transaction: invalid type"
	}

	if ok, reason := checkTransactionRoot(t); !ok {
		return true, fmt.Sprintf("transaction: %s", reason)
	}

	if t.CoordinatorStamp() != nil {
		if ok, reason := checkCoordinatorStamp(t); !ok {
			return false, fmt.Sprintf("transaction: %s", reason)
		}
	}
	if t.CrossValidations() != nil && len(t.CrossValidations()) > 0 {
		if ok, reason := checkCrossValidations(t); !ok {
			return false, fmt.Sprintf("transaction: %s", reason)
		}
	}
	return true, ""
}

func checkTransactionRoot(t transaction) (bool, string) {
	if t.Address() == nil || len(t.Address()) == 0 {
		return false, "transaction: address is missing"
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.so"))
	if err != nil {
		return false, fmt.Sprintf("transaction: %s", err.Error())
	}
	sym, err := p.Lookup("IsValidHash")
	if err != nil {
		return false, fmt.Sprintf("transaction: %s", err.Error())
	}
	if !sym.(func([]byte) bool)(t.Address()) {
		return false, "transaction: address is an invalid hash"
	}

	if t.PreviousPublicKey() == nil {
		return false, "transaction: public key is missing"
	}
	pubKey, ok := t.PreviousPublicKey().(publicKey)
	if !ok {
		return false, "transaction: public key type is invalid"
	}

	if t.Signature() == nil || len(t.Signature()) == 0 {
		return false, "transaction: signature is missing"
	}

	if t.OriginSignature() == nil || len(t.OriginSignature()) == 0 {
		return false, "transaction: origin signature is missing"
	}

	if t.Data() == nil || len(t.Data()) == 0 {
		return false, "transaction: data is missing"
	}

	if t.Timestamp().Unix() > time.Now().Unix() {
		return false, "transaction: invalid timestamp"
	}

	switch t.Type() {
	case KeychainTransactionType:
	case IDTransactionType:
	case ContractTransactionType:
	case SystemTransactionType:
	default:
		return false, "transaction: invalid type"
	}

	tJSON, err := marshalBeforeSignature(t)
	if err != nil {
		return false, fmt.Sprintf("transaction: %s", err.Error())
	}
	if ok, err := pubKey.Verify(tJSON, t.Signature()); err != nil {
		return false, fmt.Sprintf("transaction: %s", err.Error())
	} else if !ok {
		return false, "transaction: invalid signature"
	}

	return true, ""
}

func checkCoordinatorStamp(t transaction) (bool, string) {

	coorPlug, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "coordinatorStamp/plugin.so"))
	if err != nil {
		return false, err.Error()
	}
	isValidCoordSym, err := coorPlug.Lookup("IsCoordinatorStampValid")
	if err != nil {
		return false, err.Error()
	}
	isValidCoordF := isValidCoordSym.(func(interface{}) (bool, string))
	if ok, reason := isValidCoordF(t.CoordinatorStamp()); !ok {
		return false, reason
	}

	tJSON, err := marshalTransaction(t)
	if err != nil {
		return false, err.Error()
	}

	hashPlug, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.so"))
	if err != nil {
		return false, err.Error()
	}

	hashSym, err := hashPlug.Lookup("Hash")
	if err != nil {
		return false, err.Error()
	}

	hashF := hashSym.(func([]byte) []byte)
	txHash := hashF(tJSON)

	if !bytes.Equal(txHash, t.CoordinatorStamp().(coordinatorStamp).TransactionHash()) {
		return false, "integrity hash is invalid"
	}

	tJSON, err = t.MarshalBeforeOriginSignature()
	if err != nil {
		return false, err.Error()
	}

	pow := t.CoordinatorStamp().(coordinatorStamp).ProofOfWork().(publicKey)
	if ok, err := pow.Verify(tJSON, t.OriginSignature()); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "proof of work is invalid"
	}

	return true, ""
}

func checkCrossValidations(t transaction) (bool, string) {

	if t.CrossValidations() != nil && len(t.CrossValidations()) > 0 {
		vPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "validationStamp/plugin.so"))
		if err != nil {
			return false, err.Error()
		}
		validStampSym, err := vPlugin.Lookup("IsValidStamp")
		if err != nil {
			return false, err.Error()
		}
		isValidStamp := validStampSym.(func(interface{}) (bool, string))
		cv := make([]interface{}, 0)
		for _, v := range t.CrossValidations() {
			if ok, reason := isValidStamp(v); !ok {
				return false, reason
			}
			cv = append(cv, v)
		}
	}

	return true, ""
}

func marshalTransaction(t transaction) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":             t.Address(),
		"data":             t.Data(),
		"timestamp":        t.Timestamp().Unix(),
		"type":             t.Type(),
		"public_key":       t.PreviousPublicKey().(publicKey).Marshal(),
		"signature":        t.Signature(),
		"origin_signature": t.OriginSignature(),
	})
}

func marshalBeforeSignature(t transaction) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":       t.Address(),
		"data":       t.Data(),
		"timestamp":  t.Timestamp().Unix(),
		"type":       t.Type(),
		"public_key": t.PreviousPublicKey().(publicKey).Marshal(),
	})
}
