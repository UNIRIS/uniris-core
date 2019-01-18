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
}

//Service handle data storing
type Service struct {
	repo     Repository
	lister   listing.Service
	txVerif  uniris.TransactionVerifier
	txHasher uniris.TransactionHasher
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
	if err := chainedTx.CheckChainTransactionIntegrity(s.txHasher, s.txVerif); err != nil {
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

	if err := tx.CheckProofOfWork(s.txVerif); err != nil {
		return err
	}

	if err := tx.MasterValidation().Validation().CheckValidation(s.txVerif); err != nil {
		return err
	}

	for _, v := range tx.ConfirmationsValidations() {
		if err := v.CheckValidation(s.txVerif); err != nil {
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
