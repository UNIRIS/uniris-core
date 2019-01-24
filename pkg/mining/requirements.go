package mining

//GetMinimumTransactionValidation returns the validation from a transaction hash
func GetMinimumTransactionValidation(txHash string) int {
	return 1
}

func getMinimumReplicas(txHash string) int {
	return 1
}
