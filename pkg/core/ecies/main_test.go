package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/crypto/ed25519"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	os.Setenv("PLUGINS_DIR", filepath.Join(dir, "../"))
	m.Run()
}

/*
Scenario: Derivate keys from a secret
	Given a secret
	When I want to derivate X keys
	Then I get X derivated keys differents
*/
func TestDerivateKeys(t *testing.T) {
	keys, err := derivateKeys(crypto.SHA256, []byte("hello"), 2)
	assert.Nil(t, err)
	assert.Len(t, keys, 2)
	assert.NotEqual(t, keys[0], keys[1])
}

/*
Scenario: Generate a message authentication code
	Given a key and a message
	When I want generate an authenticate code message using that (i.e. HMAC)
	Then I get the MAC. If two mac generated there are the same
*/
func TestAuthenticateMessage(t *testing.T) {
	tag := authenticateMessage(crypto.SHA256, []byte("my key"), []byte("my data"))
	assert.NotEmpty(t, tag)

	tag2 := authenticateMessage(crypto.SHA256, []byte("my key"), []byte("my data"))
	assert.NotEmpty(t, tag2)

	assert.Equal(t, tag, tag2)
}

/*
Scenario: Encode a encryption
	Given cipher data, a public key, and a tag (MAC)
	When I want to encode the encryption result
	Then I get byte slice withing all these information
*/
func TestNewEncodedCipher(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	pubK := mockKey{
		b: pub,
		c: 1,
	}

	c := newEncodedCipher(pubK, []byte("cipher data"), []byte("tag MAC"))
	assert.NotEmpty(t, c)

	assert.EqualValues(t, pubK.Bytes(), c[:len(pubK.Bytes())])
	assert.EqualValues(t, []byte("cipher data"), c[len(pubK.Bytes()):len(pubK.Bytes())+len([]byte("cipher data"))])
	assert.EqualValues(t, []byte("tag MAC"), c[len(pubK.Bytes())+len([]byte("cipher data")):])
}

/*
Scenario: Decode a encryption with ecdsa
	Given encryption encoded
	When I want to decode the encryption result
	Then I get encrypted data, a public key, and a tag (MAC)
*/
func TestDecodeEncryptionWithEcdsa(t *testing.T) {
	pv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(pv.Public())

	pubK := mockKey{
		b: pub,
		c: 1,
	}

	tag := authenticateMessage(crypto.SHA256, []byte("key"), []byte("message"))
	cipher := newEncodedCipher(pubK, []byte("cipher data"), tag)
	assert.NotEmpty(t, cipher)

	rPub, em, tag, err := decodeCipher(cipher, crypto.SHA256, 1)
	assert.Nil(t, err)
	assert.Equal(t, pub, rPub)
	assert.Equal(t, tag, tag)
	assert.Equal(t, []byte("cipher data"), em)
}

/*
Scenario: Decode a encryption with Ed25519
	Given encryption encoded
	When I want to decode the encryption result
	Then I extract the encrypted data, the public key by loading plugin for its curve, and authentication tag (MAC)
*/
func TestDecodeEncryptionWithEd25519(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	pubK := mockKey{
		b: pub,
		c: 0,
	}

	tag := authenticateMessage(crypto.SHA256, []byte("key"), []byte("message"))
	cipher := newEncodedCipher(pubK, []byte("cipher data"), tag)
	assert.NotEmpty(t, cipher)

	rPub, em, tag, err := decodeCipher(cipher, crypto.SHA256, 0)
	assert.Nil(t, err)
	assert.EqualValues(t, pub, rPub)
	assert.Equal(t, tag, tag)
	assert.Equal(t, []byte("cipher data"), em)
}

/*
Scenario: Encrypt using ecies and ECDSA key
	Given a dest ECDSA public key and a data
	When I want to encrypt the data
	Then I will generate a common secret and get a cipher message
*/
func TestEciesEncryptWithECDSA(t *testing.T) {
	pv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(pv.Public())

	pubK := mockKey{
		b: pub,
		c: 1,
	}

	cipherData, err := Encrypt([]byte("hello"), pubK)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)
}

/*
Scenario: Encrypt using ecies and Ed25519 key
	Given a dest Ed25519 public key and a data
	When I want to encrypt the data
	Then I will generate a common secret and get a cipher message
*/
func TestEciesEncryptWithEd25519(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	pubK := mockKey{
		b: pub,
		c: 0,
	}

	cipherData, err := Encrypt([]byte("hello"), pubK)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)
}

/*
Scenario: Decrypt using ecies and ECDSA key
	Given a cipher data encrypted by an ECDSA key and an ECDSA private key
	When I want to decrypt the data
	Then I will generate a common secret and decrypt the cipher and get a clear data
*/
func TestEciesDecryptWithECDSA(t *testing.T) {
	pv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	pvB, _ := x509.MarshalECPrivateKey(pv)
	pub, _ := x509.MarshalPKIXPublicKey(pv.Public())

	pubK := mockKey{
		b: pub,
		c: 1,
	}

	pvK := mockKey{
		b: pvB,
		c: 1,
	}

	cipherData, err := Encrypt([]byte("hello"), pubK)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)

	data, err := Decrypt(cipherData, pvK)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), data)
}

/*
Scenario: Decrypt using ecies and Ed25519 key
	Given a cipher data encrypted by an Ed25519 key and an Ed25519 private key
	When I want to decrypt the data
	Then I will generate a common secret and decrypt the cipher and get a clear data
*/
func TestEciesDecryptWithEd25519(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	pubK := mockKey{
		b: pub,
		c: 0,
	}

	pvK := mockKey{
		b: pv,
		c: 0,
	}

	cipherData, err := Encrypt([]byte("hello"), pubK)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)

	data, err := Decrypt(cipherData, pvK)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), data)
}

type mockKey struct {
	b []byte
	c int
}

func (k mockKey) Bytes() []byte {
	return k.b
}

func (k mockKey) Curve() int {
	return k.c
}
