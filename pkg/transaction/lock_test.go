package transaction

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create lock with an empty transaction hash
	Given an empty transaction hash
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithEmptyTxHash(t *testing.T) {
	_, err := NewLock("", "", "")
	assert.EqualError(t, err, "lock: hash is empty")
}

/*
Scenario: Create lock with a transaction hash not in hexadecimal
	Given a transaction with no hexadecimal
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithNotHexTxHash(t *testing.T) {
	_, err := NewLock("hello", "", "")
	assert.EqualError(t, err, "lock: hash is not in hexadecimal format")
}

/*
Scenario: Create lock with a transaction hash which is not a valid hash
	Given a not valid transaction hash
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithNotValidHash(t *testing.T) {
	_, err := NewLock(hex.EncodeToString([]byte("hello")), "", "")
	assert.EqualError(t, err, "lock: hash is not valid")
}

/*
Scenario: Create lock with an empty address
	Given an empty address
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithEmptyAddress(t *testing.T) {

	txHash := crypto.HashString("hello")

	_, err := NewLock(txHash, "", "")
	assert.EqualError(t, err, "lock: hash is empty")
}

/*
Scenario: Create lock with an address in not hexadecimal format
	Given an address with no hexadecimal
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithNotHexAddress(t *testing.T) {

	txHash := crypto.HashString("hello")

	_, err := NewLock(txHash, "hello", "")
	assert.EqualError(t, err, "lock: hash is not in hexadecimal format")
}

/*
Scenario: Create lock with an address which is not a valid hash
	Given an invalid address
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithInvalidAddress(t *testing.T) {

	txHash := crypto.HashString("hello")

	_, err := NewLock(txHash, hex.EncodeToString([]byte("hello")), "")
	assert.EqualError(t, err, "lock: hash is not valid")
}

/*
Scenario: Create lock with an empty master public key
	Given an empty master public key
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithEmptyMasterKey(t *testing.T) {

	txHash := crypto.HashString("hello")
	addr := crypto.HashString("addr")

	_, err := NewLock(txHash, addr, "")
	assert.EqualError(t, err, "lock: public key is empty")
}

/*
Scenario: Create lock with a master public key not in hexadecimal
	Given a master public key not in hexadecimal
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithNotHexMasterKey(t *testing.T) {

	txHash := crypto.HashString("hello")
	addr := crypto.HashString("addr")

	_, err := NewLock(txHash, addr, "key")
	assert.EqualError(t, err, "lock: public key is not in hexadecimal format")
}

/*
Scenario: Create lock with an invalid master public key
	Given a master public key not valid
	When I want to create a new lock
	Then I get an error
*/
func TestNewLockWithInvalidMasterKey(t *testing.T) {

	txHash := crypto.HashString("hello")
	addr := crypto.HashString("addr")

	_, err := NewLock(txHash, addr, hex.EncodeToString([]byte("pubkey")))
	assert.EqualError(t, err, "lock: public key is not valid")
}

/*
Scenario: Create loc
	Given a transaction hash, an address and a master public key
	When I want to create a new lock
	Then I get no error and I can retrieve the lock information
*/
func TestNewLock(t *testing.T) {
	txHash := crypto.HashString("hello")
	addr := crypto.HashString("addr")

	pub, _ := crypto.GenerateKeys()

	lock, err := NewLock(txHash, addr, pub)
	assert.Nil(t, err)
	assert.Equal(t, txHash, lock.TransactionHash())
	assert.Equal(t, addr, lock.Address())
	assert.Equal(t, pub, lock.MasterRobotKey())
}
