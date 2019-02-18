package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	ecies "github.com/uniris/ecies/pkg"
)

/*
Scenario: Decrypt data
	Given encrypted bytes
	When I want decrypt it with ECIES
	Then I get the clear data
*/
func TestDecrypt(t *testing.T) {
	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, err := x509.ParsePKIXPublicKey(pbKey)
	eciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))
	cipher, err := ecies.Encrypt(rand.Reader, eciesKey, []byte("uniris"), nil, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipher)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	clear, _ := Decrypt(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Equal(t, "uniris", string(clear))
}
