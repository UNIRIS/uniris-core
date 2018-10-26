package leading

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

//Pool represents a pool containing peers
type Pool struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

//PoolFinder defines methods to find miners to perform validation
type PoolFinder interface {
	FindValidationPool() (Pool, error)
	FindStoragePool() (Pool, error)
	FindLastValidationPool(addr string) (Pool, error)
}

//PoolDispatcher defines methods to contact peers in the pool
type PoolDispatcher interface {
	RequestLock(lastValidPool Pool, txLock validating.TransactionLock, sig string) error
	RequestUnlock(lastValidPool Pool, txLock validating.TransactionLock, sig string) error

	WalletPoolDispatcher
}

//WalletPoolDispatcher defines methods to contact peers about the wallet
type WalletPoolDispatcher interface {
	RequestLastWallet(storagePool Pool, addr string) ([]*datamining.Wallet, error)
	RequestWalletValidation(validPool Pool, d *datamining.WalletData) ([]datamining.Validation, error)
	RequestBioValidation(validPool Pool, b *datamining.BioData) ([]datamining.Validation, error)
	RequestWalletStorage(storagePool Pool, wallet *datamining.Wallet) error
	RequestBioStorage(storagePool Pool, bio *datamining.BioWallet) error
}
