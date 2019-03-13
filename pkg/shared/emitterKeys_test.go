package shared

import (
	"crypto/rand"
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

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp, err := NewEmitterKeyPair([]byte("pvKey"), pub)
	assert.Nil(t, err)
	assert.Equal(t, []byte("pvKey"), kp.EncryptedPrivateKey())
	assert.Equal(t, pub, kp.PublicKey())
}

/*
Scenario: Create empty shared keypair
	Given a no encrypted private key or no a public key
	When I want to create a shared key
	Then I get an error
*/
func TestNewEmptyEmitterSharedKey(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, err := NewEmitterKeyPair([]byte(""), pub)
	assert.EqualError(t, err, "shared emitter keys encrypted private key: is empty")
}
