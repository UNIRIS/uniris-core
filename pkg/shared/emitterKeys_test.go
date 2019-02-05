package shared

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
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create new emitter shared keypair
	Given a encrypted private key and a public key
	When I want to create a emitter shared key
	Then I get the shared emitter keypair without error
*/
func TestNewEmitterSharedKey(t *testing.T) {

	pub, _ := crypto.GenerateKeys()

	kp, err := NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString([]byte("pvKey")), kp.EncryptedPrivateKey())
	assert.Equal(t, pub, kp.PublicKey())
}

/*
Scenario: Create empty shared keypair
	Given a no encrypted private key or no a public key
	When I want to create a shared key
	Then I get an error
*/
func TestNewEmptyEmitterSharedKey(t *testing.T) {
	_, err := NewEmitterKeyPair("", "")
	assert.EqualError(t, err, "shared emitter keys encrypted private key: is empty")
}

/*
Scenario: Create shared emitter keypair with not invalid public key
	Given public key is not a valid public key
	When I want to create a emitter shared key
	Then I get an error
*/
func TestNewSharedEmitterKeyWithInvalidPublicKey(t *testing.T) {
	_, err := NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), "pubKey")
	assert.EqualError(t, err, "shared emitter keys: public key is not in hexadecimal format")

	_, err = NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString([]byte("pubKey")))
	assert.EqualError(t, err, "shared emitter keys: public key is not valid")
}

/*
Scenario: Create shared emitter keypair with not hexadecimal encrypted private key
	Given private key not in hex format
	When I want to create a emitter shared key
	Then I get an error
*/
func TestNewSharedEmitterKeyWithNotHexPvKey(t *testing.T) {
	_, err := NewEmitterKeyPair("pvKey", "pvkey")
	assert.EqualError(t, err, "shared emitter keys encrypted private key: is not in hexadecimal format")
}

/*
Scenario: Marshal into a JSON a shared emitter key pair
	Given a shared emitter keypair
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalSharedEmitterKeys(t *testing.T) {

	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key))
	b, err := json.Marshal(kp)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"encrypted_private_key\":\"%s\",\"public_key\":\"%s\"}", hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key)), string(b))
}
