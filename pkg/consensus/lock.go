package consensus

import (
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//LockDatabase define methods to handle the lock storage
type LockDatabase interface {

	//WriteLock creates a new lock from the given transaction hash and transaction address
	WriteLock(txHash string, txAddr string, masterPublicKey string) error

	//RemoveLock remove the written lock for the given transaction
	RemoveLock(txHash string, txAddr string) error

	//ContainsLock determinates if a lock for the given transaction
	ContainsLock(txHash string, txAddr string) (bool, error)
}

//LockTransaction stores the lock in the locker system. If a lock exists already an error is returned
func LockTransaction(db LockDatabase, txHash string, txAddr string, masterPubk string) error {
	if _, err := crypto.IsHash(txHash); err != nil {
		return fmt.Errorf("lock transaction hash: %s", err.Error())
	}

	if _, err := crypto.IsHash(txAddr); err != nil {
		return fmt.Errorf("lock transaction address: %s", err.Error())
	}

	if _, err := crypto.IsPublicKey(masterPubk); err != nil {
		return fmt.Errorf("lock transaction public key: %s", err.Error())
	}

	exist, err := db.ContainsLock(txHash, txAddr)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("a lock already exist for this transaction")
	}
	return db.WriteLock(txHash, txAddr, masterPubk)
}

//UnlockTransaction deletes a transaction from the locker
func UnlockTransaction(db LockDatabase, txHash string, txAddr string) error {
	exist, err := db.ContainsLock(txHash, txAddr)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("no lock exist for this transaction")
	}
	return db.RemoveLock(txHash, txAddr)
}
