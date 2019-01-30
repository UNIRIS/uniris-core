package transaction

import (
	"encoding/hex"
	"errors"
)

//KeychainRepository manages the Keychain transaction storage
type KeychainRepository interface {
	GetKeychain(addr string) (*Keychain, error)
	FindKeychainByHash(txHash string) (*Keychain, error)
	FindLastKeychain(addr string) (*Keychain, error)
	StoreKeychain(kc Keychain) error
}

//Keychain represents a keychain transaction
type Keychain struct {
	encAddr   string
	encWallet string
	Transaction
}

//NewKeychain creates a keychain transaction by extracting the transaction data
func NewKeychain(tx Transaction) (Keychain, error) {

	if tx.Type() != KeychainType {
		return Keychain{}, errors.New("transaction: invalid type of transaction")
	}

	addr, exist := tx.data["encrypted_address"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_address'")
	}

	wallet, exist := tx.data["encrypted_wallet"]
	if !exist {
		return Keychain{}, errors.New("transaction: missing data keychain: 'encrypted_wallet'")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted address is not in hexadecimal format")
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

//ToTransaction converts back Keychain to transaction root
func (k Keychain) ToTransaction() (Transaction, error) {
	data := map[string]string{
		"encrypted_address": k.encAddr,
		"encrypted_wallet":  k.encWallet,
	}
	tx, err := New(k.Address(), KeychainType, data, k.Timestamp(), k.PublicKey(), k.Signature(), k.EmitterSignature(), k.Proposal(), k.TransactionHash())
	if err != nil {
		return Transaction{}, nil
	}
	if k.masterV.pow != "" && len(k.confirmValids) > 0 {
		if err = tx.AddMining(k.MasterValidation(), k.ConfirmationsValidations()); err != nil {
			return Transaction{}, nil
		}
	}

	if k.PreviousTransaction() != nil {
		if err = tx.Chain(k.PreviousTransaction()); err != nil {
			return Transaction{}, nil
		}
	}
	return tx, nil
}
