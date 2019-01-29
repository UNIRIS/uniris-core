package memstorage

import (
	"github.com/uniris/uniris-core/pkg/shared"
)

//SharedDatabase is the shared memory database
type SharedDatabase interface {
	shared.KeyRepository
}

type sharedDb struct {
	emitterKeys []shared.KeyPair
}

//NewSharedDatabase creates a new memory shared database
func NewSharedDatabase() SharedDatabase {
	return &sharedDb{}
}

func (d sharedDb) ListSharedEmitterKeyPairs() ([]shared.KeyPair, error) {
	return d.emitterKeys, nil
}

func (d *sharedDb) StoreSharedEmitterKeyPair(sk shared.KeyPair) error {
	d.emitterKeys = append(d.emitterKeys, sk)
	return nil
}
