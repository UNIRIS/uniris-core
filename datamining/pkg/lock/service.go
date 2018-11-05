package lock

import (
	"errors"
)

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//Repository defines methods to manage locks
type Repository interface {
	NewLock(TransactionLock) error
	RemoveLock(TransactionLock) error
	ContainsLock(TransactionLock) bool
}

//Service defines methods to handle lock and unlock transactions
type Service interface {
	LockTransaction(TransactionLock) error
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
