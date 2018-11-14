package mem

import (
	"log"
	"sort"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	account_adding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	account_listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	biod_adding "github.com/uniris/uniris-core/datamining/pkg/biod/adding"
	biod_listing "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

//Repo mock the entire database
type Repo interface {
	account_adding.Repository
	account_listing.Repository
	biod_listing.Repository
	biod_adding.Repository
	lock.Repository
}

type database struct {
	Biometrics   []account.Biometric
	KOBiometrics []account.Biometric
	Keychains    []account.Keychain
	KOKeychains  []account.Keychain
	BiodPubKeys  []string
	Locks        []lock.TransactionLock
}

//NewDatabase creates a new mock database
func NewDatabase() Repo {
	return &database{}
}

func (d *database) FindBiometric(hash string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			log.Printf("Biometric with hash %s retrieved\n", hash)
			return b, nil
		}
	}
	log.Printf("Biometric with hash %s not retrieved\n", hash)
	return nil, nil
}

func (d *database) FindLastKeychain(addr string) (account.Keychain, error) {
	sort.Slice(d.Keychains, func(i, j int) bool {
		iTimestamp := d.Keychains[i].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		jTimestamp := d.Keychains[j].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, k := range d.Keychains {
		if k.Address() == addr {
			log.Printf("Keychain with address %s retrieved\n", addr)
			return k, nil
		}
	}
	log.Printf("Keychain with address %s not found\n", addr)
	return nil, nil
}

func (d *database) StoreBiodPublicKey(key string) error {
	//Prevent to add multiple times
	for _, k := range d.BiodPubKeys {
		if k == key {
			return nil
		}
	}
	d.BiodPubKeys = append(d.BiodPubKeys, key)
	log.Printf("Biometric device register with key %s\n", key)
	return nil
}

func (d *database) ListBiodPubKeys() ([]string, error) {
	return d.BiodPubKeys, nil
}

func (d *database) StoreKeychain(k account.Keychain) error {
	d.Keychains = append(d.Keychains, k)
	log.Printf("New keychain stored with address %s\n", k.Address())
	return nil
}

func (d *database) StoreKOKeychain(k account.Keychain) error {
	d.KOKeychains = append(d.KOKeychains, k)
	log.Printf("New keychain stored in KO with address %s\n", k.Address())
	return nil
}

func (d *database) StoreBiometric(b account.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	log.Printf("New biometric stored with hash %s\n", b.PersonHash())
	return nil
}

func (d *database) StoreKOBiometric(b account.Biometric) error {
	d.KOBiometrics = append(d.KOBiometrics, b)
	log.Printf("New biometric stored in KO with hash %s\n", b.PersonHash())
	return nil
}

func (d *database) NewLock(txLock lock.TransactionLock) error {
	d.Locks = append(d.Locks, txLock)
	log.Printf("New lock for the tx %s\n", txLock.TxHash)
	return nil
}

func (d *database) RemoveLock(txLock lock.TransactionLock) error {
	pos := d.findLockPosition(txLock)
	if pos > -1 {
		d.Locks = append(d.Locks[:pos], d.Locks[pos+1:]...)
	}
	log.Printf("Lock removed for the tx %s\n", txLock.TxHash)
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
