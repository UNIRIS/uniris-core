package chain

import (
	"encoding/hex"
	"errors"
)

//ID represents a ID transaction
type ID struct {
	encAddrByRobot string
	encAddrByID    string
	encAesKey      string
	Transaction
}

//NewID creates a new ID transaction
func NewID(tx Transaction) (ID, error) {

	if tx.txType != IDTransactionType {
		return ID{}, errors.New("transaction: invalid type of transaction")
	}

	addrRobot, exist := tx.data["encrypted_address_by_robot"]
	if !exist {
		return ID{}, errors.New("transaction: missing data ID 'encrypted_address_by_robot'")
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

	if _, err := hex.DecodeString(addrRobot); err != nil {
		return ID{}, errors.New("transaction: id encrypted address for robot is not in hexadecimal format")
	}

	return ID{
		encAddrByID:    addrID,
		encAddrByRobot: addrRobot,
		encAesKey:      aesKey,
		Transaction:    tx,
	}, nil
}

//EncryptedAddrByRobot returns the encrypted keychain address with the robot public key
func (id ID) EncryptedAddrByRobot() string {
	return id.encAddrByRobot
}

//EncryptedAddrByID returns the encrypted keychain address with the ID public key
func (id ID) EncryptedAddrByID() string {
	return id.encAddrByID
}

//EncryptedAESKey returns the encrypted AES key with the ID public key
func (id ID) EncryptedAESKey() string {
	return id.encAesKey
}
