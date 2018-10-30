package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"
)

type locker struct {
	locks []slave.TransactionLock
}

//NewTransactionLocker creates a transaction locker
func NewTransactionLocker() slave.Locker {
	return locker{
		locks: make([]slave.TransactionLock, 0),
	}
}

func (l locker) Lock(txLock slave.TransactionLock) error {
	l.locks = append(l.locks, txLock)
	return nil
}

func (l locker) Unlock(txLock slave.TransactionLock) error {
	pos := l.findLockPosition(txLock)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}

func (l locker) ContainsLock(txLock slave.TransactionLock) bool {
	for _, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (l locker) findLockPosition(txLock slave.TransactionLock) int {
	for i, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
