package shared

import "github.com/uniris/uniris-core/pkg/crypto"

//NodeKeyPair represent a shared node keypair
type NodeKeyPair struct {
	pub crypto.PublicKey
	pv  crypto.PrivateKey
}

//NewNodeKeyPair creates a new node keypairs
func NewNodeKeyPair(pub crypto.PublicKey, pv crypto.PrivateKey) (NodeKeyPair, error) {
	return NodeKeyPair{pub, pv}, nil
}

//PublicKey returns the shared node public key
func (kp NodeKeyPair) PublicKey() crypto.PublicKey {
	return kp.pub
}

//PrivateKey returns the shared node private key
func (kp NodeKeyPair) PrivateKey() crypto.PrivateKey {
	return kp.pv
}
