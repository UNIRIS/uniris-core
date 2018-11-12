package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

type mockDecrypter struct{}

//NewDecrypter create new mocked decrypter
func NewDecrypter() mockDecrypter {
	return mockDecrypter{}
}

func (d mockDecrypter) DecryptHash(hash string, pvKey string) (string, error) {
	return "hash", nil
}

func (d mockDecrypter) DecryptKeychainData(data string, pvKey string) (account.KeychainData, error) {
	return account.NewKeychainData("", "", "", "", account.NewSignatures("", "")), nil
}

func (d mockDecrypter) DecryptBiometricData(data string, pvKey string) (account.BiometricData, error) {
	return account.NewBiometricData("personHash", "", "", "", "", "", account.NewSignatures("", "")), nil
}
