package memstorage

import "github.com/uniris/uniris-core/pkg/transaction"

//LockDatabase is the transaction memory database
type LockDatabase interface {
	transaction.LockRepository
}

type lockDb struct {
	locks []transaction.Lock
}

//NewLockDatabase creates a new memory locks database
func NewLockDatabase() LockDatabase {
	return &lockDb{}
}

func (d *lockDb) StoreLock(l transaction.Lock) error {
	d.locks = append(d.locks, l)
	return nil
}
func (d *lockDb) RemoveLock(l transaction.Lock) error {
	pos := d.findLockPosition(l)
	if pos > -1 {
		d.locks = append(d.locks[:pos], d.locks[pos+1:]...)
	}
	return nil
}
func (d *lockDb) ContainsLock(l transaction.Lock) (bool, error) {
	return d.findLockPosition(l) > -1, nil
}

func (d lockDb) findLockPosition(txLock transaction.Lock) int {
	for i, lock := range d.locks {
		if lock.TransactionHash() == txLock.TransactionHash() && txLock.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}
