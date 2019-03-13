package shared

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//KeyReader performs queries to retrieve shared keys
type KeyReader interface {

	//EmitterCrossKeypairs retrieve the list of the cross emitter keys
	EmitterCrossKeypairs() ([]EmitterCrossKeyPair, error)

	//FirstEmitterCrossKeypair retrieves the first public key
	FirstEmitterCrossKeypair() (EmitterCrossKeyPair, error)

	//CrossEmitterPublicKeys retrieves the public keys of the cross emitter keys
	CrossEmitterPublicKeys() ([]string, error)

	//FirstNodeCrossKeypair retrieve the first shared crosskeys for the nodes
	FirstNodeCrossKeypair() (NodeCrossKeyPair, error)

	//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
	LastNodeCrossKeypair() (NodeCrossKeyPair, error)

	//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
	AuthorizedNodesPublicKeys() ([]string, error)
}

//EmitterCrossKeyPair represents cross keypair with an encrypted private key
type EmitterCrossKeyPair struct {
	encPvKey string
	pubKey   string
}

//NewEmitterCrossKeyPair creates a new emitter cross keypair
func NewEmitterCrossKeyPair(encPvKey, pubKey string) (EmitterCrossKeyPair, error) {

	if encPvKey == "" {
		return EmitterCrossKeyPair{}, errors.New("missing emitter cross private key")
	}
	if _, err := hex.DecodeString(encPvKey); err != nil {
		return EmitterCrossKeyPair{}, errors.New("emitter cross private key is not in hexadecimal format")
	}

	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return EmitterCrossKeyPair{}, fmt.Errorf("emitter cross: %s", err.Error())
	}

	return EmitterCrossKeyPair{encPvKey, pubKey}, nil
}

//PublicKey returns the emitter cross public key
func (kp EmitterCrossKeyPair) PublicKey() string {
	return kp.pubKey
}

//EncryptedPrivateKey returns the emitter cross private key
func (kp EmitterCrossKeyPair) EncryptedPrivateKey() string {
	return kp.encPvKey
}

//MarshalJSON serialize the emitter cross keypair in JSON
func (kp EmitterCrossKeyPair) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"encrypted_private_key": kp.EncryptedPrivateKey(),
		"public_key":            kp.PublicKey(),
	})
}

//NodeCrossKeyPair represents a node cross keypair
type NodeCrossKeyPair struct {
	privateKey string
	publicKey  string
}

//NewNodeCrossKeyPair creates a new cross keypair
func NewNodeCrossKeyPair(pubKey string, pvKey string) (NodeCrossKeyPair, error) {
	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return NodeCrossKeyPair{}, fmt.Errorf("node cross public key: %s", err.Error())
	}

	if _, err := crypto.IsPrivateKey(pvKey); err != nil {
		return NodeCrossKeyPair{}, fmt.Errorf("node cross private key: %s", err.Error())
	}

	return NodeCrossKeyPair{
		publicKey:  pubKey,
		privateKey: pvKey,
	}, nil
}

//PublicKey returns the node cross public key
func (kp NodeCrossKeyPair) PublicKey() string {
	return kp.publicKey
}

//PrivateKey returns the node cross private key
func (kp NodeCrossKeyPair) PrivateKey() string {
	return kp.privateKey
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func IsEmitterKeyAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract

	return true, nil
}
