package chain

import (
	"encoding/hex"
	"errors"
)

//Keychain represents a keychain transaction
type Keychain struct {
	encAddr   string
	encWallet string
	Transaction
}

//NewKeychain create a new Keychain transaction
func NewKeychain(tx Transaction) (Keychain, error) {

	if tx.txType != KeychainTransactionType {
		return Keychain{}, errors.New("transaction: invalid type of transaction")
	}

	addr, exist := tx.data["encrypted_address_by_miner"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_address_by_miner'")
	}

	wallet, exist := tx.data["encrypted_wallet"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_wallet'")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted address for miner is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(wallet); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted wallet is not in hexadecimal format")
	}

	return Keychain{
		encAddr:     addr,
		encWallet:   wallet,
		Transaction: tx,
	}, nil
}

//EncryptedAddrByRobot returns the encrypted keychain address by the shared robot key
func (k Keychain) EncryptedAddrByRobot() string {
	return k.encAddr
}

//EncryptedWallet returns encrypted wallet by the person AES key
func (k Keychain) EncryptedWallet() string {
	return k.encWallet
}
