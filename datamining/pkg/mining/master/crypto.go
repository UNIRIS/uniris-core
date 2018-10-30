package master

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/checks"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

//Hasher wraps transaction type hasher
type Hasher interface {
	checks.TransactionDataHasher
}

//Signer defines methods to handle signatures
type Signer interface {
	SignLock(lock pool.TransactionLock, pvKey string) (string, error)
	PowSigner
}
