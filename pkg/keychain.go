package uniris

import (
	"encoding/json"
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
func NewKeychain(tx Transaction) (kc Keychain, err error) {
	var data keychainData
	if err = json.Unmarshal([]byte(tx.Data()), &data); err != nil {
		return
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

	tx = NewTransactionBase(k.address, KeychainTransactionType, string(b), k.timestamp, k.pubKey, k.sig, k.emSig, k.prop, k.txHash)
	if k.prevTx == nil {
		return NewMinedTransaction(tx, k.masterV, k.confirmValids), nil
	}
	tx = NewChainedTransaction(tx, *k.prevTx)
	return NewMinedTransaction(tx, k.masterV, k.confirmValids), nil
}
