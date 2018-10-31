package locking

import "errors"

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//TransactionLock represents lock data
type TransactionLock struct {
	TxHash         string
	MasterRobotKey string
	Address        string
}

//Locker defines methods to manage locks
type Locker interface {
	Lock(TransactionLock) error
	Unlock(TransactionLock) error
	ContainsLock(TransactionLock) bool
}

//Service defines methods to handle lock and unlock transactions
type Service interface {
	LockTransaction(TransactionLock) error
	UnlockTransaction(TransactionLock) error
}

type service struct {
	locker Locker
}

//NewService creates a new locking service
func NewService(l Locker) Service {
	return service{l}
}

func (s service) LockTransaction(txLock TransactionLock) error {
	if s.locker.ContainsLock(txLock) {
		return ErrLockExisting
	}

	return s.locker.Lock(txLock)
}

func (s service) UnlockTransaction(txLock TransactionLock) error {
	return s.locker.Unlock(txLock)
}
