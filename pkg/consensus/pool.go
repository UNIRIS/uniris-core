package consensus

import (
	"net"

	"github.com/uniris/uniris-core/pkg/chain"
)

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionTimeLock asks a pool to timelock a transaction using the address related
	RequestTransactionTimeLock(pool Pool, txHash string, txAddr string, masterPublicKey string) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, minStorage int, tx chain.Transaction) error

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error)
}

//FindMasterNodes finds a list of master node from a transaction hash
//TODO: To implement with AI algorithms
func FindMasterNodes(txHash string) (Pool, error) {
	return Pool{
		Node{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func FindStoragePool(address string) (Pool, error) {
	return Pool{
		Node{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func FindValidationPool(tx chain.Transaction) (Pool, error) {
	return Pool{
		Node{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//Pool represent a pool either for sharding or validation
type Pool []Node
