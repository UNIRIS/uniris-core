package datamining

import (
	"net"
)

//Pool represents a pool of peers
type Pool struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {
	FindLastValidationPool(addr string) (Pool, error)
	FindValidationPool() (Pool, error)
	FindStoragePool() (Pool, error)
}

//PoolRequester define methods to send request on pool nodes
type PoolRequester interface {
	RequestLock(lockPool Pool, lock TransactionLock, sig string) error
	RequestUnlock(lockPool Pool, lock TransactionLock, sig string) error
	RequestStorage(sPool Pool, data interface{}, endorsement Endorsement, txType TransactionType) error
	RequestValidations(sPool Pool, data interface{}, txType TransactionType) ([]Validation, error)
}
