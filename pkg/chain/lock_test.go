package chain

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create transaction lock
	Given a transaction hash, an address and a master node key
	When I want to lock the transaction
	Then the lock is stored
*/
func TestStoreLock(t *testing.T) {
	lockDB := &mockLocker{}

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	assert.Nil(t, LockTransaction(lockDB, crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")), pub))
	assert.Len(t, lockDB.locks, 1)
	assert.EqualValues(t, crypto.Hash([]byte("hash")), lockDB.locks[0]["transaction_hash"])
	assert.EqualValues(t, crypto.Hash([]byte("addr")), lockDB.locks[0]["transaction_address"])
}

/*
Scenario: Create two lock identicals
	Given a lock created
	When a want to created again the same lock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	lockDB := &mockLocker{}
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	assert.Nil(t, LockTransaction(lockDB, crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")), pub))
	assert.EqualError(t, LockTransaction(lockDB, crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")), pub), "a lock already exist for this transaction")
}

/*
Scenario: Remove a lock
	Given a lock created
	When I want to remove it
	Then I get no lock after
*/
func TestRemoveLock(t *testing.T) {
	lockDB := &mockLocker{}

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	LockTransaction(lockDB, crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")), pub)
	assert.Nil(t, unlockTransaction(lockDB, crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr"))))
	assert.Len(t, lockDB.locks, 0)
}

type mockLocker struct {
	locks []map[string][]byte
}

func (r *mockLocker) WriteLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPubk crypto.PublicKey) error {
	masterKey, _ := masterPubk.Marshal()

	r.locks = append(r.locks, map[string][]byte{
		"transaction_hash":    txHash,
		"transaction_address": txAddr,
		"master_public_key":   masterKey,
	})
	return nil
}

func (r *mockLocker) RemoveLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	pos := r.findLockPosition(txHash, txAddr)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}

func (r mockLocker) ContainsLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) (bool, error) {
	return r.findLockPosition(txHash, txAddr) > -1, nil
}

func (r mockLocker) findLockPosition(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) int {
	for i, l := range r.locks {
		if bytes.Equal(l["transaction_hash"], txHash) && bytes.Equal(l["transaction_address"], txAddr) {
			return i
		}
	}
	return -1
}
