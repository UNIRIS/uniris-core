package transaction

import (
	"encoding/hex"
	"errors"
)

//IDRepository manages the ID transaction storage
type IDRepository interface {
	FindIDByHash(txHash string) (*ID, error)
	FindIDByAddress(addr string) (*ID, error)
	StoreID(id ID) error
}

//ID represents a ID transaction
type ID struct {
	encAddrByRobot string
	encAddrByID    string
	encAesKey      string
	Transaction
}

//NewID creates an ID transaction by extracting the transaction data
func NewID(tx Transaction) (ID, error) {

	if tx.Type() != IDType {
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

//ToTransaction converts back ID to transaction root
func (id ID) ToTransaction() (tx Transaction, err error) {
	tx, err = New(id.Address(), IDType, map[string]string{
		"encrypted_address_by_robot": id.encAddrByRobot,
		"encrypted_address_by_id":    id.encAddrByID,
		"encrypted_aes_key":          id.encAesKey,
	}, id.Timestamp(), id.PublicKey(), id.Signature(), id.EmitterSignature(), id.Proposal(), id.TransactionHash())
	if err != nil {
		return
	}
	if id.masterV.pow != "" && len(id.confirmValids) > 0 {
		if err = tx.AddMining(id.MasterValidation(), id.ConfirmationsValidations()); err != nil {
			return
		}
	}
	return tx, nil
}
