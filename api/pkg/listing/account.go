package listing

//AccountResult defines the account's data returned from the robot
type AccountResult struct {
	EncryptedWallet  string `json:"encrypted_wallet"`
	EncryptedAESKey  string `json:"encrypted_aes_key"`
	AddrWalletPerson string `json:"@wallet_encrypted_person"`
}

//AccountRequest represents the data will be send to the robot
type AccountRequest struct {
	EncryptedHash    []byte
	SignatureRequest []byte
}
