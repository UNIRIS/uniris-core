package pooling

import (
	uniris "github.com/uniris/uniris-core/pkg"
)

//PoolRequester define methods to send request inside a pool
type PoolRequester interface {

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, addr string) (uniris.Transaction, error)

	//RequestTransactionLock asks a pool to lock a transaction using the address related
	RequestTransactionLock(pool Pool, txHash string, address string) error

	//RequestTransactionUnlock asks a pool to unlock a transaction using the address related
	RequestTransactionUnlock(pool Pool, txHash string, address string) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx uniris.Transaction, validChan chan<- uniris.MinerValidation)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, tx uniris.Transaction, ackChan chan<- bool)
}
