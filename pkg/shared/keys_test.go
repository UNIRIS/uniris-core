package shared

import (
	"crypto/rand"
	"encoding/base64"
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

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp, err := NewEmitterCrossKeyPair([]byte("pvKey"), pub)
	assert.Nil(t, err)
	assert.Equal(t, []byte("pvKey"), kp.EncryptedPrivateKey())
	assert.Equal(t, pub, kp.PublicKey())
}

/*
Scenario: Create empty EmitterCross keypair
	Given a no EmitterCross private key or no a public key
	When I want to create the EmitterCross keypair
	Then I get an error
*/
func TestNewEmptyEmitterCrossKey(t *testing.T) {
	_, err := NewEmitterCrossKeyPair([]byte(""), nil)
	assert.EqualError(t, err, "missing emitter cross private key")

	_, err = NewEmitterCrossKeyPair([]byte("hello"), nil)
	assert.EqualError(t, err, "missing emitter cross public key")
}

/*
Scenario: Marshal into a JSON an EmitterCross key pair
	Given an EmitterCross keypair
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalEmitterCrossKeys(t *testing.T) {

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	kp, _ := NewEmitterCrossKeyPair([]byte("pvKey"), pub)
	b, err := json.Marshal(kp)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"encrypted_private_key\":\"%s\",\"public_key\":\"%s\"}", base64.StdEncoding.EncodeToString([]byte("pvKey")), base64.StdEncoding.EncodeToString(pubB)), string(b))
}

/*
Scenario: Create a new cross node key without keys
	Given no public key or no private key
	When I Want to create a cross node key
	Then i Get errors
*/
func TestNewEmptyNodeCrossKey(t *testing.T) {
	_, err := NewNodeCrossKeyPair(nil, nil)
	assert.EqualError(t, err, "missing node cross public key")

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	_, err = NewNodeCrossKeyPair(pub, nil)
	assert.EqualError(t, err, "missing node cross private key")
}

/*
Scenario: Create a new cross node key with public and private keys
	Given a public key and a private
	When I Want to create a cross node key
	Then i Get the cross keypair
*/
func TestNewNodeCrossKey(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp, err := NewNodeCrossKeyPair(pub, pv)
	assert.Nil(t, err)
	assert.Equal(t, pv, kp.PrivateKey())
	assert.Equal(t, pub, kp.PublicKey())
}
