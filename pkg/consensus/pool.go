package consensus

import (
	"net"

	"github.com/uniris/uniris-core/pkg/chain"
)

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionLock asks a pool to lock a transaction using the address related
	RequestTransactionLock(pool Pool, txHash string, txAddr string, masterPublicKey string) error

	//RequestTransactionUnlock asks a pool to unlock a transaction using the address related
	RequestTransactionUnlock(pool Pool, txHash string, txAddr string) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.MinerValidation, error)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, minStorage int, tx chain.Transaction) error

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error)
}

//FindMasterPeers finds a list of master peer from a transaction hash
//TODO: To implement with AI algorithms
func FindMasterPeers(txHash string) (Pool, error) {
	return Pool{
		PoolMember{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func FindStoragePool(address string) (Pool, error) {
	return Pool{
		PoolMember{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func FindValidationPool(tx chain.Transaction) (Pool, error) {
	return Pool{
		PoolMember{
			ip:   net.ParseIP("127.0.0.1"),
			port: 5000,
		},
	}, nil
}

//Pool represent a pool either for sharding or validation
type Pool []PoolMember

//PoolMember represent a peer member of a pool
type PoolMember struct {
	ip   net.IP
	port int
	pubK string
}

//IP returns the pool member IP addres
func (pm PoolMember) IP() net.IP {
	return pm.ip
}

//Port returns the pool member port
func (pm PoolMember) Port() int {
	return pm.port
}

//PublicKey returns the pool member public key
func (pm PoolMember) PublicKey() string {
	return pm.pubK
}
