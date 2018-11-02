package mock

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

//NewPoolRequester create a mock pool requester
func NewPoolRequester() datamining.PoolRequester {
	return mockPoolRequester{
		Repo: mock.NewDatabase(),
	}
}

type mockPoolRequester struct {
	Repo *mock.Databasemock
}

func (r mockPoolRequester) RequestLock(datamining.Pool, datamining.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(datamining.Pool, datamining.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(sPool datamining.Pool, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool datamining.Pool, data interface{}, txType datamining.TransactionType) error {
	switch data.(type) {
	case *datamining.Keychain:
		r.Repo.StoreKeychain(data.(*account.Keychain))
	case *datamining.Biometric:
		r.Repo.StoreBiometric(data.(*account.Biometric))
	}

	return nil
}
