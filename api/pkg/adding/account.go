package adding

//AccountCreationResult represents the result of the account creation
type AccountCreationResult struct {
	Transactions AccountCreationTransactions `json:"transactions" binding:"required"`
	Signature    string                      `json:"signature" binding:"required"`
}

//AccountCreationTransactions represents the generated transactions during the account creation
type AccountCreationTransactions struct {
	Biod     string `json:"biod" binding:"required"`
	Keychain string `json:"keychain" binding:"required"`
}

//AccountCreationRequest represents the data to provide to create an account
type AccountCreationRequest struct {
	EncryptedBioData      string     `json:"encrypted_bio_data" binding:"required"`
	EncryptedKeychainData string     `json:"encrypted_keychain_data" binding:"required"`
	SignaturesBio         Signatures `json:"signatures_bio" binding:"required"`
	SignaturesKeychain    Signatures `json:"signatures_keychain" binding:"required"`
	SignatureRequest      string     `json:"signature_request" binding:"required"`
}

//AccountCreationData represents the data without signature request
type AccountCreationData struct {
	EncryptedBioData      string     `json:"encrypted_bio_data"`
	EncryptedKeychainData string     `json:"encrypted_keychain_data"`
	SignaturesBio         Signatures `json:"signatures_bio"`
	SignaturesKeychain    Signatures `json:"signatures_keychain"`
}

//Signatures represents a set of signatures for the sent data
type Signatures struct {
	BiodSig   string `json:"biod_sig" binding:"required"`
	PersonSig string `json:"person_sig" binding:"required"`
}
