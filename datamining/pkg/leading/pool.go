package leading

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//Pool represents a pool containing peers
type Pool struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	Port      int
	PublicKey string
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {
	FindValidationPool() (Pool, error)
	FindStoragePool() (Pool, error)
}

//PoolDispatcher defines methods to contact peers in the pool
type PoolDispatcher interface {
	RequestLock(pool Pool, txHash string) error
	RequestUnlock(pool Pool, txHash string) error
	RequestWalletValidation(Pool, *datamining.WalletData) ([]datamining.Validation, error)
	RequestBioValidation(Pool, *datamining.BioData) ([]datamining.Validation, error)

	RequestLastTx(pool Pool, txHash string) (oldTxHash string, validation *datamining.MasterValidation, err error)
	RequestWalletStorage(Pool, *datamining.Wallet) error
	RequestBioStorage(Pool, *datamining.BioWallet) error
}
