package mock

import (
	"time"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/leading"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

//PoolDispatcher defines methods to handle pool dispatching
type PoolDispatcher interface {
}

//NewPoolDispatcher creates a new pool dispatcher
func NewPoolDispatcher(add adding.Service, val validating.Service) leading.PoolDispatcher {
	return poolDispatcher{
		add: add,
		val: val,
	}
}

type poolDispatcher struct {
	add adding.Service
	val validating.Service
}

func (pd poolDispatcher) RequestLastTx(pool leading.Pool, txHash string) (oldTxHash string, validation *datamining.MasterValidation, err error) {
	return "", datamining.NewMasterValidation(
		[]string{
			"key1",
			"key2",
		},
		"key",
		datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"),
	), nil
}
func (pd poolDispatcher) RequestWalletStorage(p leading.Pool, w *datamining.Wallet) error {
	return pd.add.StoreDataWallet(w)
}
func (pd poolDispatcher) RequestBioStorage(p leading.Pool, b *datamining.BioWallet) error {
	return pd.add.StoreBioWallet(b)
}

func (pd poolDispatcher) RequestLock(pool leading.Pool, txHash string) error {
	return nil
}
func (pd poolDispatcher) RequestUnlock(pool leading.Pool, txHash string) error {
	return nil
}
func (pd poolDispatcher) RequestWalletValidation(p leading.Pool, w *datamining.WalletData) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)
	for range p.Peers {
		v, err := pd.val.ValidateWalletData(w)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}
	return valids, nil
}

func (pd poolDispatcher) RequestBioValidation(p leading.Pool, b *datamining.BioData) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)
	for range p.Peers {
		v, err := pd.val.ValidateBioData(b)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}
	return valids, nil
}
