package uniris

type Lock struct {
	txHash         string
	address        string
	masterRobotKey string
}

func NewLock(txHash, address, masterKey string) Lock {
	return Lock{txHash, address, masterKey}
}

func (l Lock) TransactionHash() string {
	return l.txHash
}

func (l Lock) Address() string {
	return l.address
}

func (l Lock) MasterRobotKey() string {
	return l.masterRobotKey
}
