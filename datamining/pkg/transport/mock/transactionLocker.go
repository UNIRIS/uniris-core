package mock

import "github.com/uniris/uniris-core/datamining/pkg/validating"

type locker struct {
	locks []validating.TransactionLock
}

//NewTransactionLocker creates a transaction locker
func NewTransactionLocker() validating.TransactionLocker {
	return locker{}
}

func (l locker) Lock(txLock validating.TransactionLock) error {
	l.locks = append(l.locks, txLock)
	return nil
}

func (l locker) Unlock(txLock validating.TransactionLock) error {
	pos := l.findLockPosition(txLock)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}

func (l locker) ContainsLock(txLock validating.TransactionLock) bool {
	for _, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (l locker) findLockPosition(txLock validating.TransactionLock) int {
	for i, lock := range l.locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
