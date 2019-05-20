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
	WriteTransaction(tx transaction) error
	WriteKOTransaction(tx transaction) error
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
func StoreTransaction(tx interface{}, minValidations int, w chainWriter) error {

	t, ok := tx.(transaction)
	if !ok {
		return errors.New("invalid transaction type")
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
		if err := w.WriteKOTransaction(t); err != nil {
			return err
		}
	}

	//TODO: check the elected nodes

	return w.WriteTransaction(t)
}
