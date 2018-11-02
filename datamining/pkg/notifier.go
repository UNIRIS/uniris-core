package datamining

//TransactionStatus defines the status of a transaction
type TransactionStatus int

const (

	//TxPending is defined when the is received but neither locked or approved
	TxPending TransactionStatus = 0

	//TxLocked is defined when the is received and locked
	TxLocked TransactionStatus = 1

	//TxUnlocked is defined when the is received and unlocked
	TxUnlocked TransactionStatus = 2

	//TxApproved is defined when the is received, locked and validations passed
	TxApproved TransactionStatus = 3

	//TxInvalid is defined when the is received, locked and one validation failed
	TxInvalid TransactionStatus = 4

	//TxReplicated is defined when the is received, locked, valid and replicated
	TxReplicated TransactionStatus = 5
)

func (s TransactionStatus) String() string {
	switch s {
	case TxPending:
		return "Pending"
	case TxLocked:
		return "Locked"
	case TxUnlocked:
		return "Unlocked"
	case TxApproved:
		return "Approved"
	case TxInvalid:
		return "Invalid"
	case TxReplicated:
		return "Replicated"
	default:
		panic("Unrecognized transactions status")
	}
}

//Notifier defines methods to notify transactions statuses
type Notifier interface {
	NotifyTransactionStatus(txHash string, status TransactionStatus) error
}
