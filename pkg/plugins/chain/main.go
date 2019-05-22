package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"time"
)

type chainWriter interface {
	WriteKeychain(tx interface{}) error
	WriteID(tx interface{}) error
	WriteKO(tx interface{}) error
}

type chainReader interface {
	FindKeychainByAddr(addr []byte) (interface{}, error)
	FindKeychainByHash(txHash []byte) (interface{}, error)
	FindIDByHash(txHash []byte) (interface{}, error)
	FindIDByAddr(addr []byte) (interface{}, error)
	FindKOByHash(txHash []byte) (interface{}, error)
	FindKOByAddr(addr []byte) (interface{}, error)
}

type indexReader interface {
	FindLastTransactionAddr(genesis []byte) ([]byte, error)
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
}

//StoreTransaction persist the transaction in the database and checks it before storing it
func StoreTransaction(tx interface{}, minValidations int, w interface{}) error {

	t, ok := tx.(transaction)
	if !ok {
		return errors.New("invalid transaction type")
	}

	chainW, ok := w.(chainWriter)
	if !ok {
		return errors.New("invalid chain writer type")
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return err
	}

	sym, err := p.Lookup("IsTransactionValid")
	if err != nil {
		return err
	}

	if ok, reason := sym.(func(tx interface{}) (bool, string))(tx); !ok {
		fmt.Println(reason)
		if err := chainW.WriteKO(t); err != nil {
			return err
		}
	}

	//TODO: check the elected nodes

	switch t.Type() {
	case 0:
		return chainW.WriteKeychain(t)
	case 1:
		return chainW.WriteID(t)
	default:
		return errors.New("invalid transaction type")
	}
}

func GetLastTransaction(genesis []byte, chainR interface{}, indexR interface{}) (interface{}, error) {
	iR, ok := indexR.(indexReader)
	if !ok {
		return nil, errors.New("invalid index reader")
	}

	cR, ok := chainR.(chainReader)
	if !ok {
		return nil, errors.New("invalid chain reader")
	}

	addr, err := iR.FindLastTransactionAddr(genesis)
	if err != nil {
		return nil, err
	}

	tx, err := cR.FindKeychainByAddr(addr)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	tx, err = cR.FindIDByAddr(addr)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	tx, err = cR.FindKOByAddr(addr)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	return nil, errors.New("invalid transaction")
}

func GetTransactionByHash(txHash []byte, chainR interface{}) (interface{}, error) {

	cR, ok := chainR.(chainReader)
	if !ok {
		return nil, errors.New("invalid chain reader")
	}

	tx, err := cR.FindKeychainByHash(txHash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	tx, err = cR.FindIDByHash(txHash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	tx, err = cR.FindKOByHash(txHash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx.(transaction), nil
	}

	return nil, errors.New("invalid transaction")
}

func GetTransactionStatus(txHash []byte, chainR interface{}) (int, error) {
	//TODO: check if timelocked

	cR, ok := chainR.(chainReader)
	if !ok {
		return -1, errors.New("invalid chain reader")
	}

	t, err := cR.FindKOByHash(txHash)
	if err != nil {
		return -1, err
	}
	if t != nil {
		return 0, nil
	}

	t, err = cR.FindKeychainByHash(txHash)
	if err != nil {
		return -1, err
	}
	if t != nil {
		return 0, nil
	}

	return -1, errors.New("invalid transaction")
}
