package adding

//EnrollmentResult represents the result of an enrollment
type EnrollmentResult struct {
	Hash             string `json:"hash" binding:"required"`
	SignatureRequest string `json:"signature_request" binding:"required"`
}

//EnrollmentRequest represents the data to provide to enroll an user
type EnrollmentRequest struct {
	EncryptedBioData    string     `json:"encrypted_bio_data" binding:"required"`
	EncryptedWalletData string     `json:"encrypted_wal_data" binding:"required"`
	SignaturesBio       Signatures `json:"signatures_bio" binding:"required"`
	SignaturesWallet    Signatures `json:"signatures_wal" binding:"required"`
	SignatureRequest    string     `json:"signature_request" binding:"required"`
}

//EnrollmentVerifyRequest represents the data to verify before to enroll
type EnrollmentVerifyRequest struct {
	EncryptedBioData    string     `json:"encrypted_bio_data"`
	EncryptedWalletData string     `json:"encrypted_wal_data"`
	SignaturesBio       Signatures `json:"signatures_bio"`
	SignaturesWallet    Signatures `json:"signatures_wal"`
	SignatureRequest    string     `json:"signature_request"`
}

//Signatures represents a set of signatures for the sent data
type Signatures struct {
	BiodSig   string `json:"biod_sig" binding:"required"`
	PersonSig string `json:"person_sig" binding:"required"`
}
