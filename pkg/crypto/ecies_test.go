package crypto

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Derivate keys from a secret
	Given a secret
	When I want to derivate X keys
	Then I get X derivated keys differents
*/
func TestDerivateKeys(t *testing.T) {
	keys, err := derivateKeys([]byte("hello"), 2)
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
	tag := authenticateMessage([]byte("my key"), []byte("my data"))
	assert.NotEmpty(t, tag)

	tag2 := authenticateMessage([]byte("my key"), []byte("my data"))
	assert.NotEmpty(t, tag2)

	assert.Equal(t, tag, tag2)
}

/*
Scenario: Encode a encryption with ed25519
	Given cipher data, a public key, and a tag (MAC)
	When I want to encode the encryption result
	Then I get byte slice withing all these information
*/
func TestEncodeEncryptionWithEd25519(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	c := newEncodedCipher(pub, []byte("cipher data"), []byte("tag MAC"))
	assert.NotEmpty(t, c)

	b := pub.Bytes()

	assert.EqualValues(t, b, c[:len(b)])
	assert.EqualValues(t, []byte("cipher data"), c[len(b):len(b)+len([]byte("cipher data"))])
	assert.EqualValues(t, []byte("tag MAC"), c[len(b)+len([]byte("cipher data")):])
}

/*
Scenario: Encode a encryption with ecdsa
	Given cipher data, a public key, and a tag (MAC)
	When I want to encode the encryption result
	Then I get byte slice withing all these information
*/
func TestEncodeEncryptionWithEcdsa(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	c := newEncodedCipher(pub, []byte("cipher data"), []byte("tag MAC"))
	assert.NotEmpty(t, c)

	b := pub.Bytes()

	assert.EqualValues(t, b, c[:len(b)])
	assert.EqualValues(t, []byte("cipher data"), c[len(b):len(b)+len([]byte("cipher data"))])
	assert.EqualValues(t, []byte("tag MAC"), c[len(b)+len([]byte("cipher data")):])
}

/*
Scenario: Decode a encryption with ed25519
	Given encryption encoded
	When I want to decode the encryption result
	Then I get encrypted data, a public key, and a tag (MAC)
*/
func TestDecodeEncryptionWithEd25519(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	tag := authenticateMessage([]byte("key"), []byte("message"))
	cipher := newEncodedCipher(pub, []byte("cipher data"), tag)
	assert.NotEmpty(t, cipher)

	rPub, em, tag, err := cipher.decode(ed25519ExtractRandomPublicKey)
	assert.Nil(t, err)
	assert.Equal(t, pub, rPub)
	assert.Equal(t, tag, tag)
	assert.Equal(t, []byte("cipher data"), em)
}

/*
Scenario: Decode a encryption with ecdsa
	Given encryption encoded
	When I want to decode the encryption result
	Then I get encrypted data, a public key, and a tag (MAC)
*/
func TestDecodeEncryptionWithEcdsa(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	tag := authenticateMessage([]byte("key"), []byte("message"))
	cipher := newEncodedCipher(pub, []byte("cipher data"), tag)
	assert.NotEmpty(t, cipher)

	rPub, em, tag, err := cipher.decode(ecdsaExtractRandomPublicKey(elliptic.P256()))
	assert.Nil(t, err)
	assert.Equal(t, pub, rPub)
	assert.Equal(t, tag, tag)
	assert.Equal(t, []byte("cipher data"), em)
}

/*
Scenario: Encrypt using ecies
	Given a dest public key and a data
	When I want to encrypt the data
	Then I get a cipher message
*/
func TestEciesEncrypt(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	cipherData, err := eciesEncrypt([]byte("hello"), pub, ed25519GenerateShared)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)
}

/*
Scenario: Decrypt using ecies
	Given a cipher data key and a private key
	When I want to decrypt the data
	Then I get a clear data
*/
func TestEciesDecrypt(t *testing.T) {
	pv, pub, _ := generateEd25519Keys(rand.Reader)
	cipherData, err := eciesEncrypt([]byte("hello"), pub, ed25519GenerateShared)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipherData)

	data, err := eciesDecrypt(cipherData, pv, ed25519GenerateShared, ed25519ExtractRandomPublicKey)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), data)
}
