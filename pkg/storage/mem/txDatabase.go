package memstorage

import (
	"sort"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
)

//TransactionDatabase is the transaction memory database
type TransactionDatabase interface {
	adding.TransactionRepository
	listing.TransactionRepository
}

type txDb struct {
	keychains  []uniris.Keychain
	ids        []uniris.ID
	koTxs      []uniris.Transaction
	pendingTxs []uniris.Transaction
}

//NewTransactionDatabase creates a new memory transaction database
func NewTransactionDatabase() TransactionDatabase {
	return &txDb{}
}

func (d txDb) FindPendingTransaction(txHash string) (*uniris.Transaction, error) {
	for _, tx := range d.pendingTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindKeychainByHash(txHash string) (*uniris.Keychain, error) {
	for _, tx := range d.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindKeychainByAddress(addr string) (*uniris.Keychain, error) {

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

func (d txDb) FindIDByHash(txHash string) (*uniris.ID, error) {
	for _, tx := range d.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d txDb) FindIDByAddress(addr string) (*uniris.ID, error) {

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

func (d txDb) FindKOTransaction(txHash string) (*uniris.Transaction, error) {
	for _, tx := range d.koTxs {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (d *txDb) StoreKeychain(kc uniris.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *txDb) StoreID(id uniris.ID) error {
	d.ids = append(d.ids, id)
	return nil
}

func (d *txDb) StoreKO(tx uniris.Transaction) error {
	d.koTxs = append(d.koTxs, tx)
	return nil
}
