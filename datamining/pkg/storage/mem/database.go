package mem

import (
	"sort"

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

type database struct {
	Biometrics   []account.Biometric
	KOBiometrics []account.Biometric
	Keychains    []account.Keychain
	KOKeychains  []account.Keychain
	Locks        []lock.TransactionLock
}

//NewDatabase creates a new mock database
func NewDatabase() Repo {
	return &database{}
}

func (d *database) FindBiometric(hash string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			return b, nil
		}
	}
	return nil, nil
}

func (d *database) FindLastKeychain(addr string) (account.Keychain, error) {
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

func (d *database) ListBiodPubKeys() ([]string, error) {
	return []string{"3059301306072a8648ce3d020106082a8648ce3d03010703420004061a9f65f64c701af21bc93604b93d1502cf1f30a19a37f919cf99112d5991109ea67750a7ce5ef95054a920614aa94b33148c60f34b247de62e33a1d843be21"}, nil
}

func (d *database) StoreKeychain(w account.Keychain) error {
	d.Keychains = append(d.Keychains, w)
	return nil
}

func (d *database) StoreKOKeychain(k account.Keychain) error {
	d.KOKeychains = append(d.KOKeychains, k)
	return nil
}

func (d *database) StoreBiometric(b account.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	return nil
}

func (d *database) StoreKOBiometric(b account.Biometric) error {
	d.KOBiometrics = append(d.KOBiometrics, b)
	return nil
}

func (d *database) NewLock(txLock lock.TransactionLock) error {
	d.Locks = append(d.Locks, txLock)
	return nil
}

func (d *database) RemoveLock(txLock lock.TransactionLock) error {
	pos := d.findLockPosition(txLock)
	if pos > -1 {
		d.Locks = append(d.Locks[:pos], d.Locks[pos+1:]...)
	}
	return nil
}

func (d database) ContainsLock(txLock lock.TransactionLock) bool {
	for _, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return true
		}
	}
	return false
}

func (d database) findLockPosition(txLock lock.TransactionLock) int {
	for i, lock := range d.Locks {
		if lock.TxHash == txLock.TxHash && txLock.MasterRobotKey == lock.MasterRobotKey {
			return i
		}
	}
	return -1
}
