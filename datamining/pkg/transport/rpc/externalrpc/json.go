package externalrpc

//KeychainJSON represents a keychain as JSON
type KeychainJSON struct {
	Data        KeychainDataJSON `json:"data"`
	Endorsement EndorsementJSON  `json:"endorsement"`
}

//BiometricJSON represents a biometric as JSON
type BiometricJSON struct {
	Data        BioDataJSON     `json:"data"`
	Endorsement EndorsementJSON `json:"endorsement"`
}

//EndorsementJSON represents a endorsement as JSON
type EndorsementJSON struct {
	LastTransactionHash string               `json:"last_transaction_hash"`
	TransactionHash     string               `json:"transaction_hash"`
	MasterValidation    MasterValidationJSON `json:"master_validation"`
	Validations         []ValidationJSON     `json:"validations"`
}

//MasterValidationJSON represents a master validation as JSON
type MasterValidationJSON struct {
	LastTransactionMiners []string       `json:"last_transaction_miners"`
	ProofOfWorkRobotKey   string         `json:"proof_of_work_robot_key"`
	ProofOfWorkValidation ValidationJSON `json:"proof_of_work_validation"`
}

//ValidationJSON represents a validation as JSON
type ValidationJSON struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

//KeychainDataJSON represents keychain data JSON
type KeychainDataJSON struct {
	PersonPublicKey    string `json:"person_pubk"`
	BiodPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

//BioDataJSON represents bio data JSON
type BioDataJSON struct {
	PersonPublicKey     string `json:"person_pubk"`
	BiodPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}

//Signatures represents signatures JSON
type Signatures struct {
	Person string `json:"person_sig"`
	Biod   string `json:"biod_sig"`
}
