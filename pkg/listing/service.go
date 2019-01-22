package listing

import (
	"errors"

	uniris "github.com/uniris/uniris-core/pkg"
)

//Repository defines methods to handle listing storage
type Repository interface {

	//ListSharedEmitterKeyPairs gets the shared emitter keypair
	ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error)

	FindPendingTransaction(txHash string) (*uniris.Transaction, error)
	FindKOTransaction(txHash string) (*uniris.Transaction, error)

	FindKeychainByHash(txHash string) (*uniris.Keychain, error)
	FindKeychainByAddress(addr string) (*uniris.Keychain, error)

	FindIDByHash(txHash string) (*uniris.ID, error)
	FindIDByAddress(addr string) (*uniris.ID, error)

	//ContainsLocks determines if a lock exists or not
	ContainsLock(uniris.Lock) (bool, error)
}

//Service handles data retreiving
type Service struct {
	repo Repository
}

//NewService creates a new service to retrieve data
func NewService(repo Repository) Service {
	return Service{repo}
}

//ContainsTransactionLock determines if a transaction lock exists
func (s Service) ContainsTransactionLock(l uniris.Lock) (bool, error) {
	ok, err := s.repo.ContainsLock(l)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//GetID gets an ID from its address
func (s Service) GetID(addr string) (*uniris.ID, error) {
	return s.repo.FindIDByAddress(addr)
}

//GetKeychain gets a keychain based on its address
func (s Service) GetKeychain(addr string) (*uniris.Keychain, error) {
	return s.repo.FindKeychainByAddress(addr)
}

//ListSharedEmitterKeyPairs get the shared emitter key pairs
func (s Service) ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error) {
	return s.repo.ListSharedEmitterKeyPairs()
}

//IsEmitterAuthorized checks if the emitter public key is authorized
func (s Service) IsEmitterAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract
	return true, nil
}

//GetPreviousTransaction retrieve the previous keychain from a given account's address
func (s Service) GetPreviousTransaction(addr string, txType uniris.TransactionType) (tx *uniris.Transaction, err error) {
	return nil, nil
}

//GetTransactionStatus gets the status of a transaction
//
//It lookups on Pending DB, KO DB, Keychain, ID, Smart contracts
func (s Service) GetTransactionStatus(txHash string) (uniris.TransactionStatus, error) {
	tx, err := s.repo.FindPendingTransaction(txHash)
	if err != nil {
		return uniris.UnknownTransaction, err
	}
	if tx != nil {
		return uniris.PendingTransaction, nil
	}

	tx, err = s.repo.FindKOTransaction(txHash)
	if err != nil {
		return uniris.UnknownTransaction, err
	}
	if tx != nil {
		return uniris.FailureTransaction, nil
	}

	tx, err = s.getTransactionByHash(txHash)
	if err != nil {
		if err.Error() == "Unknown transaction" {
			return uniris.UnknownTransaction, nil
		}
		return uniris.UnknownTransaction, err
	}

	if tx.IsKO() {
		return uniris.FailureTransaction, nil
	}

	return uniris.SuccessTransaction, nil
}

func (s Service) getTransactionByHash(txHash string) (*uniris.Transaction, error) {
	keychainTx, err := s.repo.FindKeychainByHash(txHash)
	if err != nil {
		return nil, err
	}
	if keychainTx != nil {
		tx, err := keychainTx.ToTransaction()
		if err != nil {
			return nil, err
		}
		return &tx, nil
	}

	idTx, err := s.repo.FindIDByHash(txHash)
	if err != nil {
		return nil, err
	}
	if idTx != nil {
		tx, err := idTx.ToTransaction()
		if err != nil {
			return nil, err
		}
		return &tx, nil
	}

	//TODO: smart contract

	return nil, errors.New("Unknown transaction")
}
