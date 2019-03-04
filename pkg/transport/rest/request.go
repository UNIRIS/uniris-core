package rest

type accountCreationRequest struct {
	EncryptedID       string `json:"encrypted_id" validate:"required,hexadecimal"`
	EncryptedKeychain string `json:"encrypted_keychain" validate:"required,hexadecimal"`
	Signature         string `json:"signature,omitempty" validate:"required,hexadecimal"`
}

type txRaw struct {
	Address                   string            `json:"addr" validate:"required,hexadecimal"`
	Data                      map[string]string `json:"data" validate:"required,gt=0,dive,keys,endkeys,required,hexadecimal"`
	Timestamp                 int64             `json:"timestamp" validate:"required"`
	Type                      int               `json:"type" validate:"oneof=0 1"`
	PublicKey                 string            `json:"public_key" validate:"required,hexadecimal"`
	SharedKeysEmitterProposal struct {
		EncryptedPrivateKey string `json:"encrypted_private_key" validate:"required,hexadecimal"`
		PublicKey           string `json:"public_key" validate:"required,hexadecimal"`
	} `json:"em_shared_keys_proposal" validate:"required"`
	Signature        string `json:"signature" validate:"required,hexadecimal"`
	EmitterSignature string `json:"em_signature" validate:"required,hexadecimal"`
}
