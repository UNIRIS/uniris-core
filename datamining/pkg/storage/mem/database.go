package mem

import (
	"sort"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	account_adding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	account_listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	em_listing "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Repo mock the entire database
type Repo interface {
	account_adding.Repository
	account_listing.Repository
	em_listing.Repository
	lock.Repository
}

type database struct {
	IDs         []account.EndorsedID
	KOIDs       []account.EndorsedID
	Keychains   []account.EndorsedKeychain
	KOKeychains []account.EndorsedKeychain
	Locks       []lock.TransactionLock
}

//NewDatabase creates a new mock database
func NewDatabase() Repo {
	return &database{}
}

func (d *database) FindID(hash string) (account.EndorsedID, error) {
	for _, id := range d.IDs {
		if id.Hash() == hash {
			return id, nil
		}
	}
	return nil, nil
}

func (d *database) FindLastKeychain(addr string) (account.EndorsedKeychain, error) {
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

func (d *database) ListEmitterPublicKeys() ([]string, error) {
	return []string{"3059301306072a8648ce3d020106082a8648ce3d03010703420004d64bb03b6dc787c848937793988ca158bdd7e9317cd0fef0a1bc5f040343af9e629e2277ffe0cbca0197ca1a443340c7f1c98b07c4e1ea056f2d99655e69c096"}, nil
}

func (d *database) StoreKeychain(k account.EndorsedKeychain) error {
	d.Keychains = append(d.Keychains, k)
	return nil
}

func (d *database) StoreKOKeychain(k account.EndorsedKeychain) error {
	d.KOKeychains = append(d.KOKeychains, k)
	return nil
}

func (d *database) StoreID(id account.EndorsedID) error {
	d.IDs = append(d.IDs, id)
	return nil
}

func (d *database) StoreKOID(id account.EndorsedID) error {
	d.KOIDs = append(d.KOIDs, id)
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
