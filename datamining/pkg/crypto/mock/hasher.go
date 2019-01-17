package mock

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/contract"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

type mockHasher struct{}

//NewHasher creates a new mocked hasher
func NewHasher() mockHasher {
	return mockHasher{}
}

func (h mockHasher) HashKeychain(account.Keychain) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashID(account.ID) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashEndorsedKeychain(account.EndorsedKeychain) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashEndorsedID(account.EndorsedID) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashLock(lock.TransactionLock) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashAPIContract(*api.Contract) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashEndorsedContract(contract.EndorsedContract) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashContract(contract.Contract) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashPublicKey(key string) string {
	return "hash"
}
