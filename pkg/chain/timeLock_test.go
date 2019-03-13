package chain

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create transaction timelock
	Given a transaction hash, an address and a master node key
	When I want to timelock the transaction
	Then the timelock is stored
*/
func TestStoreTimeLock(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	assert.Nil(t, TimeLockTransaction(crypto.Hash([]byte("hash1")), crypto.Hash([]byte("addr1")), pub))
	assert.Len(t, timeLockers, 1)
	assert.Equal(t, crypto.Hash([]byte("hash1")), timeLockers[0].txHash)
	assert.Equal(t, crypto.Hash([]byte("addr1")), timeLockers[0].txAddress)
}

/*
Scenario: Create two timelock identicals
	Given a timelock created
	When a want to created again the same timelock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	assert.Nil(t, TimeLockTransaction(crypto.Hash([]byte("hash2")), crypto.Hash([]byte("addr2")), pub))
	assert.EqualError(t, TimeLockTransaction(crypto.Hash([]byte("hash2")), crypto.Hash([]byte("addr2")), pub), "a lock already exist for this transaction")
}

/*
Scenario: Remove a timelock
	Given a timelock created
	When I want to remove it
	Then I get no timelock after
*/
func TestRemoveLock(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	TimeLockTransaction(crypto.Hash([]byte("hash3")), crypto.Hash([]byte("addr3")), pub)
	removeTimeLock(crypto.Hash([]byte("hash3")), crypto.Hash([]byte("addr3")))

	_, found, _ := findTimelock(crypto.Hash([]byte("hash3")), crypto.Hash([]byte("addr3")))
	assert.False(t, found)
}

/*
Scenario: Remove a timelock after countdown
	Given a timelock created
	When the countdown is reached
	Then the timelock is removed
*/
func TestRemoveLockAfterCountdown(t *testing.T) {
	timeLockCountdown = 1 * time.Second
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	TimeLockTransaction(crypto.Hash([]byte("hash4")), crypto.Hash([]byte("addr4")), pub)
	time.Sleep(2 * time.Second)
	_, found, _ := findTimelock(crypto.Hash([]byte("hash4")), crypto.Hash([]byte("addr4")))
	assert.False(t, found)
}
