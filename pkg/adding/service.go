package adding

import (
	"errors"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/inspecting"
	"github.com/uniris/uniris-core/pkg/listing"
)

//Repository define methods to handle storage
type Repository interface {

	//StoreKeychain persists the keychain
	StoreKeychain(uniris.Keychain) error

	//StoreID persists theID
	StoreID(uniris.ID) error

	//StoreSharedEmitterKeyPair stores a shared emitter keypair
	StoreSharedEmitterKeyPair(kp uniris.SharedKeys) error
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

//StoreKeychain handles keychain transaction storage
func (s Service) StoreKeychain(kc uniris.Keychain) error {
	if err := s.checkTransactionBeforeStorage(kc.Transaction); err != nil {
		return err
	}

	//Check integrity of the keychain
	prevKc, err := s.lister.GetLastKeychain(kc.Address())
	if err != nil {
		return err
	}
	kc.Chain(prevKc)
	if err := kc.CheckChainTransactionIntegrity(s.txHasher, s.txVerif); err != nil {
		return err
	}

	if kc.Transaction.IsKO() {
		//TODO: store on the ko database
	}

	return s.repo.StoreKeychain(kc)
}

//StoreID handles ID transaction storage
func (s Service) StoreID(id uniris.ID) error {
	if err := s.checkTransactionBeforeStorage(id.Transaction); err != nil {
		return err
	}

	if err := id.CheckTransactionIntegrity(s.txHasher, s.txVerif); err != nil {
		return err
	}

	if id.Transaction.IsKO() {
		//TODO: store on the ko database
	}

	return s.repo.StoreID(id)
}

func (s Service) checkTransactionBeforeStorage(tx uniris.Transaction) error {
	if !inspecting.IsAuthorizedToStoreTx(tx.TransactionHash()) {
		return errors.New("Not authorized storage")
	}

	minValid := inspecting.GetMinimumTransactionValidation(tx.TransactionHash())
	if len(tx.Mining().Validations()) < minValid {
		return errors.New("Invalid number of validations")
	}

	if err := tx.CheckProofOfWork(s.txVerif); err != nil {
		return err
	}

	if err := tx.Mining().MasterValidation().CheckValidation(s.txVerif); err != nil {
		return err
	}

	for _, v := range tx.Mining().Validations() {
		if err := v.CheckValidation(s.txVerif); err != nil {
			return err
		}
	}

	return nil
}
