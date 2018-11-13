package rpc

import (
	"errors"
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type mockPoolRequester struct {
	repo *mockDatabase
}

func (r mockPoolRequester) RequestBiometric(sPool datamining.Pool, personHash string) (account.Biometric, error) {
	return r.repo.FindBiometric("hash")
}
func (r mockPoolRequester) RequestKeychain(sPool datamining.Pool, addr string) (account.Keychain, error) {
	return r.repo.FindLastKeychain("hash")
}

func (r mockPoolRequester) RequestLock(lockPool datamining.Pool, txLock lock.TransactionLock) error {
	return r.repo.NewLock(txLock)
}

func (r mockPoolRequester) RequestUnlock(lockPool datamining.Pool, txLock lock.TransactionLock) error {
	return r.repo.RemoveLock(txLock)
}

func (r mockPoolRequester) RequestValidations(sPool datamining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]mining.Validation, error) {
	return []mining.Validation{
		mining.NewValidation(
			mining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolRequester) RequestStorage(sPool datamining.Pool, data interface{}, end mining.Endorsement, txType mining.TransactionType) error {
	switch txType {
	case mining.KeychainTransaction:
		keychain := account.NewKeychain("address", data.(account.KeychainData), end)
		return r.repo.StoreKeychain(keychain)
	case mining.BiometricTransaction:
		bio := account.NewBiometric(data.(account.BiometricData), end)
		return r.repo.StoreBiometric(bio)
	}

	return errors.New("Unsupported storage")
}
