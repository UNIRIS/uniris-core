package adding

import (
	uniris "github.com/uniris/uniris-core/pkg"
)

//LockRepository provides access to the lock persistence
type LockRepository interface {
	//StoreLock stores a lock
	StoreLock(l uniris.Lock) error

	//RemoveLock remove an existing lock
	RemoveLock(l uniris.Lock) error
}

//TransactionRepository provides access to the transaction persistence
type TransactionRepository interface {
	StoreKeychain(kc uniris.Keychain) error
	StoreID(id uniris.ID) error
	StoreKO(tx uniris.Transaction) error
}

//SharedRepository provides access to the shared data persistence
type SharedRepository interface {
	//StoreSharedEmitterKeyPair stores a shared emitter keypair
	StoreSharedEmitterKeyPair(kp uniris.SharedKeys) error
}
