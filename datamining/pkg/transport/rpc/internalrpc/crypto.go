package internalrpc

//Hasher define methods to hash incoming data
type Hasher interface {
	HashKeychainJSON(*KeychainDataJSON) (string, error)
	HashBiometricJSON(*BioDataJSON) (string, error)
}
