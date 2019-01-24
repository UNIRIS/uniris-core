package memstorage

import (
	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
)

//SharedDatabase is the shared memory database
type SharedDatabase interface {
	adding.SharedRepository
	listing.SharedRepository
}

type sharedDb struct {
	sharedKeys []uniris.SharedKeys
}

//NewSharedDatabase creates a new memory shared database
func NewSharedDatabase() SharedDatabase {
	return &sharedDb{}
}

func (d sharedDb) ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error) {
	return d.sharedKeys, nil
}

func (d *sharedDb) StoreSharedEmitterKeyPair(sk uniris.SharedKeys) error {
	d.sharedKeys = append(d.sharedKeys, sk)
	return nil
}
