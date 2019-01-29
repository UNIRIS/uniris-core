package memstorage

import (
	"sort"

	"github.com/uniris/uniris-core/pkg/transaction"
)

//TransactionDatabase is the transaction memory database
type TransactionDatabase interface {
	transaction.Repository
}

type txDb struct {
	keychains  []transaction.Keychain
	ids        []transaction.ID
	koTxs      []transaction.Transaction
	pendingTxs []transaction.Transaction
}

//NewTransactionDatabase creates a new memory transaction database
func NewTransactionDatabase() TransactionDatabase {
	return &txDb{}
}

func (d txDb) FindPendingTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range d.pendingTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindKeychainByHash(txHash string) (*transaction.Keychain, error) {
	for _, tx := range d.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) GetKeychain(addr string) (*transaction.Keychain, error) {
	sort.Slice(d.keychains, func(i, j int) bool {
		return d.keychains[i].Timestamp().Unix() > d.keychains[j].Timestamp().Unix()
	})

	if len(d.keychains) > 0 {
		return &d.keychains[0], nil
	}
	return nil, nil
}

func (d txDb) FindLastKeychain(addr string) (*transaction.Keychain, error) {

	sort.Slice(d.keychains, func(i, j int) bool {
		iTimestamp := d.keychains[i].Timestamp().Unix()
		jTimestamp := d.keychains[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range d.keychains {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindIDByHash(txHash string) (*transaction.ID, error) {
	for _, tx := range d.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindIDByAddress(addr string) (*transaction.ID, error) {

	sort.Slice(d.ids, func(i, j int) bool {
		iTimestamp := d.ids[i].Timestamp().Unix()
		jTimestamp := d.ids[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range d.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindKOTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range d.koTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d *txDb) StoreKeychain(kc transaction.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *txDb) StoreID(id transaction.ID) error {
	d.ids = append(d.ids, id)
	return nil
}

func (d *txDb) StoreKO(tx transaction.Transaction) error {
	d.koTxs = append(d.koTxs, tx)
	return nil
}
