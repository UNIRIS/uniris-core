package mining

//TransactionType represents the transaction type
type TransactionType int

const (
	//KeychainTransaction represents transaction related to keychain (wallet)
	KeychainTransaction TransactionType = 0

	//BiometricTransaction represents transaction related to biometric data
	BiometricTransaction TransactionType = 1
)
