package mock

import "github.com/uniris/uniris-core/datamining/pkg/mining/lock"

type locker struct {
	locks []lock.TransactionLock
}

//NewTransactionLocker creates a transaction locker
func NewTransactionLocker() lock.TransactionLocker {
	return locker{}
}

func (l locker) Lock(txLock lock.TransactionLock) error {
	l.locks = append(l.locks, txLock)
	return nil
}

func (l locker) Unlock(txLock lock.TransactionLock) error {
	pos := l.findLockPosition(txLock)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}

func (l locker) ContainsLock(txLock lock.TransactionLock) bool {
	for _, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (l locker) findLockPosition(txLock lock.TransactionLock) int {
	for i, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
