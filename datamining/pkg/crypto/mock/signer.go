package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//NewSigner creates a mock a signer
func NewSigner() crypto.Signer {
	return mockSigner{}
}

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

func (s mockSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSigner) SignLock(lock.TransactionLock, string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignValidation(v mining.UnsignedValidation, pvKey string) (string, error) {
	return "sig", nil
}
