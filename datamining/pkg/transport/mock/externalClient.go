package mock

import (
	"errors"
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

type mockExtClient struct {
	db mock.Database
}

//NewExternalClient create a new external mocked client
func NewExternalClient(db mock.Database) mockExtClient {
	return mockExtClient{db}
}

func (c mockExtClient) LeadKeychainMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error {
	return nil
}

func (c mockExtClient) LeadBiometricMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error {

	return nil
}

func (c mockExtClient) RequestBiometric(ip string, encPersonHash string) (account.Biometric, error) {
	return c.db.FindBiometric("hash")
}

func (c mockExtClient) RequestKeychain(ip string, encAddress string) (account.Keychain, error) {
	return c.db.FindLastKeychain("hash")
}

func (c mockExtClient) RequestLock(ip string, txLock lock.TransactionLock) error {
	return c.db.NewLock(txLock)
}

func (c mockExtClient) RequestUnlock(ip string, txLock lock.TransactionLock) error {
	return c.db.RemoveLock(txLock)
}

func (c mockExtClient) RequestValidation(ip string, txType mining.TransactionType, txHash string, data interface{}) (mining.Validation, error) {
	return mining.NewValidation(
		mining.ValidationOK,
		time.Now(),
		"pubkey",
		"fake sig",
	), nil
}

func (c mockExtClient) RequestStorage(ip string, txType mining.TransactionType, data interface{}, end mining.Endorsement) error {
	switch txType {
	case mining.KeychainTransaction:
		keychain := account.NewKeychain("address", data.(account.KeychainData), end)
		return c.db.StoreKeychain(keychain)
	case mining.BiometricTransaction:
		bio := account.NewBiometric(data.(account.BiometricData), end)
		return c.db.StoreBiometric(bio)
	}

	return errors.New("Unsupported storage")
}
