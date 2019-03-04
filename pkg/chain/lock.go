package chain

import (
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//Locker define methods to handle the lock storage
type Locker interface {

	//WriteLock creates a new lock from the given transaction hash and transaction address
	WriteLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error

	//RemoveLock remove the written lock for the given transaction
	RemoveLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error

	//ContainsLock determinates if a lock for the given transaction
	ContainsLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) (bool, error)
}

//LockTransaction stores the lock in the locker system. If a lock exists already an error is returned
func LockTransaction(l Locker, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPubk crypto.PublicKey) error {
	if !txHash.IsValid() {
		return fmt.Errorf("lock transaction hash: invalid tx hash")
	}

	if !txAddr.IsValid() {
		return fmt.Errorf("lock transaction address: invalid address hash")
	}

	exist, err := l.ContainsLock(txHash, txAddr)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("a lock already exist for this transaction")
	}
	return l.WriteLock(txHash, txAddr, masterPubk)
}

func unlockTransaction(l Locker, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	exist, err := l.ContainsLock(txHash, txAddr)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("no lock exist for this transaction")
	}
	return l.RemoveLock(txHash, txAddr)
}
