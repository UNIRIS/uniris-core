package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Check a public key is valid with empty string
	Given an empty string
	When I want to check if it's a public key
	Then I get an error
*/
func TestIsPublicKeyWithEmpty(t *testing.T) {
	ok, err := IsPublicKey("")
	assert.False(t, ok)
	assert.EqualError(t, err, "public key is empty")
}

/*
Scenario: Check a public key is valid with not hexadecimal
	Given a string on non -hexadecimal format
	When I want to check if it's a public key
	Then I get an error
*/
func TestIsPublicKeyWithNotHexa(t *testing.T) {
	ok, err := IsPublicKey("hello")
	assert.False(t, ok)
	assert.EqualError(t, err, "public key is not in hexadecimal format")
}

/*
Scenario: Check a public key is valid with is was badly created
	Given an invalid public key
	When I want to check if it's a public key
	Then I get an error
*/
func TestIsPublicKeyInvalid(t *testing.T) {
	ok, err := IsPublicKey(hex.EncodeToString([]byte("pubKey")))
	assert.False(t, ok)
	assert.EqualError(t, err, "public key is not valid")
}

/*
Scenario: Check a public key is valid with another generation algorithm (i.e. RSA)
	Given an RSA public key
	When I want to check if it's a public key on ECDSA algorithm
	Then I get an error
*/
func TestIsPublicKeyNotElliptic(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	ok, err := IsPublicKey(hex.EncodeToString(pub))
	assert.False(t, ok)
	assert.EqualError(t, err, "public key is not from an elliptic curve")
}

/*
Scenario: Check a public key is valid with is was well created
	Given an valid public key
	When I want to check if it's a public key
	Then I get no error
*/
func TestIsPublicKeyValid(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	ok, err := IsPublicKey(hex.EncodeToString(pub))
	assert.True(t, ok)
	assert.Nil(t, err)

}

/*
Scenario: Check a private key is valid with empty string
	Given an empty string
	When I want to check if it's a private key
	Then I get an error
*/
func TestIsPrivateKeyWithEmpty(t *testing.T) {
	ok, err := IsPrivateKey("")
	assert.False(t, ok)
	assert.EqualError(t, err, "private key is empty")
}

/*
Scenario: Check a private key is valid with not hexadecimal
	Given a string on non -hexadecimal format
	When I want to check if it's a private key
	Then I get an error
*/
func TestIsPrivateKeyWithNotHexa(t *testing.T) {
	ok, err := IsPrivateKey("hello")
	assert.False(t, ok)
	assert.EqualError(t, err, "private key is not in hexadecimal format")
}

/*
Scenario: Check a private key is valid with is was badly created
	Given an invalid public key
	When I want to check if it's a private key
	Then I get an error
*/
func TestIsPrivateKeyInvalid(t *testing.T) {
	ok, err := IsPrivateKey(hex.EncodeToString([]byte("pvKey")))
	assert.False(t, ok)
	assert.EqualError(t, err, "private key is not valid")
}

/*
Scenario: Check a public key is valid with is was well created
	Given an valid public key
	When I want to check if it's a public key
	Then I get no error
*/
func TestIsPrivateKeyValid(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pv, _ := x509.MarshalECPrivateKey(key)

	ok, err := IsPrivateKey(hex.EncodeToString(pv))
	assert.True(t, ok)
	assert.Nil(t, err)

}
