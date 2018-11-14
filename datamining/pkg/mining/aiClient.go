package mining

//AIClient define methods to communicate with AI
type AIClient interface {
	GetMininumValidations(txHash string) (int, error)
}
