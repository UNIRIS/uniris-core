package adding

import (
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//Service defines methods to adding to the blockchain
type Service interface {
	AddAccount(AccountCreationRequest) (AccountCreationResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	AddAccount(AccountCreationRequest) (AccountCreationResult, error)
}

//Signer defines methods to handle signature
type Signer interface {
	signatureVerifier
	signatureBuilder
}

type signatureVerifier interface {

	//VerifyAccountCreationRequestSignature checks the signature of the account creation request
	VerifyAccountCreationRequestSignature(req AccountCreationRequest, key string) error

	//VerifyCreationTransactionResultSignature checks the signature of a creation transaction result
	VerifyCreationTransactionResultSignature(res TransactionResult, pubKey string) error
}

type signatureBuilder interface {

	//SignAccountCreationResult signs the account creation result
	SignAccountCreationResult(data AccountCreationResult, key string) (AccountCreationResult, error)
}

type service struct {
	lister listing.Service
	client RobotClient
	sig    Signer
}

//NewService creates a new adding service
func NewService(lister listing.Service, client RobotClient, sig Signer) Service {
	return service{lister, client, sig}
}

func (s service) AddAccount(req AccountCreationRequest) (AccountCreationResult, error) {
	keys, err := s.lister.GetSafeSharedKeys()
	if err != nil {
		return nil, err
	}
	if err := s.sig.VerifyAccountCreationRequestSignature(req, keys.RequestPublicKey()); err != nil {
		return nil, err
	}

	res, err := s.client.AddAccount(req)
	if err != nil {
		return nil, err
	}

	if err := s.sig.VerifyCreationTransactionResultSignature(res.ResultTransactions().ID(), keys.RobotPublicKey()); err != nil {
		return nil, err
	}

	if err := s.sig.VerifyCreationTransactionResultSignature(res.ResultTransactions().Keychain(), keys.RobotPublicKey()); err != nil {
		return nil, err
	}

	res, err = s.sig.SignAccountCreationResult(res, keys.RobotPrivateKey())
	if err != nil {
		return nil, err
	}
	return res, nil
}
