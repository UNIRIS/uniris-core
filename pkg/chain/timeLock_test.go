package chain

import (
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
	pub, _ := crypto.GenerateKeys()

	assert.Nil(t, TimeLockTransaction(crypto.HashString("hash1"), crypto.HashString("addr1"), pub))
	assert.Len(t, timeLockers, 1)
	assert.Equal(t, crypto.HashString("hash1"), timeLockers[0].txHash)
	assert.Equal(t, crypto.HashString("addr1"), timeLockers[0].txAddress)
}

/*
Scenario: Create two timelock identicals
	Given a timelock created
	When a want to created again the same timelock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	pub, _ := crypto.GenerateKeys()
	assert.Nil(t, TimeLockTransaction(crypto.HashString("hash2"), crypto.HashString("addr2"), pub))
	assert.EqualError(t, TimeLockTransaction(crypto.HashString("hash2"), crypto.HashString("addr2"), pub), "a lock already exist for this transaction")
}

/*
Scenario: Remove a timelock
	Given a timelock created
	When I want to remove it
	Then I get no timelock after
*/
func TestRemoveLock(t *testing.T) {
	pub, _ := crypto.GenerateKeys()
	TimeLockTransaction(crypto.HashString("hash3"), crypto.HashString("addr3"), pub)
	removeTimeLock(crypto.HashString("hash3"), crypto.HashString("addr3"))

	_, found, _ := findTimelock(crypto.HashString("hash3"), crypto.HashString("addr3"))
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
	pub, _ := crypto.GenerateKeys()
	TimeLockTransaction(crypto.HashString("hash4"), crypto.HashString("addr4"), pub)
	time.Sleep(2 * time.Second)
	_, found, _ := findTimelock(crypto.HashString("hash4"), crypto.HashString("addr4"))
	assert.False(t, found)
}
