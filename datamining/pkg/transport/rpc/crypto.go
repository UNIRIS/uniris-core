package rpc

//Decrypter define methods to decrypt data for RPC methods
type Decrypter interface {
	DecryptCipherAddress(cipherAddr string, pvKey string) (string, error)
	DecryptTransactionData(data string, pvKey string) (string, error)
}
