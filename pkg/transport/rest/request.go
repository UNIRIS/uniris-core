package rest

type accountCreationRequest struct {
	EncryptedID       string `json:"encrypted_id" binding:"required"`
	EncryptedKeychain string `json:"encrypted_keychain" binding:"required"`
	Signature         string `json:"signature,omitempty" binding:"required"`
}
