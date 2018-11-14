package mock

import (
	"sort"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	account_adding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	account_listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	biod_listing "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Database mock the entire database
type Database interface {
	account_adding.Repository
	account_listing.Repository
	biod_listing.Repository
	lock.Repository
}

//NewDatabase create new mocked database
func NewDatabase() Database {
	return &mockDatabase{}
}

type mockDatabase struct {
	Biometrics   []account.Biometric
	KOBiometrics []account.Biometric
	Keychains    []account.Keychain
	KOKeychains  []account.Keychain
	Locks        []lock.TransactionLock
}

func (d mockDatabase) FindBiometric(hash string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			return b, nil
		}
	}
	return nil, nil
}

func (d mockDatabase) FindLastKeychain(addr string) (account.Keychain, error) {
	sort.Slice(d.Keychains, func(i, j int) bool {
		iTimestamp := d.Keychains[i].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		jTimestamp := d.Keychains[j].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, b := range d.Keychains {
		if b.Address() == addr {
			return b, nil
		}
	}
	return nil, nil
}

func (d mockDatabase) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

func (d *mockDatabase) StoreKeychain(k account.Keychain) error {
	d.Keychains = append(d.Keychains, k)
	return nil
}

func (d *mockDatabase) StoreKOKeychain(k account.Keychain) error {
	d.KOKeychains = append(d.KOKeychains, k)
	return nil
}

func (d *mockDatabase) StoreBiometric(b account.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	return nil
}

func (d *mockDatabase) StoreKOBiometric(b account.Biometric) error {
	d.KOBiometrics = append(d.KOBiometrics, b)
	return nil
}

func (d *mockDatabase) NewLock(txLock lock.TransactionLock) error {
	d.Locks = append(d.Locks, txLock)
	return nil
}

func (d *mockDatabase) RemoveLock(txLock lock.TransactionLock) error {
	pos := d.findLockPosition(txLock)
	if pos > -1 {
		d.Locks = append(d.Locks[:pos], d.Locks[pos+1:]...)
	}
	return nil
}

func (d mockDatabase) ContainsLock(txLock lock.TransactionLock) bool {
	for _, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (d mockDatabase) findLockPosition(txLock lock.TransactionLock) int {
	for i, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
