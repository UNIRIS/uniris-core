package memstorage

import (
	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
)

//LockDatabase is the transaction memory database
type LockDatabase interface {
	adding.LockRepository
	listing.LockRepository
}

type lockDb struct {
	locks []uniris.Lock
}

//NewLockDatabase creates a new memory locks database
func NewLockDatabase() LockDatabase {
	return &lockDb{}
}

func (d *lockDb) StoreLock(l uniris.Lock) error {
	d.locks = append(d.locks, l)
	return nil
}
func (d *lockDb) RemoveLock(l uniris.Lock) error {
	pos := d.findLockPosition(l)
	if pos > -1 {
		d.locks = append(d.locks[:pos], d.locks[pos+1:]...)
	}
	return nil
}
func (d *lockDb) ContainsLock(l uniris.Lock) (bool, error) {
	return d.findLockPosition(l) > -1, nil
}

func (d lockDb) findLockPosition(txLock uniris.Lock) int {
	for i, lock := range d.locks {
		if lock.TransactionHash() == txLock.TransactionHash() && txLock.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}
