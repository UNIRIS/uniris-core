package main

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
func TestEncryptWithInvalidKey(t *testing.T) {
	_, err := Encrypt([]byte("hello"), []byte("hello"))
	assert.NotNil(t, err)
}

/*
Scenario: Encrypt with AES
	Given a key and a message
	When I want to encrypt using AES
	Then I get a cipher message and no error
*/
func TestEncrypt(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("mykey"))
	cipher, err := Encrypt(h.Sum(nil), []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)
}

/*
Scenario: Decrypt with AES
	Given a key and a cipher message
	When I want to decrypt using AES
	Then I get the clear message and no error
*/
func TestDecrypt(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("mykey"))
	key := h.Sum(nil)
	cipher, err := Encrypt(key, []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)

	data, err := Decrypt(key, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), data)
}

/*
Scenario: Decrypt with AES with invalid key
	Given a invalid key and a cipher message
	When I want to decrypt using AES
	Then I get an error
*/
func TestDecryptWithInvalidKey(t *testing.T) {

	_, err := Decrypt([]byte("hello"), []byte("hello"))
	assert.NotNil(t, err)

	h := sha256.New()
	h.Write([]byte("mykey"))
	key := h.Sum(nil)
	h.Reset()
	cipher, err := Encrypt(key, []byte("hello"))
	assert.Nil(t, err)
	assert.NotNil(t, cipher)

	h.Write([]byte("test"))

	_, err = Decrypt(h.Sum(nil), cipher)
	assert.NotNil(t, err)
	assert.Equal(t, "cipher: message authentication failed", err.Error())
}
