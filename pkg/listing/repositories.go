package listing

import (
	uniris "github.com/uniris/uniris-core/pkg"
)

type TransactionRepository interface {
	FindPendingTransaction(txHash string) (*uniris.Transaction, error)
	FindKOTransaction(txHash string) (*uniris.Transaction, error)

	FindKeychainByHash(txHash string) (*uniris.Keychain, error)
	FindKeychainByAddress(addr string) (*uniris.Keychain, error)

	FindIDByHash(txHash string) (*uniris.ID, error)
	FindIDByAddress(addr string) (*uniris.ID, error)
}

type LockRepository interface {
	//ContainsLocks determines if a lock exists or not
	ContainsLock(uniris.Lock) (bool, error)
}

type SharedRepository interface {
	//ListSharedEmitterKeyPairs gets the shared emitter keypair
	ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error)
}
