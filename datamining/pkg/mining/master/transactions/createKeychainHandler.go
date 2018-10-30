package transactions

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

type createKeychainHandler struct {
}

//NewCreateKeychainHandler creates an transaction handler for keychain creation
func NewCreateKeychainHandler() Handler {
	return createKeychainHandler{}
}

func (h createKeychainHandler) RequestValidations(poolD pool.Requester, vPool pool.Cluster, data interface{}) ([]datamining.Validation, error) {
	valids, err := poolD.RequestValidations(vPool, data.(*datamining.KeyChainData), datamining.CreateKeychainTransaction)
	if err != nil {
		return nil, err
	}

	return valids, nil
}

func (h createKeychainHandler) RequestStorage(poolD pool.Requester, sPool pool.Cluster, data interface{}, e *datamining.Endorsement) error {
	w := datamining.NewKeychain(data.(*datamining.KeyChainData), e, "")
	if err := poolD.RequestStorage(sPool, w, datamining.CreateKeychainTransaction); err != nil {
		return err
	}
	return nil
}
