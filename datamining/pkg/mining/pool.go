package mining

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//PoolRequester define methods to send request on pool nodes
type PoolRequester interface {
	RequestLock(lockPool datamining.Pool, lock lock.TransactionLock, sig string) error
	RequestUnlock(lockPool datamining.Pool, lock lock.TransactionLock, sig string) error
	RequestValidations(vPool datamining.Pool, txHash string, data interface{}, txType TransactionType) ([]datamining.Validation, error)
	RequestStorage(sPool datamining.Pool, data interface{}, end datamining.Endorsement, txType TransactionType) error
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {
	FindLastValidationPool(addr string) (datamining.Pool, error)
	FindStoragePool() (datamining.Pool, error)
}
