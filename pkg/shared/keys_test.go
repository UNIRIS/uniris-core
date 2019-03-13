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
Scenario: Create new EmitterCross keypair
	Given a EmitterCross private key and a public key
	When I want to create a EmitterCross key
	Then I get the EmitterCross keypair without error
*/
func TestNewEmitterCrossKey(t *testing.T) {

	pub, _ := crypto.GenerateKeys()

	kp, err := NewEmitterCrossKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString([]byte("pvKey")), kp.EncryptedPrivateKey())
	assert.Equal(t, pub, kp.PublicKey())
}

/*
Scenario: Create empty EmitterCross keypair
	Given a no EmitterCross private key or no a public key
	When I want to create the EmitterCross keypair
	Then I get an error
*/
func TestNewEmptyEmitterCrossKey(t *testing.T) {
	_, err := NewEmitterCrossKeyPair("", "")
	assert.EqualError(t, err, "missing emitter cross private key")
}

/*
Scenario: Create an EmitterCross keypair with not invalid public key
	Given public key is not a valid public key
	When I want to create a, EmitterCross  key
	Then I get an error
*/
func TestNewEmitterCrossKeyWithInvalidPublicKey(t *testing.T) {
	_, err := NewEmitterCrossKeyPair(hex.EncodeToString([]byte("pvKey")), "pubKey")
	assert.EqualError(t, err, "emitter cross: public key is not in hexadecimal format")

	_, err = NewEmitterCrossKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString([]byte("pubKey")))
	assert.EqualError(t, err, "emitter cross: public key is not valid")
}

/*
Scenario: Create EmitterCross keypair with not hexadecimal EmitterCross private key
	Given private key not in hex format
	When I want to create a EmitterCross  key
	Then I get an error
*/
func TestNewEmitterCrossKeyWithNotHexPvKey(t *testing.T) {
	_, err := NewEmitterCrossKeyPair("pvKey", "pvkey")
	assert.EqualError(t, err, "emitter cross private key is not in hexadecimal format")
}

/*
Scenario: Marshal into a JSON an EmitterCross key pair
	Given an EmitterCross keypair
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalEmitterCrossKeys(t *testing.T) {

	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := NewEmitterCrossKeyPair(hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key))
	b, err := json.Marshal(kp)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"encrypted_private_key\":\"%s\",\"public_key\":\"%s\"}", hex.EncodeToString([]byte("pvKey")), hex.EncodeToString(key)), string(b))
}
