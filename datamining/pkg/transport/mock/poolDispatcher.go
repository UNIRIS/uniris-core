package mock

import (
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

func (pd poolDispatcher) RequestLastTx(pool leading.Pool, txHash string) (oldTxHash string, err error) {
	return "", nil
}
func (pd poolDispatcher) RequestWalletStorage(p leading.Pool, w *datamining.Wallet) error {
	return pd.add.StoreDataWallet(w)
}
func (pd poolDispatcher) RequestBioStorage(p leading.Pool, b *datamining.BioWallet) error {
	return pd.add.StoreBioWallet(b)
}

func (pd poolDispatcher) RequestLock(pool leading.Pool, txLock validating.TransactionLock, sig string) error {
	return pd.val.LockTransaction(txLock, sig)
}
func (pd poolDispatcher) RequestUnlock(pool leading.Pool, txLock validating.TransactionLock, sig string) error {
	return pd.val.UnlockTransaction(txLock, sig)
}
func (pd poolDispatcher) RequestWalletValidation(p leading.Pool, w *datamining.WalletData, txHash string) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)
	for range p.Peers {
		v, err := pd.val.ValidateWalletData(w, txHash)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}
	return valids, nil
}

func (pd poolDispatcher) RequestBioValidation(p leading.Pool, b *datamining.BioData, txHash string) ([]datamining.Validation, error) {
	valids := make([]datamining.Validation, 0)
	for range p.Peers {
		v, err := pd.val.ValidateBioData(b, txHash)
		if err != nil {
			return nil, err
		}
		valids = append(valids, v)
	}
	return valids, nil
}
