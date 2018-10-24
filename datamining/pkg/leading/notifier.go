package leading

//TransactionStatus defines the status of a transaction
type TransactionStatus int

const (

	//Pending is defined when the is received but neither locked or approved
	Pending TransactionStatus = 0

	//Locked is defined when the is received and locked
	Locked TransactionStatus = 1

	//Approved is defined when the is received, locked and validations passed
	Approved TransactionStatus = 2
)

func (s TransactionStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Locked:
		return "Locked"
	case Approved:
		return "Approved"
	default:
		panic("Unrecognized transactions status")
	}
}

//Notifier defines methods to notify transactions statuses
type Notifier interface {
	NotifyTransactionStatus(txHash string, status TransactionStatus) error
}
