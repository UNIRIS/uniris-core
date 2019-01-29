package transaction

import "errors"

//LockService handles the lock management
type LockService struct {
	repo LockRepository
}

//NewLockService creates a new transaction lock service
func NewLockService(r LockRepository) LockService {
	return LockService{
		repo: r,
	}
}

//StoreLock performs a lock on a transaction
//
//If a lock exist, ErrLockExisting error is returned
func (s LockService) StoreLock(l Lock) error {
	exist, err := s.repo.ContainsLock(l)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("a lock already exist for this transaction")
	}
	return s.repo.StoreLock(l)
}

//ContainsLock checks if a transaction lock exist
func (s LockService) ContainsLock(l Lock) (bool, error) {
	ok, err := s.repo.ContainsLock(l)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//RemoveLock checks if the the lock exists; if so, it's deleted
func (s LockService) RemoveLock(l Lock) error {
	exist, err := s.repo.ContainsLock(l)
	if err != nil {
		return err
	}
	if exist {
		return s.repo.RemoveLock(l)
	}
	return nil
}
