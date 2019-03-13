package shared

import (
	"errors"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//EmitterKeys define list of shared emitter keys
type EmitterKeys []EmitterKeyPair

//RequestKey returns the public key used to request between clients and nodes
func (em EmitterKeys) RequestKey() crypto.PublicKey {
	return em[0].pubKey
}

//EmitterKeyPair describe shared emitter keypair
type EmitterKeyPair struct {
	encPvKey []byte
	pubKey   crypto.PublicKey
}

//NewEmitterKeyPair creates a new shared emitter keypair
func NewEmitterKeyPair(encPvKey []byte, pubKey crypto.PublicKey) (EmitterKeyPair, error) {

	if len(encPvKey) == 0 {
		return EmitterKeyPair{}, errors.New("shared emitter keys encrypted private key: is empty")
	}

	return EmitterKeyPair{encPvKey, pubKey}, nil
}

//PublicKey returns the public key for the shared  keypair
func (kp EmitterKeyPair) PublicKey() crypto.PublicKey {
	return kp.pubKey
}

//EncryptedPrivateKey returns the encrypted private key for the shared  keypair
func (kp EmitterKeyPair) EncryptedPrivateKey() []byte {
	return kp.encPvKey
}
