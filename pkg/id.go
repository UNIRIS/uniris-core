package uniris

import (
	"encoding/json"
	"errors"
)

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
	var data idData
	if err := json.Unmarshal([]byte(tx.Data()), &data); err != nil {
		return ID{}, err
	}

	if data.EncryptedAESKey == "" || data.EncryptedAddressByID == "" || data.EncryptedAddressByRobot == "" {
		return ID{}, errors.New("Missing ID transaction data")
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

func (id ID) ToTransaction() (tx Transaction, err error) {
	data := idData{
		EncryptedAddressByID:    id.EncryptedAddrByID(),
		EncryptedAddressByRobot: id.EncryptedAddrByRobot(),
		EncryptedAESKey:         id.EncryptedAESKey(),
	}
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	tx, err = NewTransaction(id.address, IDTransactionType, string(b), id.timestamp, id.pubKey, id.sig, id.emSig, id.prop, id.txHash)
	if err != nil {
		return
	}
	if err = tx.AddMining(id.masterV, id.confirmValids); err != nil {
		return
	}
	return tx, nil
}
