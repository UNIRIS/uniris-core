package chain

import (
	"encoding/hex"
	"errors"
)

//Keychain represents a keychain transaction
type Keychain struct {
	Transaction
}

//NewKeychain create a new Keychain transaction
func NewKeychain(tx Transaction) (Keychain, error) {

	if tx.txType != KeychainTransactionType {
		return Keychain{}, errors.New("transaction: invalid type of transaction")
	}

	addr, exist := tx.data["encrypted_address_by_node"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_address_by_node'")
	}

	wallet, exist := tx.data["encrypted_wallet"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_wallet'")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted address for node is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(wallet); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted wallet is not in hexadecimal format")
	}

	return Keychain{
		Transaction: tx,
	}, nil
}

//EncryptedAddrBy returns the encrypted keychain address by the shared node key
func (k Keychain) EncryptedAddrBy() string {
	return k.data["encrypted_address_by_node"]
}

//EncryptedWallet returns encrypted wallet by the person AES key
func (k Keychain) EncryptedWallet() string {
	return k.data["encrypted_wallet"]
}
