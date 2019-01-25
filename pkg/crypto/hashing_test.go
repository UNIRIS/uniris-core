package crypto

import (
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
