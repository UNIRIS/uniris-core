package locking

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

	locker := new(mockLocker)

	service := NewService(locker)
	assert.Nil(t, service.LockTransaction(TransactionLock{
		MasterRobotKey: "robokey",
		TxHash:         "txhash",
		Address:        "address",
	}))

	assert.Len(t, locker.locks, 1)
}

/*
Scenario: Lock a already locked transaction
	Given a lock transaction request
	When there is one lock for this transaction
	Then the an error is returned
*/
func TestCannotAlreadyLockTransaction(t *testing.T) {

	locker := new(mockLocker)

	service := NewService(locker)
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

	assert.Len(t, locker.locks, 1)
}

/*
Scenario: UnLock a locked transaction
	Given a unlock transaction request
	When there is one lock for this transaction
	Then the lock is removed
*/
func TestUnlockTransaction(t *testing.T) {

	locker := new(mockLocker)

	service := NewService(locker)
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

	assert.Len(t, locker.locks, 0)
}

type mockLocker struct {
	locks []TransactionLock
}

func (l *mockLocker) Lock(txLock TransactionLock) error {
	l.locks = append(l.locks, txLock)
	return nil
}
func (l *mockLocker) Unlock(txLock TransactionLock) error {
	pos := l.findLockPosition(txLock)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l mockLocker) ContainsLock(txLock TransactionLock) bool {
	for _, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey && lock.Address == txLock.Address {
			return true
		}
	}
	return false
}

func (l mockLocker) findLockPosition(txLock TransactionLock) int {
	for i, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
