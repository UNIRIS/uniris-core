package locking

import "errors"

//ErrLockExisting is returned when a lock already exist
var ErrLockExisting = errors.New("A lock already exist for this transaction")

//Repository defines methods to manage locks
type Repository interface {

	//NewLock stores a lock
	NewLock(txHash string, addr string) error

	//RemoveLock remove an existing lock
	RemoveLock(txHash string, addr string) error

	//ContainsLocks determines if a lock exists or not
	ContainsLock(txHash string, addr string) (bool, error)
}

//Service handles transaction locks
type Service struct {
	repo Repository
}

//NewService creates a new locking service
func NewService(repo Repository) Service {
	return Service{repo: repo}
}

//LockTransaction performs a lock on a transaction
//
//If a lock exist, ErrLockExisting error is returned
func (s Service) LockTransaction(txHash string, addr string) error {
	exist, err := s.ContainsTransactionLock(txHash, addr)
	if err != nil {
		return err
	}
	if exist {
		return ErrLockExisting
	}
	s.repo.NewLock(txHash, addr)
	return nil
}

//ContainsTransactionLock determines if a transaction lock exists
func (s Service) ContainsTransactionLock(txHash string, addr string) (bool, error) {
	ok, err := s.repo.ContainsLock(txHash, addr)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//UnlockTransaction performs a unlock on a transaction
func (s Service) UnlockTransaction(txHash string, addr string) error {
	exist, err := s.repo.ContainsLock(txHash, addr)
	if err != nil {
		return err
	}
	if exist {
		return s.repo.RemoveLock(txHash, addr)
	}
	return nil
}
