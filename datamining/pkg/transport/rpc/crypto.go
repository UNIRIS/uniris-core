package rpc

//Decrypter define methods to decrypt data for RPC methods
type Decrypter interface {
	DecryptHashPerson(hash string, pvKey string) (string, error)
	DecryptCipherAddress(cipherAddr string, pvKey string) (string, error)
	DecryptTransactionData(data string, pvKey string) (string, error)
}
