package chain

import (
	"encoding/hex"
	"errors"
)

//ID represents a ID transaction
type ID struct {
	Transaction
}

//NewID creates a new ID transaction
func NewID(tx Transaction) (ID, error) {

	if tx.txType != IDTransactionType {
		return ID{}, errors.New("transaction: invalid type of transaction")
	}

	addr, exist := tx.data["encrypted_address_by_node"]
	if !exist {
		return ID{}, errors.New("transaction: missing data ID 'encrypted_address_by_node'")
	}
	addrID, exist := tx.data["encrypted_address_by_id"]
	if !exist {
		return ID{}, errors.New("transaction: missing data ID 'encrypted_address_by_id'")
	}
	aesKey, exist := tx.data["encrypted_aes_key"]
	if !exist {
		return ID{}, errors.New("transaction: missing data ID 'encrypted_aes_key'")
	}

	if _, err := hex.DecodeString(aesKey); err != nil {
		return ID{}, errors.New("transaction: id encrypted aes key is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(addrID); err != nil {
		return ID{}, errors.New("transaction: id encrypted address for id is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return ID{}, errors.New("transaction: id encrypted address for node is not in hexadecimal format")
	}

	return ID{
		Transaction: tx,
	}, nil
}

//EncryptedAddrBy returns the encrypted keychain address with the  public key
func (id ID) EncryptedAddrBy() string {
	return id.data["encrypted_address_by_node"]
}

//EncryptedAddrByID returns the encrypted keychain address with the ID public key
func (id ID) EncryptedAddrByID() string {
	return id.data["encrypted_address_by_id"]
}

//EncryptedAESKey returns the encrypted AES key with the ID public key
func (id ID) EncryptedAESKey() string {
	return id.data["encrypted_aes_key"]
}
