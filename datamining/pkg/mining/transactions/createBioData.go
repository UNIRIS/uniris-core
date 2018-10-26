package transactions

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
)

//NewCreateBioHandler creates an transaction handler for bio creation
func NewCreateBioHandler() Handler {
	return createBioData{}
}

type createBioData struct{}

func (cw createBioData) RequestValidations(poolD Requester, vPool pool.PeerCluster, data interface{}, txType Type) ([]datamining.Validation, error) {
	valids, err := poolD.RequestValidations(vPool, data.(*datamining.BioData), txType)
	if err != nil {
		return nil, err
	}

	return valids, nil
}

func (cw createBioData) RequestStorage(poolD Requester, sPool pool.PeerCluster, data interface{}, e *datamining.Endorsement, txType Type) error {
	w := datamining.NewBioWallet(data.(*datamining.BioData), e)
	if err := poolD.RequestStorage(sPool, w, txType); err != nil {
		return err
	}
	return nil
}
