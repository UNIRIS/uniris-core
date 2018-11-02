package mem

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	account_adding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	account_listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	biod_listing "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Repo mock the entire database
type Repo interface {
	account_adding.Repository
	account_listing.Repository
	biod_listing.Repository
	lock.Repository
}

type databasemock struct {
	Biometrics []account.Biometric
	Keychains  []account.Keychain
	Locks      []datamining.TransactionLock
}

//NewDatabase creates a new mock database
func NewDatabase() Repo {
	return &databasemock{
		Biometrics: make([]account.Biometric, 0),
		Keychains:  make([]account.Keychain, 0),
		Locks:      make([]datamining.TransactionLock, 0),
	}
}

func (d *databasemock) FindBiometric(hash string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindKeychain(addr string) (account.Keychain, error) {
	for _, w := range d.Keychains {
		if w.WalletAddr() == addr {
			return w, nil
		}
	}
	return nil, nil
}

func (d *databasemock) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

func (d *databasemock) StoreKeychain(w account.Keychain) error {
	d.Keychains = append(d.Keychains, w)
	return nil
}

func (d *databasemock) StoreBiometric(b account.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	return nil
}

func (d *databasemock) NewLock(txLock datamining.TransactionLock) error {
	d.Locks = append(d.Locks, txLock)
	return nil
}

func (d *databasemock) RemoveLock(txLock datamining.TransactionLock) error {
	pos := d.findLockPosition(txLock)
	if pos > -1 {
		d.Locks = append(d.Locks[:pos], d.Locks[pos+1:]...)
	}
	return nil
}

func (d databasemock) ContainsLock(txLock datamining.TransactionLock) bool {
	for _, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (d databasemock) findLockPosition(txLock datamining.TransactionLock) int {
	for i, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
