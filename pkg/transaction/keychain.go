package transaction

import (
	"encoding/hex"
	"encoding/json"
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

type keychainData struct {
	EncryptedAddress string `json:"encrypted_address"`
	EncryptedWallet  string `json:"encrypted_wallet"`
}

//NewKeychain creates a keychain transaction by extracting the transaction data
func NewKeychain(tx Transaction) (Keychain, error) {

	if tx.Type() != KeychainType {
		return Keychain{}, errors.New("transaction: invalid type of transaction")
	}

	var data keychainData

	dataBytes, err := hex.DecodeString(tx.Data())
	if err != nil {
		return Keychain{}, err
	}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return Keychain{}, err
	}

	if data.EncryptedAddress == "" || data.EncryptedWallet == "" {
		return Keychain{}, errors.New("transaction: missing keychain transaction data")
	}

	if _, err := hex.DecodeString(data.EncryptedAddress); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted address is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(data.EncryptedWallet); err != nil {
		return Keychain{}, errors.New("transaction: keychain encrypted wallet is not in hexadecimal format")
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

//ToTransaction converts back Keychain to transaction root
func (k Keychain) ToTransaction() (Transaction, error) {
	data := keychainData{
		EncryptedAddress: k.EncryptedAddrByRobot(),
		EncryptedWallet:  k.EncryptedWallet(),
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return Transaction{}, nil
	}

	tx, err := New(k.Address(), KeychainType, hex.EncodeToString(dataBytes), k.Timestamp(), k.PublicKey(), k.Signature(), k.EmitterSignature(), k.Proposal(), k.TransactionHash())
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
