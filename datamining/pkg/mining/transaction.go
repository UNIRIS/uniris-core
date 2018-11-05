package mining

//TransactionType represents the transaction type
type TransactionType int

const (
	//CreateKeychainTransaction represents a wallet creation transaction
	CreateKeychainTransaction TransactionType = 0

	//CreateBioTransaction represents a bio creation transaction
	CreateBioTransaction TransactionType = 1
)
