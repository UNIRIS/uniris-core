package transaction

import (
	"encoding/hex"
	"encoding/json"
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

type idData struct {
	EncryptedAddressByRobot string
	EncryptedAddressByID    string
	EncryptedAESKey         string
}

//NewID creates an ID transaction by extracting the transaction data
func NewID(tx Transaction) (ID, error) {

	if tx.Type() != IDType {
		return ID{}, errors.New("transaction: invalid type of transaction")
	}

	var data idData

	dataBytes, err := hex.DecodeString(tx.Data())
	if err != nil {
		return ID{}, err
	}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return ID{}, err
	}

	if data.EncryptedAESKey == "" || data.EncryptedAddressByID == "" || data.EncryptedAddressByRobot == "" {
		return ID{}, errors.New("transaction: missing id transaction data")
	}

	if _, err := hex.DecodeString(data.EncryptedAESKey); err != nil {
		return ID{}, errors.New("transaction: id encrypted aes key is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(data.EncryptedAddressByID); err != nil {
		return ID{}, errors.New("transaction: id encrypted address for id is not in hexadecimal format")
	}

	if _, err := hex.DecodeString(data.EncryptedAddressByRobot); err != nil {
		return ID{}, errors.New("transaction: id encrypted address for robot is not in hexadecimal format")
	}

	return ID{
		encAddrByID:    data.EncryptedAddressByID,
		encAddrByRobot: data.EncryptedAddressByRobot,
		encAesKey:      data.EncryptedAESKey,
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
	data := idData{
		EncryptedAddressByID:    id.EncryptedAddrByID(),
		EncryptedAddressByRobot: id.EncryptedAddrByRobot(),
		EncryptedAESKey:         id.EncryptedAESKey(),
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	tx, err = New(id.Address(), IDType, hex.EncodeToString(dataBytes), id.Timestamp(), id.PublicKey(), id.Signature(), id.EmitterSignature(), id.Proposal(), id.TransactionHash())
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
