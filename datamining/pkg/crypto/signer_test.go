package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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
	s, _ := NewSigner()
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	encData := []byte("uxazexc")
	h := sha256.Sum256(encData)
	d := h[:]
	sig, err := s.Sign(pvKey, d)
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)
	var signature ecdsaSignature
	decodesig, _ := hex.DecodeString(string(sig))
	asn1.Unmarshal(decodesig, &signature)
	assert.True(t, ecdsa.Verify(&key.PublicKey, d, signature.R, signature.S))
}

/*
Scenario: Verify encrypted data
	Given a data , a signature and a public key
	When I want verify this signature
	Then I get the supposed result
*/

func TestVerify(t *testing.T) {
	s, _ := NewSigner()
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	puKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	encData := []byte("uxazexc")
	h := sha256.Sum256(encData)
	d := h[:]
	R, S, _ := ecdsa.Sign(rand.Reader, key, d)
	sig, _ := asn1.Marshal(ecdsaSignature{R, S})
	err := s.Verify(puKey, []byte(hex.EncodeToString(sig)), d)
	assert.Nil(t, err)
}
