package transaction

import (
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

type LockRepository interface {
	StoreLock(l Lock) error
	ContainsLock(l Lock) (bool, error)
	RemoveLock(l Lock) error
}

type Lock struct {
	txHash         string
	address        string
	masterRobotKey string
}

func NewLock(txHash, address, masterKey string) (Lock, error) {

	if _, err := crypto.IsHash(txHash); err != nil {
		return Lock{}, fmt.Errorf("lock: %s", err.Error())
	}

	if _, err := crypto.IsHash(address); err != nil {
		return Lock{}, fmt.Errorf("lock: %s", err.Error())
	}

	if _, err := crypto.IsPublicKey(masterKey); err != nil {
		return Lock{}, fmt.Errorf("lock: %s", err.Error())
	}

	return Lock{txHash, address, masterKey}, nil
}

func (l Lock) TransactionHash() string {
	return l.txHash
}

func (l Lock) Address() string {
	return l.address
}

func (l Lock) MasterRobotKey() string {
	return l.masterRobotKey
}
