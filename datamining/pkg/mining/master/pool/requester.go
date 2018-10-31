package pool

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/locking"
)

//Requester define methods to send request on pool nodes
type Requester interface {
	RequestLock(lockPool PeerGroup, lock locking.TransactionLock, sig string) error
	RequestUnlock(lockPool PeerGroup, lock locking.TransactionLock, sig string) error
	RequestStorage(sPool PeerGroup, data interface{}, txType datamining.TransactionType) error
	RequestValidations(sPool PeerGroup, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error)
}
