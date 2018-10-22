package listing

//AccountResult defines the account's data returned from the robot
type AccountResult struct {
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedWallet     string `json:"encrypted_wallet"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
}

//SignedAccountResult defines the account's data returned from the robot but signed
type SignedAccountResult struct {
	EncryptedWallet     string `json:"encrypted_wallet"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	SignatureRequest    string `json:"signature_request"`
}
