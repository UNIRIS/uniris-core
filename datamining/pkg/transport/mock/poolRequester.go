package mock

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/locking"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
	"github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

//NewPoolRequester create a mock pool requester
func NewPoolRequester() pool.Requester {
	return mockPoolRequester{
		Repo: mock.NewDatabase(),
	}
}

type mockPoolRequester struct {
	Repo *mock.Databasemock
}

func (r mockPoolRequester) RequestLock(pool.PeerGroup, locking.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(pool.PeerGroup, locking.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(sPool pool.PeerGroup, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool pool.PeerGroup, data interface{}, txType datamining.TransactionType) error {
	switch data.(type) {
	case *datamining.Keychain:
		r.Repo.StoreKeychain(data.(*datamining.Keychain))
	case *datamining.Biometric:
		r.Repo.StoreBiometric(data.(*datamining.Biometric))
	}

	return nil
}
