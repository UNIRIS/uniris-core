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

type fakeObj struct {
	message string
}

/*
Scenario: Check the validity of a signature
	Given a public key and a signed data
	When I check the validity of this signature
	Then I get a positive response
*/
func TestCheckSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	data := fakeObj{message: "hi"}
	b, _ := json.Marshal(data)

	r, s, _ := ecdsa.Sign(rand.Reader, key, hash(b))

	sig, _ := asn1.Marshal(ecdsaSignature{r, s})

	val := RequestValidator{}

	valid, err := val.CheckSignature(data, []byte(hex.EncodeToString(pub)), []byte(hex.EncodeToString(sig)))
	assert.Nil(t, err)
	assert.True(t, valid)
}
