package shared

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//EmitterKeys define list of shared emitter keys
type EmitterKeys []EmitterKeyPair

//RequestKey returns the public key used to request between clients and miners
func (em EmitterKeys) RequestKey() string {
	return em[0].pubKey
}

//EmitterKeyPair describe shared emitter keypair
type EmitterKeyPair struct {
	encPvKey string
	pubKey   string
}

//NewEmitterKeyPair creates a new shared emitter keypair
func NewEmitterKeyPair(encPvKey, pubKey string) (EmitterKeyPair, error) {

	if encPvKey == "" {
		return EmitterKeyPair{}, errors.New("shared emitter keys encrypted private key: is empty")
	}
	if _, err := hex.DecodeString(encPvKey); err != nil {
		return EmitterKeyPair{}, errors.New("shared emitter keys encrypted private key: is not in hexadecimal format")
	}

	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return EmitterKeyPair{}, fmt.Errorf("shared emitter keys: %s", err.Error())
	}

	return EmitterKeyPair{encPvKey, pubKey}, nil
}

//PublicKey returns the public key for the shared  keypair
func (kp EmitterKeyPair) PublicKey() string {
	return kp.pubKey
}

//EncryptedPrivateKey returns the encrypted private key for the shared  keypair
func (kp EmitterKeyPair) EncryptedPrivateKey() string {
	return kp.encPvKey
}

//MarshalJSON serialize the keypair in JSON
func (kp EmitterKeyPair) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"encrypted_private_key": kp.EncryptedPrivateKey(),
		"public_key":            kp.PublicKey(),
	})
}
