package chain

// import (
// 	"errors"
// )

// //Keychain represents a keychain transaction
// type Keychain struct {
// 	Transaction
// }

// //NewKeychain create a new Keychain transaction
// func NewKeychain(tx Transaction) (Keychain, error) {

// 	if tx.txType != KeychainTransactionType {
// 		return Keychain{}, errors.New("invalid type of transaction")
// 	}

// 	if _, exist := tx.data["encrypted_address_by_node"]; !exist {
// 		return Keychain{}, errors.New("missing keychain data: 'encrypted_address_by_node'")
// 	}

// 	if _, exist := tx.data["encrypted_wallet"]; !exist {
// 		return Keychain{}, errors.New("missing keychain data: 'encrypted_wallet'")
// 	}

// 	return Keychain{
// 		Transaction: tx,
// 	}, nil
// }

// //EncryptedAddrBy returns the encrypted keychain address by the shared node key
// func (k Keychain) EncryptedAddrBy() []byte {
// 	return k.data["encrypted_address_by_node"]
// }

// //EncryptedWallet returns encrypted wallet by the person AES key
// func (k Keychain) EncryptedWallet() []byte {
// 	return k.data["encrypted_wallet"]
// }
