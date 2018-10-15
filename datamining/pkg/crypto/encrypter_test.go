package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/ecies/pkg"
)

/*
Scenario: Encrypt the wallet AES key
	Given the [] bytes
	When I want encrypt it with ECIES
	Then I get the key encrypted and I can decrypt with my ECDSA private key
*/
func TestEncryptext(t *testing.T) {

	e := Encrypter{}
	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())

	cipher, err := e.Encrypt([]byte(hex.EncodeToString(pbKey)), []byte("uniris"))
	assert.Nil(t, err)
	assert.NotEmpty(t, cipher)

	decodeCipher, _ := hex.DecodeString(string(cipher))

	clear, _ := ecies.ImportECDSA(superKey).Decrypt(decodeCipher, nil, nil)
	log.Print(clear)
	assert.Equal(t, []byte("uniris"), clear)
}

/*
Scenario: Encrypt the wallet AES key
	Given the [] bytes
	When I want encrypt it with ECIES
	Then I get the key encrypted and I can decrypt with my ECDSA private key
*/
func TestDecryptText(t *testing.T) {

	e := Encrypter{}

	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, err := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))
	cipher, err := ecies.Encrypt(rand.Reader, robotEciesKey, []byte("uniris"), nil, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipher)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)
	clear, _ := e.Decrypt([]byte(hex.EncodeToString(pvkey)), []byte(hex.EncodeToString(cipher)))
	assert.Equal(t, []byte("uniris"), clear)
}
