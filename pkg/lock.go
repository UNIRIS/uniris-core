package uniris

import "errors"

type Lock struct {
	txHash         string
	address        string
	masterRobotKey string
}

func NewLock(txHash, address, masterKey string) (Lock, error) {
	if txHash == "" || address == "" || masterKey == "" {
		return Lock{}, errors.New("Missing lock information")
	}

	return Lock{txHash, address, masterKey}, nil
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
