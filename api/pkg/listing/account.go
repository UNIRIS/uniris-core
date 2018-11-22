package listing

//AccountResult defines the account's data returned from the robot
type AccountResult struct {
	EncryptedAESKey  string `json:"encrypted_aes_key"`
	EncryptedWallet  string `json:"encrypted_wallet"`
	EncryptedAddress string `json:"encrypted_address"`
	Signature        string `json:"signature,omitempty"`
}
