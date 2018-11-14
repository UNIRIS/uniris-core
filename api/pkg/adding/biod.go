package adding

//BiodRegisterRequest represents the registering request of a biometric device
type BiodRegisterRequest struct {
	EncryptedPublicKey string `json:"encrypted_public_key"`
	Signature          string `json:"signature"`
}

//BiodRegisterResponse represents the registering response of a biometric device
type BiodRegisterResponse struct {
	PublicKeyHash string `json:"public_key_hash"`
	Signature     string `json:"signature"`
}
