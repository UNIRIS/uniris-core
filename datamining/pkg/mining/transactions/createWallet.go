package transactions

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
)

//NewCreateWalletHandler creates an transaction handler for wallet creation
func NewCreateWalletHandler() Handler {
	return createWallet{}
}

type createWallet struct{}

func (cw createWallet) RequestValidations(poolD Requester, vPool pool.PeerCluster, data interface{}, txType Type) ([]datamining.Validation, error) {
	valids, err := poolD.RequestValidations(vPool, data.(*datamining.WalletData), txType)
	if err != nil {
		return nil, err
	}

	return valids, nil
}

func (cw createWallet) RequestStorage(poolD Requester, sPool pool.PeerCluster, data interface{}, e *datamining.Endorsement, txType Type) error {
	w := datamining.NewWallet(data.(*datamining.WalletData), e, "")
	if err := poolD.RequestStorage(sPool, w, txType); err != nil {
		return err
	}
	return nil
}
