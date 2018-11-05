package lock

//TransactionLock represents lock data
type TransactionLock struct {
	TxHash         string
	MasterRobotKey string
	Address        string
}

//Signer define method to sign lock transaction
type Signer interface {
	SignLock(lock TransactionLock, pvKey string) (string, error)
}
