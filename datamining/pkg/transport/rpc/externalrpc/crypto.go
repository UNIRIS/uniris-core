package externalrpc

//Signer define methods to handle signatures
type Signer interface {
	SignBiometric(data BiometricJSON, pvKey string) (string, error)
	SignKeychain(data KeychainJSON, pvKey string) (string, error)
}
