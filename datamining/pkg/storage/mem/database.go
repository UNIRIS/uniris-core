package mem

import (
	"sort"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	account_adding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	account_listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/emitter"
	em_adding "github.com/uniris/uniris-core/datamining/pkg/emitter/adding"
	em_listing "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Repo mock the entire database
type Repo interface {
	account_adding.Repository
	account_listing.Repository
	em_listing.Repository
	em_adding.Repository
	lock.Repository
}

type database struct {
	IDs         []account.EndorsedID
	KOIDs       []account.EndorsedID
	Keychains   []account.EndorsedKeychain
	KOKeychains []account.EndorsedKeychain
	Locks       []lock.TransactionLock
	SharedEmKP  []emitter.SharedKeyPair
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

	for _, id := range d.KOIDs {
		if id.Hash() == hash {
			return id, nil
		}
	}
	return nil, nil
}

func (d *database) FindIDByTransaction(txHash string) (account.EndorsedID, error) {
	for _, id := range d.IDs {
		if id.Endorsement().TransactionHash() == txHash {
			return id, nil
		}
	}

	for _, id := range d.KOIDs {
		if id.Endorsement().TransactionHash() == txHash {
			return id, nil
		}
	}
	return nil, nil
}

func (d *database) FindKeychain(addr string, txHash string) (account.EndorsedKeychain, error) {
	for _, kc := range d.Keychains {
		if kc.Address() == addr && kc.Endorsement().TransactionHash() == txHash {
			return kc, nil
		}
	}

	for _, kc := range d.KOKeychains {
		if kc.Address() == addr && kc.Endorsement().TransactionHash() == txHash {
			return kc, nil
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

func (d *database) StoreSharedEmitterKeyPair(kp emitter.SharedKeyPair) error {
	d.SharedEmKP = append(d.SharedEmKP, kp)
	return nil
}

func (d *database) ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error) {
	return d.SharedEmKP, nil
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
