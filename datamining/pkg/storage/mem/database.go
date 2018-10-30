package mem

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

//Database represents a database
type Database interface {
	adding.AccountRepository
	listing.AccountRepository
	listing.TechRepository
}

type db struct {
	biometrics   []*datamining.Biometric
	keychains    []*datamining.Keychain
	sharedBioKey string
}

//NewDatabase implements the database in memory
func NewDatabase(sharedBioPubKey string) Database {
	return &db{
		sharedBioKey: sharedBioPubKey,
	}
}

//FindBiometric return the biometric from the given person hash
func (d *db) FindBiometric(personHash string) (*datamining.Biometric, error) {
	for _, b := range d.biometrics {
		if b.PersonHash() == personHash {
			return b, nil
		}
	}
	return nil, nil
}

//FindKeychain return the keychain from the given address
func (d *db) FindKeychain(addr string) (*datamining.Keychain, error) {
	for _, b := range d.keychains {
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}

//StoreKeychain add a keychain line to the database
func (d *db) StoreKeychain(kc *datamining.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

//StoreBiometric add a biometric line to the database
func (d *db) StoreBiometric(b *datamining.Biometric) error {
	d.biometrics = append(d.biometrics, b)
	return nil
}

func (d *db) ListBiodPubKeys() ([]string, error) {
	return []string{d.sharedBioKey}, nil
}
