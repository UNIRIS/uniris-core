package adding

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//ErrInvalidDataIntegrity is returned when the data does not valid the integrity hashes
var ErrInvalidDataIntegrity = errors.New("Invalid data integrity")

//ErrInvalidDataMining is returned when the validations does not match their signatures
var ErrInvalidDataMining = errors.New("Mining is invalid")

//ErrInvalidValidationNumber is returned when the validations number is not reached
var ErrInvalidValidationNumber = errors.New("Invalid validations number")

//Repository handles account storage
type Repository interface {
	//StoreKeychain persists the keychain
	StoreKeychain(account.EndorsedKeychain) error

	//StoreID persists the ID
	StoreID(account.EndorsedID) error

	//StoreKOKeychain persists the keychain in the KO database
	StoreKOKeychain(account.EndorsedKeychain) error

	//StoreKOID persists the ID in the KO database
	StoreKOID(account.EndorsedID) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {

	//StoreKeychain processes the keychain storage
	//
	//It performs checks to insure the integrity of the keychain
	//Determines if the keychain must be store in the KO, OK or pending database
	StoreKeychain(account.EndorsedKeychain) error

	//StoreID processes the ID storage
	//
	//It performs checks to insure the integrity of the ID
	//Determines if the ID must be store in the KO, OK or pending database
	StoreID(account.EndorsedID) error
}

type signatureVerifier interface {
	account.KeychainSignatureVerifier
	account.IDSignatureVerifier
	mining.ValidationVerifier
	mining.PowSigVerifier
}

type hasher interface {
	account.KeychainHasher
	account.IDHasher
}

type service struct {
	aiClient AIClient
	repo     Repository
	lister   listing.Service
	sigVerif signatureVerifier
	hasher   hasher
}

//NewService creates a new adding service
func NewService(aiClient AIClient, repo Repository, lister listing.Service, sigVerif signatureVerifier, hash hasher) Service {
	return service{aiClient, repo, lister, sigVerif, hash}
}

func (s service) StoreKeychain(kc account.EndorsedKeychain) error {
	//Checks if the storage must done on this peer
	if err := s.aiClient.CheckStorageAuthorization(kc.Endorsement().TransactionHash()); err != nil {
		return err
	}

	//Checks if the validations matches the required validations for this transaction
	minValids, err := s.aiClient.GetMininumValidations(kc.Endorsement().TransactionHash())
	if err != nil {
		return err
	}
	if len(kc.Endorsement().Validations()) < minValids {
		return ErrInvalidValidationNumber
	}

	//Check the POW
	matchedKey := kc.Endorsement().MasterValidation().ProofOfWorkKey()
	if err := s.sigVerif.VerifyTransactionDataSignature(mining.KeychainTransaction, matchedKey, kc, kc.EmitterSignature()); err != nil {
		return ErrInvalidDataMining
	}

	//Checks signatures
	if err := s.sigVerif.VerifyKeychainSignatures(kc); err != nil {
		return err
	}
	if err := s.verifyEndorsementSignatures(kc.Endorsement()); err != nil {
		return err
	}

	//Check integrity of the keychain
	prevKc, err := s.lister.GetLastKeychain(kc.Address())
	if err != nil {
		return err
	}
	if err := s.checkKeychainEndorsementHash(kc.Endorsement(), kc, prevKc); err != nil {
		return err
	}

	//If the keychain contains any KO validations, it will be stored on the KO database
	if s.isKO(kc.Endorsement()) {
		return s.repo.StoreKOKeychain(kc)
	}

	return s.repo.StoreKeychain(kc)
}

func (s service) StoreID(id account.EndorsedID) error {
	//Checks if the storage must done on this peer
	if err := s.aiClient.CheckStorageAuthorization(id.Endorsement().TransactionHash()); err != nil {
		return err
	}

	//Check the POW
	matchedKey := id.Endorsement().MasterValidation().ProofOfWorkKey()
	if err := s.sigVerif.VerifyTransactionDataSignature(mining.IDTransaction, matchedKey, id, id.EmitterSignature()); err != nil {
		return ErrInvalidDataMining
	}

	//Checks signatures
	if err := s.sigVerif.VerifyIDSignatures(id); err != nil {
		return err
	}
	if err := s.verifyEndorsementSignatures(id.Endorsement()); err != nil {
		return err
	}

	//Check integrity of the biometric
	if err := s.checkIDEndorsementHash(id.Endorsement(), id); err != nil {
		return err
	}

	//Checks if the validations matches the required validations for this transaction
	minValids, err := s.aiClient.GetMininumValidations(id.Endorsement().TransactionHash())
	if err != nil {
		return err
	}
	if len(id.Endorsement().Validations()) < minValids {
		return ErrInvalidValidationNumber
	}

	//If the biometric contains any KO validations, it will be stored on the KO database
	if s.isKO(id.Endorsement()) {
		return s.repo.StoreKOID(id)
	}

	return s.repo.StoreID(id)
}

func (s service) isKO(end mining.Endorsement) bool {
	if end.MasterValidation().ProofOfWorkValidation().Status() == mining.ValidationKO {
		return true
	}
	for _, v := range end.Validations() {
		if v.Status() == mining.ValidationKO {
			return true
		}
	}

	return false
}

func (s service) verifyEndorsementSignatures(end mining.Endorsement) error {

	if err := s.sigVerif.VerifyValidationSignature(end.MasterValidation().ProofOfWorkValidation()); err != nil {
		return ErrInvalidDataMining
	}

	for _, v := range end.Validations() {
		if err := s.sigVerif.VerifyValidationSignature(v); err != nil {
			return ErrInvalidDataMining
		}
	}

	return nil
}

func (s service) checkKeychainEndorsementHash(end mining.Endorsement, kc account.Keychain, prevKc account.Keychain) error {

	if prevKc != nil {
		if end.LastTransactionHash() == "" {
			return ErrInvalidDataIntegrity
		}
		prevHash, err := s.hasher.HashKeychain(prevKc)
		if err != nil {
			return err
		}
		if prevHash != end.LastTransactionHash() {
			return ErrInvalidDataIntegrity
		}
	}

	hash, err := s.hasher.HashKeychain(kc)
	if err != nil {
		return err
	}

	if hash != end.TransactionHash() {
		return ErrInvalidDataIntegrity
	}

	return nil
}

func (s service) checkIDEndorsementHash(end mining.Endorsement, id account.ID) error {
	hash, err := s.hasher.HashID(id)
	if err != nil {
		return err
	}

	if hash != end.TransactionHash() {
		return ErrInvalidDataIntegrity
	}

	return nil
}
