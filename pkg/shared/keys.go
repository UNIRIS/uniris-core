package shared

import (
	"encoding/json"
	"errors"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//KeyReadWriter wraps the keys reading and persisting
type KeyReadWriter interface {
	KeyReader
	KeyWriter
}

//KeyReader performs queries to retrieve shared keys
type KeyReader interface {

	//EmitterCrossKeypairs retrieve the list of the cross emitter keys
	EmitterCrossKeypairs() ([]EmitterCrossKeyPair, error)

	//FirstEmitterCrossKeypair retrieves the first public key
	FirstEmitterCrossKeypair() (EmitterCrossKeyPair, error)

	//CrossEmitterPublicKeys retrieves the public keys of the cross emitter keys
	CrossEmitterPublicKeys() ([]crypto.PublicKey, error)

	//FirstNodeCrossKeypair retrieve the first shared crosskeys for the nodes
	FirstNodeCrossKeypair() (NodeCrossKeyPair, error)

	//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
	LastNodeCrossKeypair() (NodeCrossKeyPair, error)

	//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
	AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error)

	//IsAuthorizedNode check if the public Key is on the authorized list
	IsAuthorizedNode(crypto.PublicKey) bool
}

//KeyWriter performs persistance of the shared keys
type KeyWriter interface {

	//WriteAuthorizedNode inserts a new node public key as an authorized node
	WriteAuthorizedNode(pubKey crypto.PublicKey) error
}

//EmitterCrossKeyPair represents cross keypair with an encrypted private key
type EmitterCrossKeyPair struct {
	encPvKey []byte
	pubKey   crypto.PublicKey
}

//NewEmitterCrossKeyPair creates a new emitter cross keypair
func NewEmitterCrossKeyPair(encPvKey []byte, pubKey crypto.PublicKey) (EmitterCrossKeyPair, error) {

	if encPvKey == nil || len(encPvKey) == 0 {
		return EmitterCrossKeyPair{}, errors.New("missing emitter cross private key")
	}

	if pubKey == nil {
		return EmitterCrossKeyPair{}, errors.New("missing emitter cross public key")
	}

	return EmitterCrossKeyPair{encPvKey, pubKey}, nil
}

//PublicKey returns the emitter cross public key
func (kp EmitterCrossKeyPair) PublicKey() crypto.PublicKey {
	return kp.pubKey
}

//EncryptedPrivateKey returns the emitter cross private key
func (kp EmitterCrossKeyPair) EncryptedPrivateKey() []byte {
	return kp.encPvKey
}

//MarshalJSON serialize the emitter cross keypair in JSON
func (kp EmitterCrossKeyPair) MarshalJSON() ([]byte, error) {
	pubB, err := kp.pubKey.Marshal()
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string][]byte{
		"encrypted_private_key": kp.EncryptedPrivateKey(),
		"public_key":            pubB,
	})
}

//NodeCrossKeyPair represents a node cross keypair
type NodeCrossKeyPair struct {
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

//NewNodeCrossKeyPair creates a new cross keypair
func NewNodeCrossKeyPair(pubKey crypto.PublicKey, pvKey crypto.PrivateKey) (NodeCrossKeyPair, error) {

	if pubKey == nil {
		return NodeCrossKeyPair{}, errors.New("missing node cross public key")
	}

	if pvKey == nil {
		return NodeCrossKeyPair{}, errors.New("missing node cross private key")
	}

	return NodeCrossKeyPair{
		publicKey:  pubKey,
		privateKey: pvKey,
	}, nil
}

//PublicKey returns the node cross public key
func (kp NodeCrossKeyPair) PublicKey() crypto.PublicKey {
	return kp.publicKey
}

//PrivateKey returns the node cross private key
func (kp NodeCrossKeyPair) PrivateKey() crypto.PrivateKey {
	return kp.privateKey
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func IsEmitterKeyAuthorized(emPubKey crypto.PublicKey) (bool, error) {
	//TODO: request smart contract

	return true, nil
}
