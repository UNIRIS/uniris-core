package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/locking"
)

type locker struct {
	locks []locking.TransactionLock
}

//NewTransactionLocker creates a transaction locker
func NewTransactionLocker() locking.Locker {
	return &locker{
		locks: make([]locking.TransactionLock, 0),
	}
}

func (l *locker) Lock(txLock locking.TransactionLock) error {
	l.locks = append(l.locks, txLock)
	return nil
}

func (l *locker) Unlock(txLock locking.TransactionLock) error {
	pos := l.findLockPosition(txLock)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}

func (l locker) ContainsLock(txLock locking.TransactionLock) bool {
	for _, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (l locker) findLockPosition(txLock locking.TransactionLock) int {
	for i, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
