package shared

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//KeyPair describe shared keypair
type KeyPair struct {
	encPvKey string
	pubKey   string
}

//NewKeyPair creates a new shared keypair
func NewKeyPair(encPvKey, pubKey string) (KeyPair, error) {

	if encPvKey == "" {
		return KeyPair{}, errors.New("shared keys encrypted private key: is empty")
	}
	if _, err := hex.DecodeString(encPvKey); err != nil {
		return KeyPair{}, errors.New("shared keys encrypted private key: is not in hexadecimal format")
	}

	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return KeyPair{}, fmt.Errorf("shared keys: %s", err.Error())
	}

	return KeyPair{encPvKey, pubKey}, nil
}

//PublicKey returns the public key for the shared keypair
func (kp KeyPair) PublicKey() string {
	return kp.pubKey
}

//EncryptedPrivateKey returns the encrypted private key for the shared keypair
func (kp KeyPair) EncryptedPrivateKey() string {
	return kp.encPvKey
}

func (sk KeyPair) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		EncryptedPrivateKey string `json:"encrypted_private_key"`
		PublicKey           string `json:"public_key"`
	}{
		EncryptedPrivateKey: sk.EncryptedPrivateKey(),
		PublicKey:           sk.PublicKey(),
	})
}
