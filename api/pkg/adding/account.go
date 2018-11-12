package adding

//AccountCreationResult represents the result of the account creation
type AccountCreationResult struct {
	Transactions AccountCreationTransactions `json:"transactions" binding:"required"`
	Signature    string                      `json:"signature" binding:"required"`
}

//AccountCreationTransactions represents the transactions for the account creation
type AccountCreationTransactions struct {
	Biometric TransactionResult `json:"biometric" binding:"required"`
	Keychain  TransactionResult `json:"keychain" binding:"required"`
}

//TransactionResult represents the result for a transaction
type TransactionResult struct {
	TransactionHash string `json:"transaction_hash" binding:"required"`
	MasterPeerIP    string `json:"master_peer_ip" binding:"required"`
	Signature       string `json:"signature" binding:"required"`
}

//AccountCreationRequest represents the required data to create an account
type AccountCreationRequest struct {
	EncryptedBioData      string     `json:"encrypted_bio_data" binding:"required"`
	EncryptedKeychainData string     `json:"encrypted_keychain_data" binding:"required"`
	SignaturesBio         Signatures `json:"signatures_bio" binding:"required"`
	SignaturesKeychain    Signatures `json:"signatures_keychain" binding:"required"`
	SignatureRequest      string     `json:"signature_request" binding:"required"`
}

//Signatures represents a common set of signatures for requests
type Signatures struct {
	BiodSig   string `json:"biod_sig" binding:"required"`
	PersonSig string `json:"person_sig" binding:"required"`
}
