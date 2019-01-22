package adding

import (
	"errors"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/inspecting"
	"github.com/uniris/uniris-core/pkg/listing"
)

//Repository define methods to handle storage
type Repository interface {

	//StoreSharedEmitterKeyPair stores a shared emitter keypair
	StoreSharedEmitterKeyPair(kp uniris.SharedKeys) error

	StoreKeychain(kc uniris.Keychain) error
	StoreID(id uniris.ID) error
	StoreKO(tx uniris.Transaction) error

	//StoreLock stores a lock
	StoreLock(l uniris.Lock) error

	//RemoveLock remove an existing lock
	RemoveLock(l uniris.Lock) error
}

//Service handle data storing
type Service struct {
	repo   Repository
	lister listing.Service
}

//NewService creates a new data storage service
func NewService(r Repository, l listing.Service) Service {
	return Service{
		repo:   r,
		lister: l,
	}
}

//StoreLock performs a lock on a transaction
//
//If a lock exist, ErrLockExisting error is returned
func (s Service) StoreLock(l uniris.Lock) error {
	exist, err := s.lister.ContainsTransactionLock(l)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("A lock already exist for this transaction")
	}
	s.repo.StoreLock(l)
	return nil
}

//RemoveLock performs a unlock on a transaction
func (s Service) RemoveLock(l uniris.Lock) error {
	exist, err := s.lister.ContainsTransactionLock(l)
	if err != nil {
		return err
	}
	if exist {
		return s.repo.RemoveLock(l)
	}
	return nil
}

//StoreSharedEmitterKeyPair handles emitter shared key storage
func (s Service) StoreSharedEmitterKeyPair(kp uniris.SharedKeys) error {
	return s.repo.StoreSharedEmitterKeyPair(kp)
}

//StoreTransaction handles the transaction storage
//
//It ensures the miner has the authorized to store the transaction
//It checks the transaction validations (master and confirmations)
//It's building the transaction chain and verify its integrity
//Then finally store in the right database
func (s Service) StoreTransaction(tx uniris.Transaction) error {
	if err := s.checkTransactionBeforeStorage(tx); err != nil {
		return err
	}

	//Check integrity of the keychain
	chainedTx, err := s.getChainedTransaction(tx)
	if err != nil {
		return err
	}
	if err := chainedTx.CheckChainTransactionIntegrity(); err != nil {
		return err
	}

	return s.storeTransaction(tx)
}

func (s Service) checkTransactionBeforeStorage(tx uniris.Transaction) error {
	if !inspecting.IsAuthorizedToStoreTx(tx.TransactionHash()) {
		return errors.New("Not authorized storage")
	}

	minValid := inspecting.GetMinimumTransactionValidation(tx.TransactionHash())
	if len(tx.ConfirmationsValidations()) < minValid {
		return errors.New("Invalid number of validations")
	}

	if err := tx.CheckProofOfWork(); err != nil {
		return err
	}

	if err := tx.MasterValidation().Validation().CheckValidation(); err != nil {
		return err
	}

	for _, v := range tx.ConfirmationsValidations() {
		if err := v.CheckValidation(); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) getChainedTransaction(tx uniris.Transaction) (chainedTx uniris.Transaction, err error) {
	prev, err := s.lister.GetPreviousTransaction(tx.Address(), tx.Type())
	if err != nil {
		return
	}
	if prev == nil {
		return tx, nil
	}

	prevTx, err := s.getChainedTransaction(*prev)
	if err != nil {
		return chainedTx, err
	}

	return uniris.NewChainedTransaction(tx, prevTx), nil
}

func (s Service) storeTransaction(tx uniris.Transaction) error {
	if tx.IsKO() {
		return s.repo.StoreKO(tx)
	}

	switch tx.Type() {
	case uniris.KeychainTransactionType:
		{
			kc, err := uniris.NewKeychain(tx)
			if err != nil {
				return err
			}
			return s.repo.StoreKeychain(kc)
		}
	case uniris.IDTransactionType:
		{
			id, err := uniris.NewID(tx)
			if err != nil {
				return err
			}
			return s.repo.StoreID(id)
		}
	}

	return nil
}
