package mock

import "github.com/uniris/uniris-core/datamining/pkg/crypto"

//NewHasher create a mock of hasher
func NewHasher() crypto.Hasher {
	return mockHasher{}
}

type mockHasher struct{}

func (h mockHasher) HashTransactionData(data interface{}) (string, error) {
	return crypto.NewHasher().HashTransactionData(data)
}
