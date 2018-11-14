package mock

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type mockPoolRequester struct {
	cli mockExtClient
}

func NewPoolRequester(cli mockExtClient) mockPoolRequester {
	return mockPoolRequester{cli}
}

func (r mockPoolRequester) RequestBiometric(sPool datamining.Pool, personHash string) (account.Biometric, error) {
	return r.cli.RequestBiometric("127.0.0.1", personHash)
}
func (r mockPoolRequester) RequestKeychain(sPool datamining.Pool, addr string) (account.Keychain, error) {
	return r.cli.RequestKeychain("127.0.0.1", addr)
}

func (r mockPoolRequester) RequestLock(lockPool datamining.Pool, txLock lock.TransactionLock) error {
	return r.cli.RequestLock("127.0.0.1", txLock)
}

func (r mockPoolRequester) RequestUnlock(lockPool datamining.Pool, txLock lock.TransactionLock) error {
	return r.cli.RequestUnlock("127.0.0.1", txLock)
}

func (r mockPoolRequester) RequestValidations(minValid int, sPool datamining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]mining.Validation, error) {
	v, _ := r.cli.RequestValidation("127.0.0.1", txType, txHash, data)
	return []mining.Validation{v}, nil
}

func (r mockPoolRequester) RequestStorage(sPool datamining.Pool, data interface{}, end mining.Endorsement, txType mining.TransactionType) error {
	return r.cli.RequestStorage("127.0.0.1", txType, data, end)
}
