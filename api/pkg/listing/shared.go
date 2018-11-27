package listing

//SharedKeysResult describes the shared keys result used internally
type SharedKeysResult struct {
	RobotPublicKey  string          `json:"shared_robot_pubkey"`
	RobotPrivateKey string          `json:"shared_robot_privatekey,omitempty"`
	EmitterKeys     []SharedKeyPair `json:"shared_emitter_keys"`
}

//RequestPublicKey returns the public key for emitter request
func (sk SharedKeysResult) RequestPublicKey() string {
	return sk.EmitterKeys[0].PublicKey
}

//SharedKeyPair represent a shared keypair
type SharedKeyPair struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}
