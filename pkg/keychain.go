package uniris

import (
	"encoding/json"
	"errors"
)

//Keychain represents a keychain transaction
type Keychain struct {
	encAddr   string
	encWallet string
	Transaction
}

type keychainData struct {
	EncryptedAddress string `json:"encrypted_address"`
	EncryptedWallet  string `json:"encrypted_wallet"`
}

//NewKeychain creates a keychain transaction by extracting the transaction data
func NewKeychain(tx Transaction) (Keychain, error) {
	var data keychainData
	if err := json.Unmarshal([]byte(tx.Data()), &data); err != nil {
		return Keychain{}, err
	}

	if data.EncryptedAddress == "" || data.EncryptedWallet == "" {
		return Keychain{}, errors.New("Missing Keychain transaction data")
	}

	return Keychain{
		encAddr:     data.EncryptedAddress,
		encWallet:   data.EncryptedWallet,
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

func (k Keychain) ToTransaction() (tx Transaction, err error) {
	data := keychainData{
		EncryptedAddress: k.EncryptedAddrByRobot(),
		EncryptedWallet:  k.EncryptedWallet(),
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}

	tx, err = NewTransaction(k.address, KeychainTransactionType, string(b), k.timestamp, k.pubKey, k.sig, k.emSig, k.prop, k.txHash)
	if err != nil {
		return
	}
	if err = tx.AddMining(k.masterV, k.confirmValids); err != nil {
		return
	}
	if k.prevTx != nil {
		tx.Chain(k.prevTx)
	}
	return tx, nil
}
