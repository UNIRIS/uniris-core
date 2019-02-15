package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create transaction lock
	Given a transaction hash, an address and a master miner key
	When I want to lock the transaction
	Then the lock is stored
*/
func TestStoreLock(t *testing.T) {
	lockDB := &mockDB{}

	pub, _ := crypto.GenerateKeys()

	assert.Nil(t, LockTransaction(lockDB, crypto.HashString("hash"), crypto.HashString("addr"), pub))
	assert.Len(t, lockDB.locks, 1)
	assert.Equal(t, crypto.HashString("hash"), lockDB.locks[0]["transaction_hash"])
	assert.Equal(t, crypto.HashString("addr"), lockDB.locks[0]["transaction_address"])
}

/*
Scenario: Create two lock identicals
	Given a lock created
	When a want to created again the same lock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	lockDB := &mockDB{}
	pub, _ := crypto.GenerateKeys()
	assert.Nil(t, LockTransaction(lockDB, crypto.HashString("hash"), crypto.HashString("addr"), pub))
	assert.EqualError(t, LockTransaction(lockDB, crypto.HashString("hash"), crypto.HashString("addr"), pub), "a lock already exist for this transaction")
}

/*
Scenario: Remove a lock
	Given a lock created
	When I want to remove it
	Then I get no lock after
*/
func TestRemoveLock(t *testing.T) {
	lockDB := &mockDB{}

	pub, _ := crypto.GenerateKeys()
	LockTransaction(lockDB, crypto.HashString("hash"), crypto.HashString("addr"), pub)
	assert.Nil(t, UnlockTransaction(lockDB, crypto.HashString("hash"), crypto.HashString("addr")))
	assert.Len(t, lockDB.locks, 0)
}

type mockDB struct {
	locks []map[string]interface{}
}

func (r *mockDB) WriteLock(txHash string, txAddr string, masterPubk string) error {
	r.locks = append(r.locks, map[string]interface{}{
		"transaction_hash":    txHash,
		"transaction_address": txAddr,
		"master_public_key":   masterPubk,
	})
	return nil
}

func (r *mockDB) RemoveLock(txHash string, txAddr string) error {
	pos := r.findLockPosition(txHash, txAddr)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}

func (r mockDB) ContainsLock(txHash string, txAddr string) (bool, error) {
	return r.findLockPosition(txHash, txAddr) > -1, nil
}

func (r mockDB) findLockPosition(txHash string, txAddr string) int {
	for i, l := range r.locks {
		if l["transaction_hash"] == txHash && l["transaction_address"] == txAddr {
			return i
		}
	}
	return -1
}
