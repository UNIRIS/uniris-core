package shared

import (
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//KeyPair represent a shared node keypair
type KeyPair struct {
	pubKey string
	pvKey  string
}

//NewKeyPair creates a new node keypairs
func NewKeyPair(pubKey string, pvKey string) (KeyPair, error) {
	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return KeyPair{}, fmt.Errorf("shared node keys: %s", err.Error())
	}

	if _, err := crypto.IsPrivateKey(pvKey); err != nil {
		return KeyPair{}, fmt.Errorf("shared node keys: %s", err.Error())
	}

	return KeyPair{
		pubKey: pubKey,
		pvKey:  pvKey,
	}, nil
}

//PublicKey returns the shared node public key
func (mKP KeyPair) PublicKey() string {
	return mKP.pubKey
}

//PrivateKey returns the shared node private key
func (mKP KeyPair) PrivateKey() string {
	return mKP.pvKey
}
