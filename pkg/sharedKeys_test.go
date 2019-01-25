package uniris

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create new shared keypair
	Given a encrypted private key and a public key
	When I want to create a shared key
	Then I get the shared keypair without error
*/
func TestNewSharedKey(t *testing.T) {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	sk, err := NewSharedKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key))
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString([]byte("pvKey")), sk.EncryptedPrivateKey())
	assert.Equal(t, hex.EncodeToString(key), sk.PublicKey())
}

/*
Scenario: Create empty shared keypair
	Given a no encrypted private key or no a public key
	When I want to create a shared key
	Then I get an error
*/
func TestNewEmptySharedKey(t *testing.T) {
	_, err := NewSharedKeyPair("", "")
	assert.EqualError(t, err, "Shared keys encrypted private key: is empty")
}

/*
Scenario: Create shared keypair with not invalid public key
	Given public key is not a valid public key
	When I want to create a shared key
	Then I get an error
*/
func TestNewSharedKeyWithInvalidPublicKey(t *testing.T) {
	_, err := NewSharedKeyPair(hex.EncodeToString([]byte("pvKey")), "pubKey")
	assert.EqualError(t, err, "Shared keys: Public key is not in hexadecimal format")

	_, err = NewSharedKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString([]byte("pubKey")))
	assert.EqualError(t, err, "Shared keys: Public key is not valid")
}

/*
Scenario: Create shared keypair with not hexadecimal encrypted private key
	Given private key not in hex format
	When I want to create a shared key
	Then I get an error
*/
func TestNewSharedKeyWithNotHexPvKey(t *testing.T) {
	_, err := NewSharedKeyPair("pvKey", "pvkey")
	assert.EqualError(t, err, "Shared keys encrypted private key: is not in hexadecimal format")
}

/*
Scenario: Marshal into a JSON a sharedkey pair
	Given a shared keypair
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalSharedKeys(t *testing.T) {

	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := NewSharedKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key))
	b, err := json.Marshal(kp)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"encrypted_private_key\":\"%s\",\"public_key\":\"%s\"}", hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key)), string(b))
}
