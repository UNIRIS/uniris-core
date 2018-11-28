package crypto

type accountCreationRequest struct {
	EncryptedID       string `json:"encrypted_id"`
	EncryptedKeychain string `json:"encrypted_keychain"`
	Signature         string `json:"signature,omitempty"`
}

type accountResult struct {
	EncryptedAESKey  string `json:"encrypted_aes_key"`
	EncryptedWallet  string `json:"encrypted_wallet"`
	EncryptedAddress string `json:"encrypted_address"`
	Signature        string `json:"signature,omitempty"`
}

type accountCreationResult struct {
	Transactions accountCreationTransactionsResult `json:"transactions" binding:"required"`
	Signature    string                            `json:"signature,omitempty" binding:"required"`
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
