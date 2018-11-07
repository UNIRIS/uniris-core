package internalrpc

import "github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

type decrypter interface {
	Decrypter
	rpc.Decrypter
}

//Hasher define methods to hash incoming data
type Hasher interface {
	HashKeychainJSON(*KeychainDataFromJSON) (string, error)
	HashBiometricJSON(*BioDataFromJSON) (string, error)
}

//Decrypter define decryption methods for the internal rpc methods
type Decrypter interface {
	DecryptHashPerson(hash string, pvKey string) (string, error)
}
