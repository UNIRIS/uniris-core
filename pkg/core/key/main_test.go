package main

import (
	"bytes"
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
 Scenario: Generate an ECDSA key by loading its plugin
	Given a entropy source and an curve choice
	When I want to generate the keys, it loads the plugin ECDSA
	Then I get a valid ECDSA key (same if using the ECSDA standard lib)
*/
func TestGenerateECDSAKey(t *testing.T) {

	src1 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	pvKey, pubKey, err := GenerateKeys(P256Curve, src1)
	assert.Nil(t, err)
	assert.NotNil(t, pvKey)
	assert.NotNil(t, pubKey)

	src2 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	ecdsaPvKey, err := ecdsa.GenerateKey(elliptic.P256(), src2)
	assert.Nil(t, err)
	pvBytes, err := x509.MarshalECPrivateKey(ecdsaPvKey)
	assert.Nil(t, err)
	assert.Equal(t, pvBytes, pvKey.(ECKey).Bytes())
}

/*
 Scenario: Generate an Ed25519 key by loading its plugin
	Given a entropy source and an curve choice
	When I want to generate the keys, it loads the plugin Ed25519
	Then I get a valid Ed25519 key (same if using the Ed25519 standard lib)
*/
func TestGenerateEd25519Key(t *testing.T) {

	src1 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	pvKey, pubKey, err := GenerateKeys(Ed25519Curve, src1)
	assert.Nil(t, err)
	assert.NotNil(t, pvKey)
	assert.NotNil(t, pubKey)

	src2 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	_, pv2, err := ed25519.GenerateKey(src2)
	assert.Nil(t, err)
	assert.EqualValues(t, pv2, pvKey.(ECKey).Bytes())
}

/*
Scenario: Marshal ECDSA keys
	Given an ECDSA keypair
	When I want to marshal them
	Then I get as first bytes the curve and the rest as validate key
*/
func TestMarshalECDSAKeys(t *testing.T) {

	src1 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	pvKey, pubKey, err := GenerateKeys(P256Curve, src1)
	assert.Nil(t, err)

	assert.Equal(t, int(pvKey.(ECKey).Marshal()[0]), P256Curve)
	assert.Equal(t, int(pubKey.(ECKey).Marshal()[0]), P256Curve)

	src2 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	ecdsaPvKey, err := ecdsa.GenerateKey(elliptic.P256(), src2)

	pvBytes, _ := x509.MarshalECPrivateKey(ecdsaPvKey)
	assert.Equal(t, pvKey.(ECKey).Marshal()[1:], pvBytes)

	pubBytes, _ := x509.MarshalPKIXPublicKey(ecdsaPvKey.Public())
	assert.Equal(t, pubKey.(ECKey).Marshal()[1:], pubBytes)
}

/*
Scenario: Marshal Ed25519 keys
	Given an Ed25519 keypair
	When I want to marshal them
	Then I get as first bytes the curve and the rest as validate key
*/
func TestMarshalEd25519Keys(t *testing.T) {

	src1 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	pvKey, pubKey, err := GenerateKeys(Ed25519Curve, src1)
	assert.Nil(t, err)

	assert.Equal(t, int(pvKey.(ECKey).Marshal()[0]), Ed25519Curve)
	assert.Equal(t, int(pubKey.(ECKey).Marshal()[0]), Ed25519Curve)

	src2 := bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	pub, _, _ := ed25519.GenerateKey(src2)

	assert.EqualValues(t, pubKey.(ECKey).Marshal()[1:], pub)
}

/*
Scenario: Parse an marshalled ECDSA public key
	Given an marshalled ECDSA public key
	When I want to parse it
	Then I get a valid ECDSA public key with its identified curve
*/
func TestParseMarshalledECDSAPublicKey(t *testing.T) {
	_, pubKey, err := GenerateKeys(P256Curve, rand.Reader)
	assert.Nil(t, err)
	pub, err := ParsePublicKey(pubKey.(ECKey).Marshal())
	assert.Nil(t, err)
	assert.Equal(t, pubKey.(ECKey).Bytes(), pub.(ECKey).Bytes())
	assert.Equal(t, P256Curve, pub.(ECKey).Curve())
}

/*
Scenario: Parse an marshalled Ed25519 public key
	Given an marshalled Ed25519 public key
	When I want to parse it
	Then I get a valid Ed25519 public key
*/
func TestParseMarshalledEd25519PublicKey(t *testing.T) {
	_, pubKey, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)
	pub, err := ParsePublicKey(pubKey.(ECKey).Marshal())
	assert.Nil(t, err)
	assert.Equal(t, pubKey.(ECKey).Bytes(), pub.(ECKey).Bytes())
	assert.Equal(t, Ed25519Curve, pub.(ECKey).Curve())
}

/*
Scenario: Parse an marshalled ECDSA private key
	Given an marshalled ECDSA private key
	When I want to parse it
	Then I get a valid ECDSA private key with its identified curve
*/
func TestParseMarshalledECDSAPrivateKey(t *testing.T) {
	pvKey, _, err := GenerateKeys(P256Curve, rand.Reader)
	assert.Nil(t, err)
	pv, err := ParsePrivateKey(pvKey.(ECKey).Marshal())
	assert.Nil(t, err)
	assert.Equal(t, pvKey.(ECKey).Bytes(), pv.(ECKey).Bytes())
	assert.Equal(t, P256Curve, pv.(ECKey).Curve())
}

/*
Scenario: Parse an marshalled Ed25519 private key
	Given an marshalled Ed25519 private key
	When I want to parse it
	Then I get a valid Ed25519 private key
*/
func TestParseMarshalledEd25519PrivateKey(t *testing.T) {
	pvKey, _, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)
	pv, err := ParsePrivateKey(pvKey.(ECKey).Marshal())
	assert.Nil(t, err)
	assert.Equal(t, pvKey.(ECKey).Bytes(), pv.(ECKey).Bytes())
	assert.Equal(t, Ed25519Curve, pv.(ECKey).Curve())
}

/*
Scenario: Generate a signature with a private key by loading its plugin
	Given a private key and a message
	When I want to sign the message, it loads the plugin of the private key curve
	Then I get a signature
*/
func TestSign(t *testing.T) {
	pvKey, pub, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)
	sig, err := pvKey.(PrivateKey).Sign([]byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, sig)
	assert.True(t, ed25519.Verify(pub.(ECKey).Bytes(), []byte("hello"), sig))
}

/*
Scenario: Verify a signature with a public key by loading its plugin
	Given an public key, a message and a signature
	When I want to verify the signature, it loads the plugin relative of the public key curve
	Then I get not error and a true response
*/
func TestVerify(t *testing.T) {
	pvKey, pub, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)

	sig, err := pvKey.(PrivateKey).Sign([]byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, sig)

	ok, err := pub.(PublicKey).Verify([]byte("hello"), sig)
	assert.Nil(t, err)
	assert.True(t, ok)
}

/*
Scenario: Encrypt a message with a public key by loading the ECIES plugin
	Given a public key
	When I want to encrypt a message
	Then I call the ECIES plugin to encrypt it and return the cipher message
*/
func TestEncrypt(t *testing.T) {
	_, pub, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)

	cipher, err := pub.(PublicKey).Encrypt([]byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)
}

/*
Scenario: Decrypt an encrypted message with a private key by loading the ECIES plugin
	Given a private key and an encrypted message
	When I want to decrypt the message
	Then I call the ECIES plugin to decrypt it and return the clear message
*/
func TestDecrypt(t *testing.T) {
	pv, pub, err := GenerateKeys(Ed25519Curve, rand.Reader)
	assert.Nil(t, err)

	cipher, err := pub.(PublicKey).Encrypt([]byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)

	msg, err := pv.(PrivateKey).Decrypt(cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), msg)

}
