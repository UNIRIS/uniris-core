package contract

type Hasher interface {
	HashEndorsedContract(EndorsedContract) (string, error)
	HashContract(Contract) (string, error)
	HashContractMessage(Message) (string, error)
	HashEndorsedContractMessage(EndorsedMessage) (string, error)
}

type SignatureVerifier interface {
	VerifyContractSignature(Contract) error
	VerifyContractMessageSignature(Message) error
}
