package master

//TransactionStatus defines the status of a transaction
type TransactionStatus int

const (

	//Pending is defined when the is received but neither locked or approved
	Pending TransactionStatus = 0

	//Locked is defined when the is received and locked
	Locked TransactionStatus = 1

	//Unlocked is defined when the is received and unlocked
	Unlocked TransactionStatus = 2

	//Approved is defined when the is received, locked and validations passed
	Approved TransactionStatus = 3

	//Invalid is defined when the is received, locked and one validation failed
	Invalid TransactionStatus = 4

	//Replicated is defined when the is received, locked, valid and replicated
	Replicated TransactionStatus = 5
)

func (s TransactionStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Locked:
		return "Locked"
	case Unlocked:
		return "Unlocked"
	case Approved:
		return "Approved"
	case Invalid:
		return "Invalid"
	case Replicated:
		return "Replicated"
	default:
		panic("Unrecognized transactions status")
	}
}

//Notifier defines methods to notify transactions statuses
type Notifier interface {
	NotifyTransactionStatus(txHash string, status TransactionStatus) error
}
