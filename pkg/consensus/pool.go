package consensus

import (
	"encoding/hex"
	"net"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionTimeLock asks a pool to timelock a transaction using the address related
	RequestTransactionTimeLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, minStorage int, tx chain.Transaction) error

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error)
}

//FindMasterNodes finds a list of master node from a transaction hash
//TODO: To implement with AI algorithms
func FindMasterNodes(txHash crypto.VersionnedHash, txType chain.TransactionType) (Pool, error) {

	b, err := hex.DecodeString("0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	if err != nil {
		return Pool{}, err
	}
	pub, err := crypto.ParsePublicKey(b)

	return Pool{
		Node{
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
			publicKey: pub,
		},
	}, nil
}

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func FindStoragePool(address []byte) (Pool, error) {
	b, err := hex.DecodeString("0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	if err != nil {
		return Pool{}, err
	}
	pub, err := crypto.ParsePublicKey(b)

	return Pool{
		Node{
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
			publicKey: pub,
		},
	}, nil
}

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func FindValidationPool(tx chain.Transaction) (Pool, error) {
	b, err := hex.DecodeString("0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	if err != nil {
		return Pool{}, err
	}
	pub, err := crypto.ParsePublicKey(b)
	if err != nil {
		return Pool{}, err
	}

	return Pool{
		Node{
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
			publicKey: pub,
		},
	}, nil
}

//Pool represent a pool either for sharding or validation
type Pool []Node
