package shared

import "github.com/uniris/uniris-core/pkg/crypto"

//TechDatabaseReader wraps the shared emitter and node storage
type TechDatabaseReader interface {
	EmitterDatabaseReader

	//NodeLastKeys retrieve the last shared node keys from the Tech DB
	NodeLastKeys() (NodeKeyPair, error)
}

//EmitterDatabaseReader handles queries for the shared emitter information
type EmitterDatabaseReader interface {
	//EmitterKeys retrieve the shared emitter key from the Tech DB
	EmitterKeys() (EmitterKeys, error)
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func IsEmitterKeyAuthorized(emPubKey crypto.PublicKey) (bool, error) {
	//TODO: request smart contract

	return true, nil
}
