package lock

import "github.com/uniris/uniris-core/datamining/pkg/mining/pool"

type Locker interface {
	RequestLock(lastValidPool pool.PeerCluster, lock TransactionLock, sig string) error
	RequestUnlock(lastValidPool pool.PeerCluster, lock TransactionLock, sig string) error
}

type TransactionLocker interface {
	Lock(TransactionLock) error
	Unlock(TransactionLock) error
	ContainsLock(TransactionLock) bool
}

type TransactionLock struct {
	TxHash         string `json:"tx_hash"`
	MasterRobotKey string `json:"master_robot_key"`
}
