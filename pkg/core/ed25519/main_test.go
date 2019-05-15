package main

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	ed25519 "golang.org/x/crypto/ed25519"
)

/*
Scenario: Generate a new ed25519 key
	Given a seed
	When I want to generate keys and create a key object with the curve identifier in the "key" plugin
	Then I get the same generate keys from the same seed
*/
func TestGenerateKeys(t *testing.T) {

	src1 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	pv, pub, err := GenerateKeys(src1)
	assert.Nil(t, err)
	assert.NotNil(t, pv)
	assert.NotNil(t, pub)

	src2 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	_, pv2, err := ed25519.GenerateKey(src2)
	assert.Nil(t, err)
	assert.EqualValues(t, pv, pv2)
}

/*
Scenario: Sign with an Ed25519 private key
	Given an Ed25519 private key and some data
	When I want to sign this data
	Then I get a signature valid
*/
func TestSign(t *testing.T) {
	pv, pub, _ := GenerateKeys(rand.Reader)
	sig, err := Sign(pv, []byte("hello"))
	assert.Nil(t, err)
	assert.True(t, ed25519.Verify(pub, []byte("hello"), sig))
}

/*
Scenario: Verify a signature with an Ed25519 public key
	Given an ECDSA public key, a signature and the related data
	When I want to verify the signature
	Then I get not error and return true (if the data changed or the signature is invalid I get a false)
*/
func TestVerify(t *testing.T) {
	pv, pub, _ := GenerateKeys(rand.Reader)
	sig, err := Sign(pv, []byte("hello"))
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)

	ok, err := Verify(pub, []byte("hello"), sig)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, _ = Verify(pub, []byte("other data"), sig)
	assert.False(t, ok)

	ok, _ = Verify(pub, []byte("hello"), []byte("fake sig"))
	assert.False(t, ok)
}

/*
Scenario: Generate a shared secret using elliptic curve scalar multiplication
	Given a set of public key and private key
	When I want to create shared secret
	Then I get a shared secret using scalar multiplication of the given public key and private key (and vice-versa).
	The secret is the same by key inversion
*/
func TestGenerateSharedSecret(t *testing.T) {
	pv1, pub1, _ := GenerateKeys(rand.Reader)
	pv2, pub2, _ := GenerateKeys(rand.Reader)

	secret1, err := GenerateSharedSecret(pub1, pv2)
	assert.Nil(t, err)

	secret2, err := GenerateSharedSecret(pub2, pv1)
	assert.Nil(t, err)

	assert.Equal(t, secret1, secret2)
}

/*
Scenario: Extract public key from a cipher message
	Given a cipher message
	When I want to extract the public key inside
	Then I get the public key (91 first bytes)
*/
func TestExtractMessagePublicKey(t *testing.T) {

	_, pub, _ := GenerateKeys(rand.Reader)

	test := make([]byte, 0)
	test = append(test, pub...)
	test = append(test, []byte("blabla")...)

	rPub, pos, err := ExtractMessagePublicKey(test)
	assert.Nil(t, err)
	assert.NotEmpty(t, rPub)
	assert.NotEqual(t, 0, pos)

	assert.Equal(t, pub, rPub)
}
