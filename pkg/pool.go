package uniris

type Pool []PeerIdentity

//PoolRequester define methods to send request inside a pool
type PoolRequester interface {

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, addr string) (Transaction, error)

	//RequestTransactionLock asks a pool to lock a transaction using the address related
	RequestTransactionLock(pool Pool, txLock Lock) error

	//RequestTransactionUnlock asks a pool to unlock a transaction using the address related
	RequestTransactionUnlock(pool Pool, txLock Lock) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx Transaction, masterValid MasterValidation, validChan chan<- MinerValidation, replyChan chan<- bool)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, tx Transaction, ackChan chan<- bool)
}
