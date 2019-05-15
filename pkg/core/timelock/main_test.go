package main

import (
	"crypto/rand"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create transaction timelock
	Given a transaction hash, an address and a master node key
	When I want to timelock the transaction
	Then the timelock is stored
*/
func TestStoreTimeLock(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, TimeLockTransaction([]byte("hash1"), []byte("addr1"), pub))
	assert.Len(t, timeLockers, 1)
	assert.Equal(t, []byte("hash1"), timeLockers[0].txHash)
	assert.Equal(t, []byte("addr1"), timeLockers[0].txAddress)
}

/*
Scenario: Create two timelock identicals
	Given a timelock created
	When a want to created again the same timelock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, TimeLockTransaction([]byte("hash2"), []byte("addr2"), pub))
	assert.EqualError(t, TimeLockTransaction([]byte("hash2"), []byte("addr2"), pub), "a lock already exist for this transaction")
}

/*
Scenario: Remove a timelock
	Given a timelock created
	When I want to remove it
	Then I get no timelock after
*/
func TestRemoveLock(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	TimeLockTransaction([]byte("hash3"), []byte("addr3"), pub)
	removeTimeLock([]byte("hash3"), []byte("addr3"))

	_, found, _ := findTimelock([]byte("hash3"), []byte("addr3"))
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
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	TimeLockTransaction([]byte("hash4"), []byte("addr4"), pub)
	time.Sleep(2 * time.Second)
	_, found, _ := findTimelock([]byte("hash4"), []byte("addr4"))
	assert.False(t, found)
}
