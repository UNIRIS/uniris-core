package transactions

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

type createBiometricHandler struct {
}

//NewCreateBiometricHandler creates an transaction handler for biometric creation
func NewCreateBiometricHandler() Handler {
	return createBiometricHandler{}
}

func (h createBiometricHandler) RequestValidations(poolD pool.Requester, vPool pool.PeerGroup, data interface{}) ([]datamining.Validation, error) {
	valids, err := poolD.RequestValidations(vPool, data.(*datamining.BioData), datamining.CreateBioTransaction)
	if err != nil {
		return nil, err
	}

	return valids, nil
}

func (h createBiometricHandler) RequestStorage(poolD pool.Requester, sPool pool.PeerGroup, data interface{}, e *datamining.Endorsement) error {
	b := datamining.NewBiometric(data.(*datamining.BioData), e)
	if err := poolD.RequestStorage(sPool, b, datamining.CreateBioTransaction); err != nil {
		return err
	}
	return nil
}
