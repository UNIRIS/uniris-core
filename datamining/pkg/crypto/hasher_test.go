package crypto

import (
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"

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
Scenario: Hash a ID
	Given ID
	When I want to hash it, I create a JSON of it
	Then it produces a hash
*/
func TestHashID(t *testing.T) {
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	id := account.NewID("hash", "addr", "addr", "aesKey", "id pub", "id sig", "em sig", prop)
	hash, err := NewHasher().HashID(id)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash a keychain
	Given a keychain
	When I want to hash it, I create a JSON of it
	Then it produces a hash
*/
func TestHashKeychain(t *testing.T) {
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	kc := account.NewKeychain("addr", "enc wallet", "id pub", "id sig", "em sig", prop)
	hash, err := NewHasher().HashKeychain(kc)
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
Scenario: Hash a endorsed id
	Given an endorsed id
	When I want to hash it, I create JSON of it
	Then it produces a hash
*/
func TestHashEndorsedID(t *testing.T) {
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))
	id := account.NewID("hash", "addr", "addr", "aesKey", "id pub", "id sig", "em sig", prop)
	end := mining.NewEndorsement("last hash", "hash",
		mining.NewMasterValidation([]string{"pubkey"}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")},
	)

	eID := account.NewEndorsedID(id, end)
	hash, err := NewHasher().HashEndorsedID(eID)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

/*
Scenario: Hash an endorsed keychain
	Given an endorsed keychain
	When I want to hash it, I create JSON of it
	Then it produces a hash
*/
func TestHashEndorsedKeychain(t *testing.T) {
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	kc := account.NewKeychain("addr", "enc wallet", "id pub", "id sig", "em sig", prop)
	end := mining.NewEndorsement("last hash", "hash",
		mining.NewMasterValidation([]string{"pubkey"}, "pubkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature")},
	)

	eKC := account.NewEndorsedKeychain("address", kc, end)
	hash, err := NewHasher().HashEndorsedKeychain(eKC)
	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}
