package lock

import (
	"errors"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//Repository defines methods to manage locks
type Repository interface {
	NewLock(datamining.TransactionLock) error
	RemoveLock(datamining.TransactionLock) error
	ContainsLock(datamining.TransactionLock) bool
}

//Service defines methods to handle lock and unlock transactions
type Service interface {
	LockTransaction(datamining.TransactionLock) error
	UnlockTransaction(datamining.TransactionLock) error
}

type service struct {
	repo Repository
}

//NewService creates a new locking service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) LockTransaction(txLock datamining.TransactionLock) error {
	if s.repo.ContainsLock(txLock) {
		return ErrLockExisting
	}

	return s.repo.NewLock(txLock)
}

func (s service) UnlockTransaction(txLock datamining.TransactionLock) error {
	return s.repo.RemoveLock(txLock)
}
