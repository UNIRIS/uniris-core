package mock

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

type repo interface {
	adding.AccountRepository
	listing.AccountRepository
	listing.TechRepository
}

type Databasemock struct {
	repo
	Biometrics []*datamining.Biometric
	Keychains  []*datamining.Keychain
}

//NewDatabase creates a new mock database
func NewDatabase() *Databasemock {
	return &Databasemock{
		Biometrics: make([]*datamining.Biometric, 0),
		Keychains:  make([]*datamining.Keychain, 0),
	}
}

func (d *Databasemock) FindBiometric(hash string) (*datamining.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			return b, nil
		}
	}
	return nil, nil
}

func (d *Databasemock) FindKeychain(addr string) (*datamining.Keychain, error) {
	for _, w := range d.Keychains {
		if w.WalletAddr() == addr {
			return w, nil
		}
	}
	return nil, nil
}

func (d *Databasemock) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

func (d *Databasemock) StoreKeychain(w *datamining.Keychain) error {
	d.Keychains = append(d.Keychains, w)
	return nil
}

func (d *Databasemock) StoreBiometric(bw *datamining.Biometric) error {
	d.Biometrics = append(d.Biometrics, bw)
	return nil
}
