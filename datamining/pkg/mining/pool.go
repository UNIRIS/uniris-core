package mining

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//PoolRequester define methods to send request on pool nodes
type PoolRequester interface {

	//RequestLock assk a pool pool to lock a transaction
	RequestLock(lockPool datamining.Pool, lock lock.TransactionLock) error

	//RequestUnlock asks a pool to unlock a transaction
	RequestUnlock(lockPool datamining.Pool, lock lock.TransactionLock) error

	//RequestValidations asks a validation pool to perform checks on a transaction
	RequestValidations(vPool datamining.Pool, txHash string, data interface{}, txType TransactionType) ([]Validation, error)

	//RequestStorage asks a storage pool to store the transaction
	RequestStorage(sPool datamining.Pool, data interface{}, end Endorsement, txType TransactionType) error
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {

	//FindLastValidationPool searches a validation pool for the given address
	FindLastValidationPool(addr string) (datamining.Pool, error)

	//FindStoragePool searches a storage pool for the given address
	FindStoragePool(addr string) (datamining.Pool, error)
}
