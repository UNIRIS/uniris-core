package transactions

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
)

//Requester defines methods to contact peers in the pool
type Requester interface {
	RequestStorage(sPool pool.PeerCluster, data interface{}, txType Type) error
	RequestValidations(sPool pool.PeerCluster, data interface{}, txType Type) ([]datamining.Validation, error)
}

//Handler defines methods every transaction must define
type Handler interface {
	RequestValidations(pd Requester, vPool pool.PeerCluster, data interface{}, txType Type) ([]datamining.Validation, error)
	RequestStorage(pd Requester, sPool pool.PeerCluster, data interface{}, e *datamining.Endorsement, txType Type) error
}

//Type represents the transaction type
type Type int

const (
	//CreateWallet represents a wallet creation transaction
	CreateWallet Type = 0

	//CreateBio represents a bio creation transaction
	CreateBio Type = 1
)
