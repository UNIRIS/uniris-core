package checkers

//BioData represents the data by will be encrypted by the biometric device
type BioData struct {
	PersonPublicKey     string `json:"person_pubk"`
	BIODPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}

//WalletData represents the data will be encrypted by the person
type WalletData struct {
	PersonPublicKey    string `json:"person_pubk"`
	BIODPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}
