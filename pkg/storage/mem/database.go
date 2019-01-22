package mem

import (
	"sort"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
)

//Database is the memory database
type Database interface {
	listing.Repository
	adding.Repository
}

type db struct {
	keychains  []uniris.Keychain
	ids        []uniris.ID
	sharedKeys []uniris.SharedKeys
	koTxs      []uniris.Transaction
	pendingTxs []uniris.Transaction
	locks      []uniris.Lock

	Database
}

//NewDatabase creates a new memory database
func NewDatabase() Database {
	return &db{}
}

func (d db) ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error) {
	return d.sharedKeys, nil
}

func (d db) FindPendingTransaction(txHash string) (*uniris.Transaction, error) {
	for _, tx := range d.pendingTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d db) FindKeychainByHash(txHash string) (*uniris.Keychain, error) {
	for _, tx := range d.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d db) FindKeychainByAddress(addr string) (*uniris.Keychain, error) {

	sort.Slice(d.keychains, func(i, j int) bool {
		iTimestamp := d.keychains[i].Timestamp().Unix()
		jTimestamp := d.keychains[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range d.keychains {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d db) FindIDByHash(txHash string) (*uniris.ID, error) {
	for _, tx := range d.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d db) FindIDByAddress(addr string) (*uniris.ID, error) {

	sort.Slice(d.ids, func(i, j int) bool {
		iTimestamp := d.ids[i].Timestamp().Unix()
		jTimestamp := d.ids[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range d.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d db) FindKOTransaction(txHash string) (*uniris.Transaction, error) {
	for _, tx := range d.koTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d *db) StoreSharedEmitterKeyPair(sk uniris.SharedKeys) error {
	d.sharedKeys = append(d.sharedKeys, sk)
	return nil
}

func (d *db) StoreKeychain(kc uniris.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *db) StoreID(id uniris.ID) error {
	d.ids = append(d.ids, id)
	return nil
}

func (d *db) StoreKO(tx uniris.Transaction) error {
	d.koTxs = append(d.koTxs, tx)
	return nil
}

func (d *db) StoreLock(l uniris.Lock) error {
	d.locks = append(d.locks, l)
	return nil
}
func (d *db) RemoveLock(l uniris.Lock) error {
	pos := d.findLockPosition(l)
	if pos > -1 {
		d.locks = append(d.locks[:pos], d.locks[pos+1:]...)
	}
	return nil
}
func (d *db) ContainsLock(l uniris.Lock) (bool, error) {
	return d.findLockPosition(l) > -1, nil
}

func (d db) findLockPosition(txLock uniris.Lock) int {
	for i, lock := range d.locks {
		if lock.TransactionHash() == txLock.TransactionHash() && txLock.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}
