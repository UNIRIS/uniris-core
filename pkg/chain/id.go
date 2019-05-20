package chain

// import (
// 	"errors"
// )

// //ID represents a ID transaction
// type ID struct {
// 	Transaction
// }

// //NewID creates a new ID transaction
// func NewID(tx Transaction) (ID, error) {

// 	if tx.txType != IDTransactionType {
// 		return ID{}, errors.New("invalid type of transaction")
// 	}

// 	if _, exist := tx.data["encrypted_address_by_node"]; !exist {
// 		return ID{}, errors.New("missing ID data: 'encrypted_address_by_node'")
// 	}
// 	if _, exist := tx.data["encrypted_address_by_id"]; !exist {
// 		return ID{}, errors.New("missing ID data: 'encrypted_address_by_id'")
// 	}
// 	if _, exist := tx.data["encrypted_aes_key"]; !exist {
// 		return ID{}, errors.New("missing ID data: 'encrypted_aes_key'")
// 	}

// 	return ID{
// 		Transaction: tx,
// 	}, nil
// }

// //EncryptedAddrBy returns the encrypted keychain address with the  public key
// func (id ID) EncryptedAddrBy() []byte {
// 	return id.data["encrypted_address_by_node"]
// }

// //EncryptedAddrByID returns the encrypted keychain address with the ID public key
// func (id ID) EncryptedAddrByID() []byte {
// 	return id.data["encrypted_address_by_id"]
// }

// //EncryptedAESKey returns the encrypted AES key with the ID public key
// func (id ID) EncryptedAESKey() []byte {
// 	return id.data["encrypted_aes_key"]
// }
