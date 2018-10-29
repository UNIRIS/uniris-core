package adding

//EnrollmentResult represents the result of an enrollment
type EnrollmentResult struct {
	Transactions EnrollmentTransactions `json:"transactions" binding:"required"`
	Signature    string                 `json:"signature" binding:"required"`
}

//EnrollmentTransactions represents the generated transactions during the enrollment
type EnrollmentTransactions struct {
	Biod string `json:"biod" binding:"required"`
	Data string `json:"data" binding:"required"`
}

//EnrollmentRequest represents the data to provide to enroll an user
type EnrollmentRequest struct {
	EncryptedBioData    string     `json:"encrypted_bio_data" binding:"required"`
	EncryptedWalletData string     `json:"encrypted_wal_data" binding:"required"`
	SignaturesBio       Signatures `json:"signatures_bio" binding:"required"`
	SignaturesWallet    Signatures `json:"signatures_wal" binding:"required"`
	SignatureRequest    string     `json:"signature_request" binding:"required"`
}

//EnrollmentData represents the data without signature request
type EnrollmentData struct {
	EncryptedBioData    string     `json:"encrypted_bio_data"`
	EncryptedWalletData string     `json:"encrypted_wal_data"`
	SignaturesBio       Signatures `json:"signatures_bio"`
	SignaturesWallet    Signatures `json:"signatures_wal"`
}

//Signatures represents a set of signatures for the sent data
type Signatures struct {
	BiodSig   string `json:"biod_sig" binding:"required"`
	PersonSig string `json:"person_sig" binding:"required"`
}
