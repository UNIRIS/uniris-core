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
	StoreKeychain(account.Keychain) error

	//StoreBiometric persists the biometric
	StoreBiometric(account.Biometric) error

	//StoreKOKeychain persists the keychain in the KO database
	StoreKOKeychain(account.Keychain) error

	//StoreKOBiometric persists the biometric in the KO database
	StoreKOBiometric(account.Biometric) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {

	//StoreKeychain processes the keychain storage
	//
	//It performs checks to insure the integrity of the keychain
	//Determines if the keychain must be store in the KO, OK or pending database
	StoreKeychain(account.Keychain) error

	//StoreBiometric processes the biometric storage
	//
	//It performs checks to insure the integrity of the biometric
	//Determines if the biometric must be store in the KO, OK or pending database
	StoreBiometric(account.Biometric) error
}

type signatureVerifier interface {
	account.KeychainSignatureVerifier
	account.BiometricSignatureVerifier
	mining.ValidationVerifier
	mining.PowSigVerifier
}

type hasher interface {
	account.KeychainHasher
	account.BiometricHasher
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

func (s service) StoreKeychain(kc account.Keychain) error {
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
	biodSig := kc.Signatures().Biod()
	if err := s.sigVerif.VerifyTransactionDataSignature(mining.KeychainTransaction, matchedKey, kc, biodSig); err != nil {
		return ErrInvalidDataMining
	}

	//Checks signatures
	if err := s.sigVerif.VerifyKeychainDataSignatures(kc); err != nil {
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

func (s service) StoreBiometric(bio account.Biometric) error {
	//Checks if the storage must done on this peer
	if err := s.aiClient.CheckStorageAuthorization(bio.Endorsement().TransactionHash()); err != nil {
		return err
	}

	//Check the POW
	matchedKey := bio.Endorsement().MasterValidation().ProofOfWorkKey()
	biodSig := bio.Signatures().Biod()
	if err := s.sigVerif.VerifyTransactionDataSignature(mining.BiometricTransaction, matchedKey, bio, biodSig); err != nil {
		return ErrInvalidDataMining
	}

	//Checks signatures
	if err := s.sigVerif.VerifyBiometricDataSignatures(bio); err != nil {
		return err
	}
	if err := s.verifyEndorsementSignatures(bio.Endorsement()); err != nil {
		return err
	}

	//Check integrity of the biometric
	if err := s.checkBiometricEndorsementHash(bio.Endorsement(), bio); err != nil {
		return err
	}

	//Checks if the validations matches the required validations for this transaction
	minValids, err := s.aiClient.GetMininumValidations(bio.Endorsement().TransactionHash())
	if err != nil {
		return err
	}
	if len(bio.Endorsement().Validations()) < minValids {
		return ErrInvalidValidationNumber
	}

	//If the biometric contains any KO validations, it will be stored on the KO database
	if s.isKO(bio.Endorsement()) {
		return s.repo.StoreKOBiometric(bio)
	}

	return s.repo.StoreBiometric(bio)
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

func (s service) checkKeychainEndorsementHash(end mining.Endorsement, data account.KeychainData, previousData account.KeychainData) error {

	if previousData != nil {
		if end.LastTransactionHash() == "" {
			return ErrInvalidDataIntegrity
		}
		prevHash, err := s.hasher.HashKeychainData(previousData)
		if err != nil {
			return err
		}
		if prevHash != end.LastTransactionHash() {
			return ErrInvalidDataIntegrity
		}
	}

	hash, err := s.hasher.HashKeychainData(data)
	if err != nil {
		return err
	}

	if hash != end.TransactionHash() {
		return ErrInvalidDataIntegrity
	}

	return nil
}

func (s service) checkBiometricEndorsementHash(end mining.Endorsement, data account.BiometricData) error {
	hash, err := s.hasher.HashBiometricData(data)
	if err != nil {
		return err
	}

	if hash != end.TransactionHash() {
		return ErrInvalidDataIntegrity
	}

	return nil
}
