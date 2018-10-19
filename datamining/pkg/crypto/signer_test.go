package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"

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
	encData := "uxazexc"

	sig, err := Sign(hex.EncodeToString(pvKey), encData)
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)
	var signature ecdsaSignature
	decodesig, _ := hex.DecodeString(string(sig))
	asn1.Unmarshal(decodesig, &signature)

	hash := []byte(HashString(encData))

	assert.True(t, ecdsa.Verify(&key.PublicKey, hash, signature.R, signature.S))
}

/*
Scenario: Verify encrypted data
	Given a data , a signature and a public key
	When I want verify this signature
	Then I get the supposed result
*/
func TestVerify(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	puKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	encData := fakeData{Message: "hello"}

	b, _ := json.Marshal(encData)

	s := NewSigner()
	sig, _ := Sign(hex.EncodeToString(pvKey), string(b))

	err := s.CheckSignature(
		hex.EncodeToString(puKey),
		encData,
		sig,
	)
	assert.Nil(t, err)
}

type fakeData struct {
	Message string
}
