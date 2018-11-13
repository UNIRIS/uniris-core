package crypto

import (
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

/*
Scenario: Hash a string
	Given a string
	When I want hash it
	Then I get a hash and I can retrieve the same hash by giving the same input
*/
func TestHashString(t *testing.T) {
	hash := hashString("hello")
	assert.NotEmpty(t, hash)

	hash2 := hashString("hello")
	assert.Equal(t, hash, hash2)
}

/*
Scenario: Hash a byte slice
	Given a byte slice
	When I want hash it
	Then I get a hash and I can retrieve the same hash by giving the same input
*/
func TestHashBytes(t *testing.T) {
	hash := hashBytes([]byte("hello"))
	assert.NotEmpty(t, hash)

	hash2 := hashBytes([]byte("hello"))
	assert.Equal(t, hash, hash2)
}

/*
Scenario: Hash a biometric data
	Given biometric data
	When I want to hash it, I create a JSON of it
	Then it produces a hash
*/
func TestHashBiometricData(t *testing.T) {
	bio := account.NewBiometricData("hash", "addr", "addr", "aesKey", "pub", "pub", account.NewSignatures("sig", "sig"))
	hash, err := NewHasher().NewBiometricDataHash(bio)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash a keychain data
	Given a keychain data
	When I want to hash it, I create a JSON of it
	Then it produces a hash
*/
func TestHashKeychainData(t *testing.T) {
	kc := account.NewKeychainData("addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))
	hash, err := NewHasher().NewKeychainDataHash(kc)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash a lock
	Given a lock for a transaction
	When I want to hash it, I create JSON of it
	Then it produces a hash
*/
func TestHashLock(t *testing.T) {
	txLock := lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}

	hash, err := NewHasher().HashLock(txLock)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash a biometric
	Given a biometric transaction
	When I want to hash it, I create JSON of it
	Then it produces a hash
*/
func TestHashBiometric(t *testing.T) {
	data := account.NewBiometricData("hash", "addr", "addr", "aesKey", "pub", "pub", account.NewSignatures("sig", "sig"))
	end := mining.NewEndorsement("last hash", "hash",
		mining.NewMasterValidation([]string{"pubkey"}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")},
	)

	bio := account.NewBiometric(data, end)
	hash, err := NewHasher().HashBiometric(bio)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash a keychain
	Given a keychain transaction
	When I want to hash it, I create JSON of it
	Then it produces a hash
*/
func TestHashKeychain(t *testing.T) {
	data := account.NewKeychainData("addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))
	end := mining.NewEndorsement("last hash", "hash",
		mining.NewMasterValidation([]string{"pubkey"}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")},
	)

	bio := account.NewKeychain("address", data, end)
	hash, err := NewHasher().HashKeychain(bio)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}
