package uniris

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create a new miner validation
	Given a public key, a status, a timestamp and signature
	When I want to create a miner validation
	Then I get the validation
*/
func TestNewMinerValidation(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pv, _ := x509.MarshalECPrivateKey(key)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	v, err := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)
	assert.Nil(t, err)
	assert.Equal(t, ValidationOK, v.Status())
	assert.Equal(t, time.Now().Unix(), v.Timestamp().Unix())
	assert.Equal(t, hex.EncodeToString(pub), v.MinerPublicKey())
	assert.Equal(t, sig, v.MinerSignature())
}

/*
Scenario: Create a new miner validation with a timestamp later than now
	Given a public key, a status and a timestamp (now + 2 sec)
	When I want to create a miner validation
	Then I get an error
*/
func TestNewMinerValidationWithInvalidTimestamp(t *testing.T) {
	_, err := NewMinerValidation(ValidationOK, time.Now().Add(2*time.Second), "", "")
	assert.EqualError(t, err, "Miner validation: timestamp must be lower than now")
}

/*
Scenario: Create a new miner validation with invalid public key
	Given no public key or no hex or not valid public key
	When I want to create a miner validation
	Then I get an error
*/
func TestNewMinerValidationWithInvalidPublicKey(t *testing.T) {
	_, err := NewMinerValidation(ValidationOK, time.Now(), "", "sig")
	assert.EqualError(t, err, "Miner validation: Public key is empty")

	_, err = NewMinerValidation(ValidationOK, time.Now(), "key", "sig")
	assert.EqualError(t, err, "Miner validation: Public key is not in hexadecimal format")

	_, err = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString([]byte("key")), "sig")
	assert.EqualError(t, err, "Miner validation: Public key is not valid")
}

/*
Scenario: Create a new miner validation with invalid signature
	Given no hex or not valid signature
	When I want to create a miner validation
	Then I get an error
*/
func TestNewMinerValidationWithInvalidSignature(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	_, err := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), "sig")
	assert.EqualError(t, err, "Miner validation: Signature is not in hexadecimal format")

	_, err = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), hex.EncodeToString([]byte("sig")))
	assert.EqualError(t, err, "Miner validation: Signature is not valid")
}

/*
Scenario: Create a new miner validation with an invalid status
	Given public key, signature, timestamp and an invalid validation status
	When I want to create a miner validation
	Then I get an error
*/
func TestNewMinerValidationWithInvalidStatus(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	_, err := NewMinerValidation(10, time.Now(), hex.EncodeToString(pub), sig)
	assert.EqualError(t, err, "Miner validation: status not allowed")
}

/*
Scenario: Check a miner validation
	Given a miner validation
	When I want to validate
	Then I get an error if it's not valid
*/
func TestCheckMinerValidation(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("fake validation marshalling", hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)
	assert.EqualError(t, v.CheckValidation(), "Miner validation: signature is invalid")

	v, _ = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), "")
	b, _ := json.Marshal(v)
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)
	assert.Nil(t, v.CheckValidation())
}

/*
Scenario: Create a new master validation
	Given a proof of work and miner validation
	When I want to create the master validation
	Then I get it
*/
func TestNewMasterValidation(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pv, _ := x509.MarshalECPrivateKey(key)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	sig, _ := crypto.Sign("hello", hex.EncodeToString(pv))

	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)
	mv, err := NewMasterValidation([]PeerIdentity{}, hex.EncodeToString(pub), v)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString(pub), mv.ProofOfWork())
	assert.Equal(t, v.MinerPublicKey(), mv.Validation().MinerPublicKey())
	assert.Equal(t, v.Timestamp(), mv.Validation().Timestamp())
	assert.Empty(t, mv.PreviousTransactionMiners())
}

/*
Scenario: Create a master validation with POW invalid
	Given a no POW or not hex or invalid public key
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithInvalidPOW(t *testing.T) {
	_, err := NewMasterValidation([]PeerIdentity{}, "", MinerValidation{})
	assert.EqualError(t, err, "Master validation Proof of work: Public key is empty")

	_, err = NewMasterValidation([]PeerIdentity{}, "key", MinerValidation{})
	assert.EqualError(t, err, "Master validation Proof of work: Public key is not in hexadecimal format")

	_, err = NewMasterValidation([]PeerIdentity{}, hex.EncodeToString([]byte("key")), MinerValidation{})
	assert.EqualError(t, err, "Master validation Proof of work: Public key is not valid")
}

/*
Scenario: Create a master validation without miner validation
	Given a no validation
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithoutValidation(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	_, err := NewMasterValidation([]PeerIdentity{}, hex.EncodeToString(pub), MinerValidation{})
	assert.EqualError(t, err, "Master validation: Missing pre-validation")
}
