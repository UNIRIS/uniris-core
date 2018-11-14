package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
)

type mockHasher struct{}

//NewHasher creates a new mocked hasher
func NewHasher() mockHasher {
	return mockHasher{}
}

func (h mockHasher) HashBiodPublicKey(pubKey string) string {
	return "hash"
}

func (h mockHasher) HashKeychain(account.Keychain) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashBiometric(account.Biometric) (string, error) {
	return "hash", nil
}

func (h mockHasher) NewKeychainDataHash(account.KeychainData) (string, error) {
	return "hash", nil
}
func (h mockHasher) NewBiometricDataHash(account.BiometricData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashKeychainData(account.KeychainData) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashBiometricData(account.BiometricData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashLock(lock.TransactionLock) (string, error) {
	return "hash", nil
}
