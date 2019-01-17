package contract

type Hasher interface {
	HashEndorsedContract(EndorsedContract) (string, error)
	HashContract(Contract) (string, error)
}

type SignatureVerifier interface {
	VerifyContractSignature(Contract) error
}
