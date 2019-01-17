package rest

type accountRequest struct {
	EncryptedID       string `json:"encrypted_id" binding:"required"`
	EncryptedKeychain string `json:"encrypted_keychain" binding:"required"`
	Signature         string `json:"signature" binding:"required"`
}

type accountCreationResult struct {
	Transactions accountCreationTransactionsResult `json:"transactions" binding:"required"`
	Signature    string                            `json:"signature" binding:"required"`
}

type accountCreationTransactionsResult struct {
	ID       transactionResult `json:"id" binding:"required"`
	Keychain transactionResult `json:"keychain" binding:"required"`
}

type transactionResult struct {
	TransactionHash string `json:"transaction_hash" binding:"required"`
	MasterPeerIP    string `json:"master_peer_ip" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

type accountResult struct {
	EncryptedAESKey  string `json:"encrypted_aes_key" binding:"required"`
	EncryptedWallet  string `json:"encrypted_wallet" binding:"required"`
	EncryptedAddress string `json:"encrypted_address" binding:"required"`
	Signature        string `json:"signature" binding:"required"`
}

type sharedEmitterKeys struct {
	PublicKey           string `json:"public_key" binding:"required"`
	EncryptedPrivateKey string `json:"encrypted_private_key" binding:"required"`
}

type sharedKeys struct {
	SharedRobotPublicKey string              `json:"shared_robot_pubkey" binding:"required"`
	SharedEmitterKeys    []sharedEmitterKeys `json:"shared_emitter_keys" binding:"required"`
}

type contractCreationRequest struct {
	EncryptedContract string `json:"encrypted_contract"`
	Signature         string `json:"signature"`
}

type contractMessageRequest struct {
	EncryptedMessage string `json:"encrypted_message"`
	Signature        string `json:"signature"`
}

type contractState struct {
	Data      string `json:"data"`
	Signature string `json:"signature"`
}
