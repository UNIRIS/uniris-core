package slave

//TransactionLock represents lock data
type TransactionLock struct {
	TxHash         string `json:"tx_hash"`
	MasterRobotKey string `json:"master_robot_key"`
}

//Locker defines methods to manage locks
type Locker interface {
	Lock(TransactionLock) error
	Unlock(TransactionLock) error
	ContainsLock(TransactionLock) bool
}
