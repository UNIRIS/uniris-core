package memstorage

import (
	"github.com/uniris/uniris-core/pkg/consensus"
)

type lockDB struct {
	locks []map[string]string
}

//NewLockDatabase creates a lock database in memory
func NewLockDatabase() consensus.LockDatabase {
	return &lockDB{}
}

func (l *lockDB) WriteLock(txHash string, txAddr string) error {
	l.locks = append(l.locks, map[string]string{
		"transaction_address": txAddr,
		"transaction_hash":    txHash,
	})
	return nil
}
func (l *lockDB) RemoveLock(txHash string, txAddr string) error {
	pos := l.findLockPosition(txHash, txAddr)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l *lockDB) ContainsLock(txHash string, txAddr string) (bool, error) {
	return l.findLockPosition(txHash, txAddr) > -1, nil
}

func (l lockDB) findLockPosition(txHash string, txAddr string) int {
	for i, lock := range l.locks {
		if lock["transaction_hash"] == txHash && lock["transaction_address"] == txAddr {
			return i
		}
	}
	return -1
}
