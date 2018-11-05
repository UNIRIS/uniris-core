package mining

import (
	"net"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Pool represents a pool of peers
type Pool interface {
	Peers() []Peer
}

type pool struct {
	peers []Peer
}

func (p pool) Peers() []Peer {
	return p.peers
}

//NewPool creates a new pool
func NewPool(pp ...Peer) Pool {
	return pool{pp}
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {
	FindLastValidationPool(addr string) (Pool, error)
	FindStoragePool() (Pool, error)
}

//PoolRequester define methods to send request on pool nodes
type PoolRequester interface {
	RequestLock(lockPool Pool, lock lock.TransactionLock, sig string) error
	RequestUnlock(lockPool Pool, lock lock.TransactionLock, sig string) error
	RequestValidations(vPool Pool, data interface{}, txType TransactionType) ([]datamining.Validation, error)
	RequestStorage(sPool Pool, data interface{}, end datamining.Endorsement, txType TransactionType) error
}
