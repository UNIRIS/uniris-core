package checks

//Handler defines methods every checker has to define
type Handler interface {
	CheckData(data interface{}, txHash string) error
}

//TransactionDataHasher define methods to hash transaction data
type TransactionDataHasher interface {
	HashTransactionData(data interface{}) (string, error)
}

//Signer defines methods to handle signatures
type Signer interface {
	CheckSignature(pubKey string, data interface{}, sig string) error
}

type rawKeychainData struct {
	PersonPublicKey    string `json:"person_pubk"`
	BiodPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

type rawBiometricData struct {
	PersonPublicKey     string `json:"person_pubk"`
	BiodPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}
