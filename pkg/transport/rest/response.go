package rest

type transactionResponse struct {
	TransactionReceipt string `json:"transaction_receipt"`
	Timestamp          int64  `json:"timestamp"`
	Signature          string `json:"signature,omitempty"`
}

type accountCreationResponse struct {
	IDTransaction       transactionResponse `json:"id_transaction"`
	KeychainTransaction transactionResponse `json:"keychain_transaction"`
}

type sharedKeysResponse struct {
	NodePublicKey string              `json:"shared_node_public_key"`
	EmitterKeys   []emitterSharedKeys `json:"shared_emitter_keys"`
}

type emitterSharedKeys struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}

type accountFindResponse struct {
	EncryptedWallet string `json:"encrypted_wallet"`
	EncryptedAESKey string `json:"encrypted_aes_key"`
	Timestamp       int64  `json:"timestamp"`
	Signature       string `json:"signature,omitempty"`
}

type transactionStatusResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}
