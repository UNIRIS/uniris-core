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
