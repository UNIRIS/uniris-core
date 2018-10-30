package pool

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//PeerGroup represents a group of peers in a pool
type PeerGroup struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

//Requester define methods to send request on pool nodes
type Requester interface {
	RequestLock(lastValidPool PeerGroup, lock TransactionLock, sig string) error
	RequestUnlock(lastValidPool PeerGroup, lock TransactionLock, sig string) error
	RequestStorage(sPool PeerGroup, data interface{}, txType datamining.TransactionType) error
	RequestValidations(sPool PeerGroup, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error)
}

//Finder defines methods to find miners to perform validation
type Finder interface {
	FindValidationPool() (PeerGroup, error)
	FindStoragePool() (PeerGroup, error)
	FindLastValidationPool(addr string) (PeerGroup, error)
}

//TransactionLock represents a transaction lock
type TransactionLock struct {
	TxHash         string
	MasterRobotKey string
}

//Lookup find storage, validation and last validation pool
func Lookup(addr string, f Finder) (lastVPool PeerGroup, vPool PeerGroup, sPool PeerGroup, err error) {
	sPool, err = f.FindStoragePool()
	if err != nil {
		return
	}

	lastVPool, err = f.FindLastValidationPool(addr)
	if err != nil {
		return
	}

	vPool, err = f.FindValidationPool()
	if err != nil {
		return
	}

	return
}
