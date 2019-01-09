package mining

//AIClient define methods to communicate with AI
type AIClient interface {
	//GetMininumValidations asks the AI service to retreive the minimum of validations based on a transaction hash
	GetMininumValidations(txHash string) (int, error)

	//GetMininumReplications asks the AI service to retreive the minimum of storage replications based on a transaction hash
	GetMininumReplications(txHash string) (int, error)
}
