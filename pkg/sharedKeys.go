package uniris

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//SharedKeys describe shared keypair
type SharedKeys struct {
	encPvKey string
	pubKey   string
}

//NewSharedKeyPair creates a new proposed keypair
func NewSharedKeyPair(encPvKey, pubKey string) (SharedKeys, error) {

	if encPvKey == "" {
		return SharedKeys{}, errors.New("Shared keys encrypted private key: is empty")
	}
	if _, err := hex.DecodeString(encPvKey); err != nil {
		return SharedKeys{}, errors.New("Shared keys encrypted private key: is not in hexadecimal format")
	}

	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return SharedKeys{}, fmt.Errorf("Shared keys: %s", err.Error())
	}

	return SharedKeys{encPvKey, pubKey}, nil
}

//PublicKey returns the public key for the proposed keypair
func (sK SharedKeys) PublicKey() string {
	return sK.pubKey
}

//EncryptedPrivateKey returns the encrypted private key for the proposed keypair
func (sK SharedKeys) EncryptedPrivateKey() string {
	return sK.encPvKey
}

func (sk SharedKeys) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		EncryptedPrivateKey string `json:"encrypted_private_key"`
		PublicKey           string `json:"public_key"`
	}{
		EncryptedPrivateKey: sk.EncryptedPrivateKey(),
		PublicKey:           sk.PublicKey(),
	})
}
