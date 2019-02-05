package chain

import (
	"errors"
	"fmt"
)

//Database wrap chain database reader and writer
type Database interface {
	DatabaseReader
	DatabaseWriter
}

//DatabaseReader handles transaction chain queries to retreive a transaction from its related database
type DatabaseReader interface {

	//PendingByHash retrieves a transaction from the pending database by hash
	PendingByHash(txHash string) (*Transaction, error)

	//KOByHash retrieves a transactionfrom the KO database by hash
	KOByHash(txHash string) (*Transaction, error)

	//FullKeychain retrieves the entire transaction chain from the Keychain database by addres
	FullKeychain(addr string) (*Keychain, error)

	//LastKeychain retrieves the last keychain transaction from the Keychain database by address
	LastKeychain(addr string) (*Keychain, error)

	//KeychainByHash a transaction from the Keychain database by hash
	KeychainByHash(txHash string) (*Keychain, error)

	//ID retrieves a transaction from the ID database by address
	ID(addr string) (*ID, error)

	//ReadIDByHash retrieves a transaction from the ID database by hash
	IDByHash(txHash string) (*ID, error)
}

//DatabaseWriter handles transaction persistence by writing in the right database the related transaction
type DatabaseWriter interface {
	//WritePending stores the transaction in the pending storage
	WritePending(tx Transaction) error

	//WriteKO stores the transaction in the KO storage
	WriteKO(tx Transaction) error

	//WriteKeychain stores the transaction in the keychain storage
	WriteKeychain(kc Keychain) error

	//WriteID stores the transaction in the ID storage
	WriteID(id ID) error
}

//ErrUnknownTransaction is returned when the transaction is not found in all the ledgers (pending, ko, keychain, id, contracts)
var ErrUnknownTransaction = errors.New("unknown transaction")

//WriteTransaction stores the transaction
//
//It ensures the miner has the authorized to store the transaction
//It checks the transaction validations (master and confirmations)
//It's building the transaction chain and verify its integrity
//Then finally store in the right database
func WriteTransaction(db Database, tx Transaction, minValids int) error {
	if err := checkTransactionBeforeStorage(tx, minValids); err != nil {
		return err
	}

	if tx.IsKO() {
		return db.WriteKO(tx)
	}

	chain, err := getFullChain(db, tx.addr, tx.txType)
	if err != nil {
		return err
	}
	if err := tx.Chain(chain); err != nil {
		return err
	}

	switch tx.txType {
	case KeychainTransactionType:
		keychain, err := NewKeychain(tx)
		if err != nil {
			return err
		}
		return db.WriteKeychain(keychain)
	case IDTransactionType:
		id, err := NewID(tx)
		if err != nil {
			return err
		}
		return db.WriteID(id)
	}

	return nil
}

func checkTransactionBeforeStorage(tx Transaction, minValids int) error {
	if _, err := tx.IsValid(); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	if len(tx.ConfirmationsValidations()) < minValids {
		return errors.New("transaction: invalid number of validations")
	}

	if err := tx.CheckMasterValidation(); err != nil {
		return err
	}

	for _, v := range tx.ConfirmationsValidations() {
		if _, err := v.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

func getFullChain(db DatabaseReader, txAddr string, txType TransactionType) (*Transaction, error) {
	switch txType {
	case KeychainTransactionType:
		keychain, err := db.FullKeychain(txAddr)
		if err != nil {
			return nil, err
		}
		if keychain == nil {
			return nil, nil
		}
		return &keychain.Transaction, nil
	}

	return nil, nil
}

//LastTransaction retrieves the last transaction from the database
func LastTransaction(db DatabaseReader, txAddr string, txType TransactionType) (*Transaction, error) {
	switch txType {
	case KeychainTransactionType:
		keychain, err := db.LastKeychain(txAddr)
		if err != nil {
			return nil, err
		}
		if keychain == nil {
			return nil, nil
		}
		return &keychain.Transaction, nil
	case IDTransactionType:
		id, err := db.ID(txAddr)
		if err != nil {
			return nil, err
		}
		if id == nil {
			return nil, nil
		}
		return &id.Transaction, nil
	}

	return nil, nil
}

//GetTransactionStatus gets the status of a transaction
//It lookups on Pending DB, KO DB, Keychain, ID, Smart contracts
func GetTransactionStatus(db DatabaseReader, txHash string) (TransactionStatus, error) {
	tx, err := db.PendingByHash(txHash)
	if err != nil {
		return TransactionStatusSuccess, err
	}
	if tx != nil {
		return TransactionStatusPending, nil
	}

	tx, err = db.KOByHash(txHash)
	if err != nil {
		return TransactionStatusUnknown, err
	}
	if tx != nil {
		return TransactionStatusFailure, nil
	}

	tx, err = getTransactionByHash(db, txHash)
	if err != nil {
		if err == ErrUnknownTransaction {
			return TransactionStatusUnknown, nil
		}
		return TransactionStatusUnknown, err
	}

	return TransactionStatusSuccess, nil
}

func getTransactionByHash(db DatabaseReader, txHash string) (*Transaction, error) {
	keychain, err := db.KeychainByHash(txHash)
	if err != nil {
		return nil, err
	}
	if keychain != nil {
		return &keychain.Transaction, nil
	}
	id, err := db.IDByHash(txHash)
	if err != nil {
		return nil, err
	}
	if id != nil {
		return &id.Transaction, nil
	}

	//TODO: smart contract

	return nil, ErrUnknownTransaction
}
