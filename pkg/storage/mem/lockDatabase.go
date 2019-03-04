package memstorage

import (
	"bytes"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

type locker struct {
	locks []map[string][]byte
}

//NewLocker creates a transaction locker in memory
func NewLocker() chain.Locker {
	return &locker{}
}

func (l *locker) WriteLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPubk crypto.PublicKey) error {
	mPub, _ := masterPubk.Marshal()
	l.locks = append(l.locks, map[string][]byte{
		"transaction_address": txAddr,
		"transaction_hash":    txHash,
		"master_public_key":   mPub,
	})
	return nil
}
func (l *locker) RemoveLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	pos := l.findLockPosition(txHash, txAddr)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l *locker) ContainsLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) (bool, error) {
	return l.findLockPosition(txHash, txAddr) > -1, nil
}

func (l locker) findLockPosition(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) int {
	for i, lock := range l.locks {
		if bytes.Equal(lock["transaction_hash"], txHash) && bytes.Equal(lock["transaction_address"], txAddr) {
			return i
		}
	}
	return -1
}
