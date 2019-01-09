package uniris

//ID represents a ID transaction
type ID struct {
	hash           string
	encAddrByRobot string
	encAddrByID    string
	encAesKey      string
	Transaction
}

//NewID creates a new transaction related to an ID
func NewID(hash string, encAddrByRobot string, encAddrByID string, encAesKey string, rootTx Transaction) ID {
	return ID{
		hash:           hash,
		encAddrByRobot: encAddrByRobot,
		encAddrByID:    encAddrByID,
		encAesKey:      encAesKey,
		Transaction:    rootTx,
	}
}

//Hash return the ID hash (address)
func (id ID) Hash() string {
	return id.hash
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
