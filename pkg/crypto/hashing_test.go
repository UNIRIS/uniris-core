package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Hash a string
	Given a string
	When I want hash it
	Then I get a hash and I can retrieve the same hash by giving the same input
*/
func TestHashString(t *testing.T) {
	hash := HashString("hello")
	assert.NotEmpty(t, hash)

	hash2 := HashString("hello")
	assert.Equal(t, hash, hash2)
}

/*
Scenario: Hash a byte slice
	Given a byte slice
	When I want hash it
	Then I get a hash and I can retrieve the same hash by giving the same input
*/
func TestHashBytes(t *testing.T) {
	hash := HashBytes([]byte("hello"))
	assert.NotEmpty(t, hash)

	hash2 := HashBytes([]byte("hello"))
	assert.Equal(t, hash, hash2)
}

/*
Scenario: Check the hash with empty string
	Given an empty hash
	When I want to check if it's an hash
	Then I get an error
*/
func TestIsHashWhenEmpty(t *testing.T) {
	ok, err := IsHash("")
	assert.False(t, ok)
	assert.EqualError(t, err, "hash is empty")
}

/*
Scenario: Check the hash with non hexadecimal
	Given an hash on non hexadicimal
	When I want to check if it's an hash
	Then I get an error
*/
func TestIsHashWhenNotHexadecimal(t *testing.T) {
	ok, err := IsHash("hello")
	assert.False(t, ok)
	assert.EqualError(t, err, "hash is not in hexadecimal format")
}

/*
Scenario: Check the hash with invalid format
	Given an hash without the required size
	When I want to check if it's an hash
	Then I get an error
*/
func TestIsHashNotValid(t *testing.T) {
	ok, err := IsHash(hex.EncodeToString([]byte("my hash")))
	assert.False(t, ok)
	assert.EqualError(t, err, "hash is not valid")
}

/*
Scenario: Check the hash
	Given an hash valid
	When I want to check if it's an hash
	Then I get no error
*/
func TestIsHashValid(t *testing.T) {
	ok, err := IsHash(HashString("hello"))
	assert.True(t, ok)
	assert.Nil(t, err)
}
