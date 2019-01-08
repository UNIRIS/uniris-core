package listing

//TransactionStatus represents the status for the transaction endorsement
type TransactionStatus int

const (

	//TransactionPending represents a pending status for the transaction
	TransactionPending TransactionStatus = 0

	//TransactionSuccess represents a success status for the transaction
	TransactionSuccess TransactionStatus = 1

	//TransactionFailure represents a failure status for the transaction
	TransactionFailure TransactionStatus = 2

	//TransactionUnknown represents a unknown status for the transaction
	TransactionUnknown TransactionStatus = 3
)

func (s TransactionStatus) String() string {
	switch s {
	case TransactionPending:
		return "Pending"
	case TransactionSuccess:
		return "Success"
	case TransactionFailure:
		return "Failure"
	case TransactionUnknown:
		return "Unknown"
	}

	return ""
}
