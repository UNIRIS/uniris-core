package main

import (
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a new node ValidationStamp
	Given a public key, a status, a timestamp and signature
	When I want to create a node ValidationStamp
	Then I get the ValidationStamp
*/
func TestNewTransactionValidationStamp(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	pubK := mockPublicKey{bytes: pub, curve: 0}

	b, _ := json.Marshal(vStamp{
		nodePubk:  pubK,
		status:    ValidationOK,
		timestamp: time.Now(),
	})

	pvKey := mockPrivateKey{bytes: pv}
	sig, _ := pvKey.Sign(b)

	v, err := NewValidationStamp(ValidationOK, time.Now(), pubK, sig)
	assert.Nil(t, err)
	assert.Equal(t, ValidationOK, v.(validationStamp).Status())
	assert.Equal(t, time.Now().Unix(), v.(validationStamp).Timestamp().Unix())
	assert.Equal(t, pubK, v.(validationStamp).NodePublicKey())
	assert.Equal(t, sig, v.(validationStamp).NodeSignature())
}

/*
Scenario: Create a new node ValidationStamp with a timestamp later than now
	Given a public key, a status and a timestamp (now + 2 sec)
	When I want to create a node ValidationStamp
	Then I get an error
*/
func TestNewTransactionValidationStampWithInvalidTimestamp(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	_, err := NewValidationStamp(ValidationOK, time.Now().Add(2*time.Second), mockPublicKey{bytes: pub}, []byte("sig"))
	assert.EqualError(t, err, "validation stamp: timestamp must be anterior or equal to now")
}

/*
Scenario: Create a new node ValidationStamp with invalid public key
	Given no public key or no hex or not valid public key
	When I want to create a node ValidationStamp
	Then I get an error
*/
func TestNewTransactionValidationStampWithInvalidPublicKey(t *testing.T) {
	_, err := NewValidationStamp(ValidationOK, time.Now(), nil, []byte("sig"))
	assert.EqualError(t, err, "validation stamp: public key is missing")
}

/*
Scenario: Create a new node ValidationStamp with invalid signature
	Given no hex or not valid signature
	When I want to create a node ValidationStamp
	Then I get an error
*/
func TestNewTransactionValidationStampWithInvalidSignature(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
	pubK := mockPublicKey{bytes: pub}
	pvK := mockPrivateKey{bytes: pv}

	_, err := NewValidationStamp(ValidationOK, time.Now(), pubK, nil)
	assert.EqualError(t, err, "validation stamp: signature is missing")

	_, err = NewValidationStamp(ValidationOK, time.Now(), pubK, []byte("sig"))
	assert.EqualError(t, err, "validation stamp: signature is not valid")

	sig, _ := pvK.Sign([]byte("hello"))
	_, err = NewValidationStamp(ValidationOK, time.Now(), pubK, sig)
	assert.EqualError(t, err, "validation stamp: signature is not valid")
}

/*
Scenario: Create a new node ValidationStamp with an invalid status
	Given public key, signature, timestamp and an invalid ValidationStamp status
	When I want to create a node ValidationStamp
	Then I get an error
*/
func TestNewTransactionValidationStampWithInvalidStatus(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
	pvKey := mockPrivateKey{bytes: pv}
	pubKey := mockPublicKey{bytes: pub}

	sig, _ := pvKey.Sign([]byte("hello"))

	_, err := NewValidationStamp(10, time.Now(), pubKey, sig)
	assert.EqualError(t, err, "validation stamp: invalid status")
}

type mockPublicKey struct {
	bytes []byte
	curve int
}

func (pb mockPublicKey) Marshal() []byte {
	out := make([]byte, 1+len(pb.bytes))
	out[0] = byte(int(pb.curve))
	copy(out[1:], pb.bytes)
	return out
}

func (pb mockPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	return ed25519.Verify(pb.bytes, data, sig), nil
}

type mockPrivateKey struct {
	bytes []byte
}

func (pv mockPrivateKey) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(pv.bytes, data), nil
}
