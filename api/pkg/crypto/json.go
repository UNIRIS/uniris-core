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
	Transactions accountCreationTransactionsResult `json:"transactions"`
	Signature    string                            `json:"signature,omitempty"`
}

type accountCreationTransactionsResult struct {
	ID       transactionResult `json:"id"`
	Keychain transactionResult `json:"keychain"`
}

type transactionResult struct {
	TransactionHash string `json:"transaction_hash"`
	MasterPeerIP    string `json:"master_peer_ip"`
	Signature       string `json:"signature,omitempty"`
}

type contractCreationRequest struct {
	Address      string `json:"address"`
	Code         string `json:"code"`
	Event        string `json:"event"`
	PublicKey    string `json:"public_key"`
	Signature    string `json:"signature"`
	EmSig        string `json:"em_signature"`
	ReqSignature string `json:"request_signature,omitempty"`
}
