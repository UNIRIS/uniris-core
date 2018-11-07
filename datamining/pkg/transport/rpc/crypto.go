package rpc

//Decrypter define methods to decrypt data for RPC methods
type Decrypter interface {
	DecryptHashPerson(hash string, pvKey string) (string, error)
	DecryptCipherAddress(cipherAddr string, pvKey string) (string, error)
	DecryptTransactionData(data string, pvKey string) (string, error)
}

//Hasher define methods to hash incoming data
type Hasher interface {
	HashKeychainJSON(*KeychainDataJSON) (string, error)
	HashBiometricJSON(*BioDataJSON) (string, error)
}

//Signer define methods to handle signatures
type Signer interface {
	SignBiometric(data BiometricJSON, pvKey string) (string, error)
	SignKeychain(data KeychainJSON, pvKey string) (string, error)
}
