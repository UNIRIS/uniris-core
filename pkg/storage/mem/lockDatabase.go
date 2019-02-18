package memstorage

import (
	"github.com/uniris/uniris-core/pkg/chain"
)

type locker struct {
	locks []map[string]string
}

//NewLocker creates a transaction locker in memory
func NewLocker() chain.Locker {
	return &locker{}
}

func (l *locker) WriteLock(txHash string, txAddr string, masterPubk string) error {
	l.locks = append(l.locks, map[string]string{
		"transaction_address": txAddr,
		"transaction_hash":    txHash,
		"master_public_key":   masterPubk,
	})
	return nil
}
func (l *locker) RemoveLock(txHash string, txAddr string) error {
	pos := l.findLockPosition(txHash, txAddr)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l *locker) ContainsLock(txHash string, txAddr string) (bool, error) {
	return l.findLockPosition(txHash, txAddr) > -1, nil
}

func (l locker) findLockPosition(txHash string, txAddr string) int {
	for i, lock := range l.locks {
		if lock["transaction_hash"] == txHash && lock["transaction_address"] == txAddr {
			return i
		}
	}
	return -1
}
