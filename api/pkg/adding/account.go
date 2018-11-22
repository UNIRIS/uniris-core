package adding

//AccountCreationResult represents the result of the account creation
type AccountCreationResult struct {
	Transactions AccountCreationTransactionsResult `json:"transactions" binding:"required"`
	Signature    string                            `json:"signature,omitempty" binding:"required"`
}

//AccountCreationTransactionsResult represents the transactions for the account creation
type AccountCreationTransactionsResult struct {
	ID       TransactionResult `json:"id" binding:"required"`
	Keychain TransactionResult `json:"keychain" binding:"required"`
}

//TransactionResult represents the result for a transaction
type TransactionResult struct {
	TransactionHash string `json:"transaction_hash" binding:"required"`
	MasterPeerIP    string `json:"master_peer_ip" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

//ProposedKeyPair represent a key pair for a renew proposal
type ProposedKeyPair struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}

//AccountCreationRequest represents the required data to create an account
type AccountCreationRequest struct {
	EncryptedID       string `json:"encrypted_id"`
	EncryptedKeychain string `json:"encrypted_keychain"`
	Signature         string `json:"signature,omitempty"`
}
