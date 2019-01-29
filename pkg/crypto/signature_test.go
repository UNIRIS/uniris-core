package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Sign encrypted data
	Given an encrypted data and a private key
	When I want sign this data
	Then I get the signature and can be verify by the public key associated
*/
func TestSign(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	encData := "hello"

	sig, err := Sign(encData, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)
	var signature ecdsaSignature
	decodesig, _ := hex.DecodeString(string(sig))
	asn1.Unmarshal(decodesig, &signature)

	assert.True(t, ecdsa.Verify(&key.PublicKey, []byte(HashString(encData)), signature.R, signature.S))
}

/*
Scenario: Verify signature
	Given a data , a signature and a public key
	When I want verify this signature
	Then I get the supposed result
*/
func TestVerify(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	encData := "hello"

	sig, err := Sign(encData, hex.EncodeToString(pvKey))
	assert.Nil(t, err)
	assert.Nil(t, VerifySignature(encData, hex.EncodeToString(pubKey), sig))
}

/*
Scenario: Verify bad signature
	Given a data , a bad signature and a public key
	When I want verify this signature
	Then I get an error
*/
func TestVerifyBadSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	encData := "hello"
	assert.Equal(t, ErrInvalidSignature, VerifySignature(encData, hex.EncodeToString(pubKey), hex.EncodeToString([]byte("fake sig"))))
}

/*
Scenario: Check is signature with empty string
	Given an empty signature
	When I want to check it's a signature
	Then I get an error
*/
func TestIsSignatureWithEmpty(t *testing.T) {
	ok, err := IsSignature("")
	assert.False(t, ok)
	assert.EqualError(t, err, "signature is empty")
}

/*
Scenario: Check is signature with not hexadecimal string
	Given a signature with non hexa format
	When I want to check it's a signature
	Then I get an error
*/
func TestIsSignatureWithNotHexadecimal(t *testing.T) {
	ok, err := IsSignature("hello")
	assert.False(t, ok)
	assert.EqualError(t, err, "signature is not in hexadecimal format")
}

/*
Scenario: Check Is signature when the signature is invalid
	Given a signature badly created
	When I want to check it's a signature
	Then I get an error
*/
func TestIsSignatureInvalid(t *testing.T) {
	ok, err := IsSignature(hex.EncodeToString([]byte("signature")))
	assert.False(t, ok)
	assert.EqualError(t, err, "signature is not valid")
}

/*
Scenario: Check Is signature when the signature is valid
	Given a valid signature
	When I want to check it's a signature
	Then I get no error
*/
func TestIsSignatureValid(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pv, _ := x509.MarshalECPrivateKey(key)
	sig, _ := Sign("hello", hex.EncodeToString(pv))

	ok, err := IsSignature(sig)
	assert.True(t, ok)
	assert.Nil(t, err)
}
