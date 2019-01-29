package transaction

import (
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//Repository handle transaction storage
type Repository interface {
	FindPendingTransaction(txHash string) (*Transaction, error)
	FindKOTransaction(txHash string) (*Transaction, error)
	StoreKO(tx Transaction) error

	KeychainRepository
	IDRepository
}

//StorageService handles transaction storage
type StorageService struct {
	repo    Repository
	mineSrv MiningService
}

//NewStorageService creates a new transaction storage service
func NewStorageService(r Repository, mS MiningService) StorageService {
	return StorageService{
		repo:    r,
		mineSrv: mS,
	}
}

//StoreTransaction handles the transaction storage
//
//It ensures the miner has the authorized to store the transaction
//It checks the transaction validations (master and confirmations)
//It's building the transaction chain and verify its integrity
//Then finally store in the right database
func (s StorageService) StoreTransaction(tx Transaction) error {
	if err := s.checkTransactionBeforeStorage(tx); err != nil {
		return err
	}

	if tx.IsKO() {
		return s.repo.StoreKO(tx)
	}

	chain, err := s.GetTransactionChain(tx.Address(), tx.Type())
	if err != nil {
		return err
	}
	if err := tx.Chain(chain); err != nil {
		return err
	}

	switch tx.Type() {
	case KeychainType:
		keychain, err := NewKeychain(tx)
		if err != nil {
			return err
		}
		return s.repo.StoreKeychain(keychain)
	case IDType:
		id, err := NewID(tx)
		if err != nil {
			return err
		}
		return s.repo.StoreID(id)
	}

	return nil
}

func (s StorageService) checkTransactionBeforeStorage(tx Transaction) error {
	if !s.isAuthorizedToStoreTx(tx.txHash) {
		return errors.New("transaction: not authorized storage")
	}

	if _, err := tx.IsValid(); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	minValid := s.mineSrv.GetMinimumTransactionValidation(tx.txHash)
	if len(tx.ConfirmationsValidations()) < minValid {
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

func (s StorageService) isAuthorizedToStoreTx(txHash string) bool {
	return true
}

//GetTransactionChain returns the entire chain for a type of transaction
func (s StorageService) GetTransactionChain(txAddr string, txType Type) (*Transaction, error) {
	switch txType {
	case KeychainType:
		keychain, err := s.repo.GetKeychain(txAddr)
		if err != nil {
			return nil, err
		}
		if keychain == nil {
			return nil, nil
		}
		tx, err := keychain.ToTransaction()
		return &tx, err
	}

	return nil, nil
}

//GetLastTransaction returns the last transaction for specific type
func (s StorageService) GetLastTransaction(txAddr string, txType Type) (*Transaction, error) {
	if _, err := crypto.IsHash(txAddr); err != nil {
		return nil, fmt.Errorf("get last transaction: %s", err.Error())
	}
	switch txType {
	case KeychainType:
		keychain, err := s.repo.FindLastKeychain(txAddr)
		if err != nil {
			return nil, err
		}
		tx, err := keychain.ToTransaction()
		if err != nil {
			return nil, err
		}
		return &tx, nil
	case IDType:
		id, err := s.repo.FindIDByAddress(txAddr)
		if err != nil {
			return nil, err
		}
		tx, err := id.ToTransaction()
		if err != nil {
			return nil, err
		}
		return &tx, nil
	}

	return nil, nil
}

//GetTransactionStatus gets the status of a transaction
//
//It lookups on Pending DB, KO DB, Keychain, ID, Smart contracts
func (s StorageService) GetTransactionStatus(txHash string) (Status, error) {
	if _, err := crypto.IsHash(txHash); err != nil {
		return StatusUnknown, fmt.Errorf("get transaction status: %s", err.Error())
	}

	tx, err := s.repo.FindPendingTransaction(txHash)
	if err != nil {
		return StatusSuccess, err
	}
	if tx != nil {
		return StatusPending, nil
	}

	tx, err = s.repo.FindKOTransaction(txHash)
	if err != nil {
		return StatusUnknown, err
	}
	if tx != nil {
		return StatusFailure, nil
	}

	tx, err = s.getTransactionByHash(txHash)
	if err != nil {
		if err.Error() == "unknown transaction" {
			return StatusUnknown, nil
		}
		return StatusUnknown, err
	}

	return StatusSuccess, nil
}

func (s StorageService) getTransactionByHash(txHash string) (*Transaction, error) {

	keychain, err := s.getTransactionByHashAndType(txHash, KeychainType)
	if err != nil {
		return nil, err
	}
	if keychain != nil {
		return keychain, nil
	}

	id, err := s.getTransactionByHashAndType(txHash, IDType)
	if err != nil {
		return nil, err
	}
	if id != nil {
		return id, nil
	}

	//TODO: smart contract

	return nil, errors.New("unknown transaction")
}

func (s StorageService) getTransactionByHashAndType(txHash string, txType Type) (*Transaction, error) {
	if _, err := crypto.IsHash(txHash); err != nil {
		return nil, fmt.Errorf("get transaction hash by type: %s", err.Error())
	}

	switch txType {
	case KeychainType:
		keychain, err := s.repo.FindKeychainByHash(txHash)
		if err != nil {
			return nil, err
		}
		if keychain == nil {
			break
		}
		tx, err := keychain.ToTransaction()
		return &tx, err
	case IDType:
		id, err := s.repo.FindIDByHash(txHash)
		if err != nil {
			return nil, err
		}
		if id == nil {
			break
		}
		tx, err := id.ToTransaction()
		return &tx, err
	}

	return nil, nil
}
