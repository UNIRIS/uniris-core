package mock

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//NewPoolRequester create a mock pool requester
func NewPoolRequester(db Repo) mining.PoolRequester {
	return mockPoolRequester{db}
}

type mockPoolRequester struct {
	Repo Repo
}

func (r mockPoolRequester) RequestLock(mining.Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestUnlock(mining.Pool, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolRequester) RequestValidations(sPool mining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool mining.Pool, data interface{}, end datamining.Endorsement, txType mining.TransactionType) error {
	switch data.(type) {
	case *account.KeyChainData:
		data := data.(*account.KeyChainData)
		kc := account.NewKeychain(data, end)
		r.Repo.StoreKeychain(kc)
	case *account.BioData:
		data := data.(*account.BioData)
		bio := account.NewBiometric(data, end)
		r.Repo.StoreBiometric(bio)
	}

	return nil
}
