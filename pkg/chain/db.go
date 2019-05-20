package chain

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/uniris/uniris-cli/pkg/crypto"
// )

// //Database wrap chain database reader and writer
// type Database interface {
// 	DatabaseReader
// 	DatabaseWriter
// }

// //DatabaseReader handles transaction chain queries to retreive a transaction from its related database
// type DatabaseReader interface {

// 	//KOByHash retrieves a transactionfrom the KO database by hash
// 	KOByHash(txHash crypto.VersionnedHash) (*Transaction, error)

// 	//FullKeychain retrieves the entire transaction chain from the Keychain database by addres
// 	FullKeychain(addr crypto.VersionnedHash) (*Keychain, error)

// 	//LastKeychain retrieves the last keychain transaction from the Keychain database by address
// 	LastKeychain(addr crypto.VersionnedHash) (*Keychain, error)

// 	//KeychainByHash a transaction from the Keychain database by hash
// 	KeychainByHash(txHash crypto.VersionnedHash) (*Keychain, error)

// 	//ID retrieves a transaction from the ID database by address
// 	ID(addr crypto.VersionnedHash) (*ID, error)

// 	//ReadIDByHash retrieves a transaction from the ID database by hash
// 	IDByHash(txHash crypto.VersionnedHash) (*ID, error)
// }

// //DatabaseWriter handles transaction persistence by writing in the right database the related transaction
// type DatabaseWriter interface {

// 	//WriteKO stores the transaction in the KO storage
// 	WriteKO(tx Transaction) error

// 	//WriteKeychain stores the transaction in the keychain storage
// 	WriteKeychain(kc Keychain) error

// 	//WriteID stores the transaction in the ID storage
// 	WriteID(id ID) error
// }

// //ErrUnknownTransaction is returned when the transaction is not found in all the ledgers (in progress, ko, keychain, id, contracts)
// var ErrUnknownTransaction = errors.New("unknown transaction")

// //WriteTransaction stores the transaction
// //
// //It ensures the node has the authorized to store the transaction
// //It checks the transaction validations (master and confirmations)
// //It's building the transaction chain and verify its integrity
// //Then finally store in the right database
// func WriteTransaction(chainDB Database, tx Transaction, minValids int) error {
// 	if err := checkTransactionBeforeStorage(tx, minValids); err != nil {
// 		return err
// 	}

// 	if tx.IsKO() {
// 		return chainDB.WriteKO(tx)
// 	}

// 	chain, err := getFullChain(chainDB, tx.addr, tx.txType)
// 	if err != nil {
// 		return err
// 	}
// 	if err := tx.Chain(chain); err != nil {
// 		return err
// 	}

// 	switch tx.txType {
// 	case KeychainTransactionType:
// 		keychain, err := NewKeychain(tx)
// 		if err != nil {
// 			return err
// 		}
// 		if err := chainDB.WriteKeychain(keychain); err != nil {
// 			return err
// 		}
// 	case IDTransactionType:
// 		id, err := NewID(tx)
// 		if err != nil {
// 			return err
// 		}
// 		if err := chainDB.WriteID(id); err != nil {
// 			return err
// 		}
// 	}

// 	removeTimeLock(tx.TransactionHash(), tx.Address())
// 	return nil
// }

// func checkTransactionBeforeStorage(tx Transaction, minValids int) error {
// 	if _, err := tx.IsValid(); err != nil {
// 		return fmt.Errorf("transaction: %s", err.Error())
// 	}

// 	if len(tx.ConfirmationsValidations()) < minValids {
// 		return errors.New("transaction: invalid number of validations")
// 	}

// 	if err := tx.CheckMasterValidation(); err != nil {
// 		return err
// 	}

// 	for _, v := range tx.ConfirmationsValidations() {
// 		if _, err := v.IsValid(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func getFullChain(db DatabaseReader, txAddr crypto.VersionnedHash, txType TransactionType) (*Transaction, error) {
// 	switch txType {
// 	case KeychainTransactionType:
// 		keychain, err := db.FullKeychain(txAddr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if keychain == nil {
// 			return nil, nil
// 		}
// 		return &keychain.Transaction, nil
// 	}

// 	return nil, nil
// }

// //LastTransaction retrieves the last transaction from the database
// func LastTransaction(db DatabaseReader, txAddr crypto.VersionnedHash, txType TransactionType) (*Transaction, error) {
// 	switch txType {
// 	case KeychainTransactionType:
// 		keychain, err := db.LastKeychain(txAddr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if keychain == nil {
// 			return nil, nil
// 		}
// 		return &keychain.Transaction, nil
// 	case IDTransactionType:
// 		id, err := db.ID(txAddr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if id == nil {
// 			return nil, nil
// 		}
// 		return &id.Transaction, nil
// 	}

// 	return nil, nil
// }

// //GetTransactionStatus gets the status of a transaction
// //It lookups on timelockers, KO DB, Keychain, ID, Smart contracts
// func GetTransactionStatus(db DatabaseReader, txHash crypto.VersionnedHash) (TransactionStatus, error) {
// 	if transactionHashTimeLocked(txHash) {
// 		return TransactionStatusInProgress, nil
// 	}

// 	tx, err := db.KOByHash(txHash)
// 	if err != nil {
// 		return TransactionStatusUnknown, err
// 	}
// 	if tx != nil {
// 		return TransactionStatusFailure, nil
// 	}

// 	tx, err = getTransactionByHash(db, txHash)
// 	if err != nil {
// 		if err == ErrUnknownTransaction {
// 			return TransactionStatusUnknown, nil
// 		}
// 		return TransactionStatusUnknown, err
// 	}

// 	return TransactionStatusSuccess, nil
// }

// func getTransactionByHash(db DatabaseReader, txHash crypto.VersionnedHash) (*Transaction, error) {
// 	keychain, err := db.KeychainByHash(txHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if keychain != nil {
// 		return &keychain.Transaction, nil
// 	}
// 	id, err := db.IDByHash(txHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if id != nil {
// 		return &id.Transaction, nil
// 	}

// 	//TODO: smart contract

// 	return nil, ErrUnknownTransaction
// }
