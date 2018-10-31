package mock

import (
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/locking"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"
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

func (s mockSigner) SignLock(locking.TransactionLock, string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignMasterValidation(v master.Validation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignValidation(v slave.Validation, pvKey string) (string, error) {
	return "sig", nil
}
