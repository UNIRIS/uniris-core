package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg"
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

func (d mockDecrypter) DecryptKeychain(data string, pvKey string) (account.Keychain, error) {
	return account.NewKeychain("", "", "", datamining.NewProposal(
		datamining.NewProposedKeyPair("enc pv key", "pub key"),
	), "id sig", "em sig"), nil
}

func (d mockDecrypter) DecryptID(data string, pvKey string) (account.ID, error) {
	return account.NewID("hash", "", "", "", "", datamining.NewProposal(
		datamining.NewProposedKeyPair("enc pv key", "pub key"),
	), "id sig", "em sig"), nil
}
