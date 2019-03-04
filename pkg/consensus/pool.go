package consensus

import (
	"net"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
<<<<<<< HEAD
	//RequestTransactionTimeLock asks a pool to lock a transaction using the address related
	RequestTransactionTimeLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error
=======
	//RequestTransactionLock asks a pool to lock a transaction using the address related
	RequestTransactionLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash

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

	pub, err := crypto.ParsePublicKey([]byte("3059301306072a8648ce3d020106082a8648ce3d0301070342000408f4b4026d2560aaa552244bdf8ec421bb41378b56487d9d4ca5a57fd6e64ef7ae2f2c6530f18bd0f359342b4fa7fdaeaa60c45a1197260eb1c267cc996bec81"))
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

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func FindStoragePool(address []byte) (Pool, error) {
	pub, err := crypto.ParsePublicKey([]byte("3059301306072a8648ce3d020106082a8648ce3d0301070342000408f4b4026d2560aaa552244bdf8ec421bb41378b56487d9d4ca5a57fd6e64ef7ae2f2c6530f18bd0f359342b4fa7fdaeaa60c45a1197260eb1c267cc996bec81"))
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

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func FindValidationPool(tx chain.Transaction) (Pool, error) {
	pub, err := crypto.ParsePublicKey([]byte("3059301306072a8648ce3d020106082a8648ce3d0301070342000408f4b4026d2560aaa552244bdf8ec421bb41378b56487d9d4ca5a57fd6e64ef7ae2f2c6530f18bd0f359342b4fa7fdaeaa60c45a1197260eb1c267cc996bec81"))
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
<<<<<<< HEAD
type Pool []Node
=======
type Pool []PoolMember

//PoolMember represent a node member of a pool
type PoolMember struct {
	ip   net.IP
	port int
	pubK crypto.PublicKey
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
func (pm PoolMember) PublicKey() crypto.PublicKey {
	return pm.pubK
}
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash
