package shared

import (
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//MinerKeyPair represent a shared miner keypair
type MinerKeyPair struct {
	pubKey string
	pvKey  string
}

//NewMinerKeyPair creates a new miner keypair
func NewMinerKeyPair(pubKey string, pvKey string) (MinerKeyPair, error) {
	if _, err := crypto.IsPublicKey(pubKey); err != nil {
		return MinerKeyPair{}, fmt.Errorf("shared miner keys: %s", err.Error())
	}

	if _, err := crypto.IsPrivateKey(pvKey); err != nil {
		return MinerKeyPair{}, fmt.Errorf("shared miner keys: %s", err.Error())
	}

	return MinerKeyPair{
		pubKey: pubKey,
		pvKey:  pvKey,
	}, nil
}

//PublicKey returns the shared miner public key
func (mKP MinerKeyPair) PublicKey() string {
	return mKP.pubKey
}

//PrivateKey returns the shared miner private key
func (mKP MinerKeyPair) PrivateKey() string {
	return mKP.pvKey
}
