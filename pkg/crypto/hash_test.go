package crypto

import (
	"crypto"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a new versionned hash
	Given a hash digest
	When I want to create as hash digest
	Then I get a hash preceed by its algorithm
*/
func TestNewVersionnedHash(t *testing.T) {
	h := sha256Func([]byte("hello"))
	vh := NewVersionnedHash(crypto.SHA256, h)
	assert.Equal(t, crypto.SHA256, vh.Algorithm())
	assert.Equal(t, h, vh.Digest())
}

/*
Scenario: Hash a data using default algo
	Given an input data
	When I want hash it
	Then I get a hash and I can retrieve the same hash by giving the same input
*/
func TestHashBytes(t *testing.T) {
	vHash := Hash([]byte("hello"))
	assert.NotEmpty(t, vHash)

	h := DefaultHashAlgo.New()
	h.Write([]byte("hello"))
	hash := h.Sum(nil)
	assert.Equal(t, NewVersionnedHash(DefaultHashAlgo, hash), vHash)
}

/*
Scenario: Check if a hash is valid
	Given a valid hash
	When I want to check if it's valid
	Then I get not error and returns true
*/
func TestIsValidHash(t *testing.T) {
	h := sha256Func([]byte("hello"))
	vh := NewVersionnedHash(crypto.SHA256, h)
	assert.True(t, vh.IsValid())
}

/*
Scenario: Check if a hash is valid when it's empty
	Given an empty hash
	When I want to check if it's valid
	Then I get an error
*/
func TestIsValidHashWithEmptyDigest(t *testing.T) {
	vh := NewVersionnedHash(crypto.SHA256, []byte(""))
	assert.False(t, vh.IsValid())
}

/*
Scenario: Check if a hash is valid with invalid size
	Given an hash with invalid size
	When I want to check if it's valid
	Then I get an error
*/
func TestIsValidHashWithInvalidSize(t *testing.T) {
	vh := NewVersionnedHash(crypto.SHA256, []byte("abc"))
	assert.False(t, vh.IsValid())
}

/*
Scenario: Hash with SHA256
	Given differents input
	When I want to generate a hash with SHA256
	Then I get differents hash
*/
func TestSHA256Hash(t *testing.T) {
	hash := sha256Func([]byte("hello"))
	assert.NotEmpty(t, hash)

	h2 := sha256.New()
	h2.Write([]byte("hello"))
	assert.Equal(t, h2.Sum(nil), hash)

	hash2 := sha256Func([]byte("helloX"))
	assert.NotEqual(t, hash, hash2)
}
