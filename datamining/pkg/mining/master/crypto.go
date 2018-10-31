package master

import (
	"github.com/uniris/uniris-core/datamining/pkg/locking"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/checks"
)

//Hasher wraps transaction type hasher
type Hasher interface {
	checks.TransactionDataHasher
}

//Signer defines methods to handle signatures
type Signer interface {
	SignLock(lock locking.TransactionLock, pvKey string) (string, error)
	PowSigner
}
