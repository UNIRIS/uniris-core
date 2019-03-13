package memstorage

import (
	"bytes"
	"sort"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

type chainDB struct {
	keychains []chain.Keychain
	ids       []chain.ID
	koTxs     []chain.Transaction
}

//NewchainDatabase creates a new chain database in memory
func NewchainDatabase() chain.Database {
	return &chainDB{}
}

func (db chainDB) KeychainByHash(txHash crypto.VersionnedHash) (*chain.Keychain, error) {
	for _, tx := range db.keychains {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db chainDB) FullKeychain(addr crypto.VersionnedHash) (*chain.Keychain, error) {
	sort.Slice(db.keychains, func(i, j int) bool {
		return db.keychains[i].Timestamp().Unix() > db.keychains[j].Timestamp().Unix()
	})

	if len(db.keychains) > 0 {
		return &db.keychains[0], nil
	}
	return nil, nil
}

func (db chainDB) LastKeychain(addr crypto.VersionnedHash) (*chain.Keychain, error) {

	sort.Slice(db.keychains, func(i, j int) bool {
		iTimestamp := db.keychains[i].Timestamp().Unix()
		jTimestamp := db.keychains[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range db.keychains {
		if bytes.Equal(tx.Address(), addr) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db chainDB) IDByHash(txHash crypto.VersionnedHash) (*chain.ID, error) {
	for _, tx := range db.ids {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db chainDB) ID(addr crypto.VersionnedHash) (*chain.ID, error) {

	sort.Slice(db.ids, func(i, j int) bool {
		iTimestamp := db.ids[i].Timestamp().Unix()
		jTimestamp := db.ids[j].Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, tx := range db.ids {
		if bytes.Equal(tx.Address(), addr) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db chainDB) KOByHash(txHash crypto.VersionnedHash) (*chain.Transaction, error) {
	for _, tx := range db.koTxs {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db *chainDB) WriteKeychain(kc chain.Keychain) error {
	db.keychains = append(db.keychains, kc)
	return nil
}

func (db *chainDB) WriteID(id chain.ID) error {
	db.ids = append(db.ids, id)
	return nil
}

func (db *chainDB) WriteKO(tx chain.Transaction) error {
	db.koTxs = append(db.koTxs, tx)
	return nil
}
