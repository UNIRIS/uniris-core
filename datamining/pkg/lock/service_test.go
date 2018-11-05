package lock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Lock a incoming transaction
	Given a lock transaction request
	When there is no lock for this transaction
	Then the lock is created
*/
func TestLockTransaction(t *testing.T) {

	repo := new(mockRepository)

	service := NewService(repo)
	assert.Nil(t, service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	assert.Len(t, repo.locks, 1)
}

/*
Scenario: Lock a already locked transaction
	Given a lock transaction request
	When there is one lock for this transaction
	Then the an error is returned
*/
func TestCannotAlreadyLockTransaction(t *testing.T) {

	repo := new(mockRepository)

	service := NewService(repo)
	assert.Nil(t, service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	assert.Equal(t, ErrLockExisting, service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	assert.Len(t, repo.locks, 1)
}

/*
Scenario: UnLock a locked transaction
	Given a unlock transaction request
	When there is one lock for this transaction
	Then the lock is removed
*/
func TestUnlockTransaction(t *testing.T) {

	repo := new(mockRepository)

	service := NewService(repo)
	assert.Nil(t, service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	})

	assert.Nil(t, service.UnlockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	assert.Len(t, repo.locks, 0)
}

type mockRepository struct {
	locks []TransactionLock
}

func (r *mockRepository) NewLock(txLock TransactionLock) error {
	r.locks = append(r.locks, txLock)
	return nil
}

func (r *mockRepository) RemoveLock(txLock TransactionLock) error {
	pos := r.findLockPosition(txLock)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}
func (r mockRepository) ContainsLock(txLock TransactionLock) bool {
	for _, lock := range r.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey && lock.Address == txLock.Address {
			return true
		}
	}
	return false
}

func (r mockRepository) findLockPosition(txLock TransactionLock) int {
	for i, lock := range r.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
