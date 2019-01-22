package inspecting

//IsAuthorizedToStoreTx checks if the miner can store the transaction
func IsAuthorizedToStoreTx(txHash string) bool {
	return true
}

//GetMinimumTransactionValidation returns the validation from a transaction hash
func GetMinimumTransactionValidation(txHash string) int {
	return 1
}

func FindTransactionMasterPeer(txHash string) string {
	return "127.0.0.1"
}
