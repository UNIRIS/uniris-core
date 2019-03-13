package crypto

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Encrypt with AES with invalid key
	Given a invalid key and a message
	When I want to encrypt using AES
	Then I get an error
*/
func TestAESEncryptWithInvalidKey(t *testing.T) {
	_, err := AESEncrypt([]byte("hello"), []byte("hello"))
	assert.NotNil(t, err)
}

/*
Scenario: Encrypt with AES
	Given a key and a message
	When I want to encrypt using AES
	Then I get a cipher message and no error
*/
func TestAESEncrypt(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("mykey"))
	cipher, err := AESEncrypt(h.Sum(nil), []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)
}

/*
Scenario: Decrypt with AES
	Given a key and a cipher message
	When I want to decrypt using AES
	Then I get the clear message and no error
*/
func TestAESDecrypt(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("mykey"))
	key := h.Sum(nil)
	cipher, err := AESEncrypt(key, []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)

	data, err := AESDecrypt(key, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), data)
}

/*
Scenario: Decrypt with AES with invalid key
	Given a invalid key and a cipher message
	When I want to decrypt using AES
	Then I get an error
*/
func TestAESDecryptWithInvalidKey(t *testing.T) {

	_, err := AESDecrypt([]byte("hello"), []byte("hello"))
	assert.NotNil(t, err)

	h := sha256.New()
	h.Write([]byte("mykey"))
	key := h.Sum(nil)
	h.Reset()
	cipher, err := AESEncrypt(key, []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)

	h.Write([]byte("test"))

	_, err = AESDecrypt(h.Sum(nil), cipher)
	assert.NotNil(t, err)
	assert.Equal(t, "cipher: message authentication failed", err.Error())
}
