package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

//Transaction describe a root tx
type Transaction interface {
	Address() []byte
	Type() int
	Data() map[string]interface{}
	Timestamp() time.Time
	PreviousPublicKey() interface{}
	Signature() []byte
	OriginSignature() []byte
	CoordinatorStamp() interface{}
	CrossValidations() []interface{}
	IsCoordinatorStampValid() (bool, string)
}

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Marshal() []byte
}

type coordinatorStamp interface {
	ProofOfWork() interface{}
	ValidationStamp() validationStamp
	TransactionHash() []byte
	IsValid() (bool, string)
}

type validationStamp interface {
	IsValid() (bool, string)
}

type tx struct {
	addr          []byte
	txType        int
	data          map[string]interface{}
	timestamp     time.Time
	pubKey        publicKey
	sig           []byte
	originSig     []byte
	coordStmp     coordinatorStamp
	confirmValids []validationStamp
}

//NewTransaction creates a new transaction and checks its integrity as well as the validation stamps
func NewTransaction(addr []byte, txType int, data map[string]interface{}, timestamp time.Time, pubK interface{}, sig []byte, originSig []byte, coordS interface{}, crossV []interface{}) (interface{}, error) {

	if len(addr) == 0 {
		return nil, errors.New("transaction: address is missing")
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.so"))
	if err != nil {
		return nil, fmt.Errorf("transaction: %s", err.Error())
	}
	sym, err := p.Lookup("IsValidHash")
	if err != nil {
		return nil, fmt.Errorf("transaction: %s", err.Error())
	}
	if !sym.(func([]byte) bool)(addr) {
		return nil, errors.New("transaction: address is an invalid hash")
	}

	if pubK == nil {
		return nil, errors.New("transaction: public key is missing")
	}
	pubKey, ok := pubK.(publicKey)
	if !ok {
		return nil, errors.New("transaction: public key type is invalid")
	}

	if len(sig) == 0 {
		return nil, errors.New("transaction: signature is missing")
	}

	if len(originSig) == 0 {
		return nil, errors.New("transaction: origin signature is missing")
	}

	if len(data) == 0 {
		return nil, errors.New("transaction: data is missing")
	}

	if timestamp.Unix() > time.Now().Unix() {
		return nil, errors.New("transaction: invalid timestamp")
	}

	switch txType {
	case KeychainTransactionType:
	case IDTransactionType:
	case ContractTransactionType:
	case SystemTransactionType:
	default:
		return nil, errors.New("transaction: invalid type")
	}

	t := tx{
		addr:      addr,
		txType:    txType,
		data:      data,
		timestamp: timestamp,
		pubKey:    pubKey,
		sig:       sig,
		originSig: originSig,
	}

	tJSON, err := t.marshalBeforeSignature()
	if err != nil {
		return nil, fmt.Errorf("transaction: %s", err.Error())
	}
	if ok, err := pubKey.Verify(tJSON, t.sig); err != nil {
		return nil, fmt.Errorf("transaction: %s", err.Error())
	} else if !ok {
		return nil, errors.New("transaction: invalid signature")
	}

	//Check the coordinator stamp
	if coordS != nil {
		coordStmp, ok := coordS.(coordinatorStamp)
		log.Print(coordS)
		if !ok {
			return nil, errors.New("transaction: coordinator stamp type is invalid")
		}
		t.coordStmp = coordStmp
		if ok, reason := t.IsCoordinatorStampValid(); !ok {
			return nil, fmt.Errorf("transaction: %s", reason)
		}

	}

	//Check the cross validation stamps
	if crossV != nil && len(crossV) > 0 {
		cv := make([]validationStamp, 0)
		for _, v := range crossV {
			vstmp, ok := v.(validationStamp)
			if !ok {
				return nil, errors.New("transaction: cross validation type is invalid")
			}
			if ok, reason := vstmp.IsValid(); !ok {
				return nil, fmt.Errorf("transaction: %s", reason)
			}
			cv = append(cv, vstmp)
		}
		t.confirmValids = cv
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

func (t tx) IsCoordinatorStampValid() (bool, string) {

	if ok, reason := t.coordStmp.IsValid(); !ok {
		return false, reason
	}

	tJSON, err := t.marshal()
	if err != nil {
		return false, err.Error()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.so"))
	if err != nil {
		return false, err.Error()
	}

	hashSym, err := p.Lookup("Hash")
	if err != nil {
		return false, err.Error()
	}

	f := hashSym.(func([]byte) []byte)
	txHash := f(tJSON)

	if !bytes.Equal(txHash, t.coordStmp.TransactionHash()) {
		return false, "integrity hash is invalid"
	}

	tJSON, err = t.marshalBeforeOriginSignature()
	if err != nil {
		return false, err.Error()
	}

	pow := t.coordStmp.ProofOfWork().(publicKey)
	if ok, err := pow.Verify(tJSON, t.originSig); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "proof of work is invalid"
	}

	return true, ""
}

func (t tx) marshalBeforeSignature() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":       t.addr,
		"data":       t.data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": t.pubKey.Marshal(),
	})
}

func (t tx) marshalBeforeOriginSignature() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":       t.addr,
		"data":       t.data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": t.pubKey.Marshal(),
		"signature":  t.sig,
	})
}

func (t tx) marshal() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":             t.addr,
		"data":             t.data,
		"timestamp":        t.timestamp.Unix(),
		"type":             t.txType,
		"public_key":       t.pubKey.Marshal(),
		"signature":        t.sig,
		"origin_signature": t.originSig,
	})
}
