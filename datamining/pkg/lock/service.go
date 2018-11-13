package lock

import (
	"errors"
)

//TransactionLock represents lock data
type TransactionLock struct {
	TxHash         string
	MasterRobotKey string
	Address        string
}

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//Hasher defines methods to handle lock hashing
type Hasher interface {

	//HashLock produces a hash of the lock transaction
	HashLock(txLock TransactionLock) (string, error)
}

//Repository defines methods to manage locks
type Repository interface {

	//NewLock stores a lock
	NewLock(TransactionLock) error

	//RemoveLock remove an existing lock
	RemoveLock(TransactionLock) error

	//ContainsLocks determines if a lock exists or not
	ContainsLock(TransactionLock) bool
}

//Service defines methods to handle lock and unlock transactions
type Service interface {
	//LockTransaction performs a lock on a transaction
	//
	//If a lock exist, ErrLockExisting error is returned
	LockTransaction(TransactionLock) error

	//UnlockTransaction performs a unlock on a transaction
	UnlockTransaction(TransactionLock) error
}

type service struct {
	repo Repository
}

//NewService creates a new locking service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) LockTransaction(txLock TransactionLock) error {
	if s.repo.ContainsLock(txLock) {
		return ErrLockExisting
	}

	return s.repo.NewLock(txLock)
}

func (s service) UnlockTransaction(txLock TransactionLock) error {
	return s.repo.RemoveLock(txLock)
}
