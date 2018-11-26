package mock

import (
	"errors"
	"time"

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

func (c mockExtClient) LeadKeychainMining(ip string, txHash string, encData string, validators []string) error {
	return nil
}

func (c mockExtClient) LeadIDMining(ip string, txHash string, encData string, validators []string) error {
	return nil
}

func (c mockExtClient) RequestID(ip string, encPersonHash string) (account.EndorsedID, error) {
	return c.db.FindID("hash")
}

func (c mockExtClient) RequestKeychain(ip string, encAddress string) (account.EndorsedKeychain, error) {
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
		keychain := account.NewEndorsedKeychain("address", data.(account.Keychain), end)
		return c.db.StoreKeychain(keychain)
	case mining.IDTransaction:
		bio := account.NewEndorsedID(data.(account.ID), end)
		return c.db.StoreID(bio)
	}

	return errors.New("Unsupported storage")
}
