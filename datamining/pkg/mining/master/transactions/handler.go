package transactions

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

//Handler defines methods a transaction handler have to define
type Handler interface {
	RequestValidations(poolD pool.Requester, vPool pool.PeerGroup, data interface{}) ([]datamining.Validation, error)
	RequestStorage(poolD pool.Requester, sPool pool.PeerGroup, data interface{}, e *datamining.Endorsement) error
}
