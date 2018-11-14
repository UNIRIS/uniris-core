package adding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Enroll an user
	Given a account creation request and a signature
	When I want to get the account details
	Then I can get the encrypted data from the roboto
*/
func TestAddAccount(t *testing.T) {
	s := NewService("my key", mockClient{}, mockGoodSignatureVerifer{})

	req := AccountCreationRequest{
		EncryptedBioData:      "encrypted bio data",
		EncryptedKeychainData: "encrypted wallet data",
		SignatureRequest:      "signature request",
		SignaturesBio: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
		SignaturesKeychain: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
	}

	res, err := s.AddAccount(req)
	assert.Nil(t, err)
	assert.Equal(t, "transaction hash", res.Transactions.Biometric.TransactionHash)
	assert.Equal(t, "transaction hash", res.Transactions.Keychain.TransactionHash)
	assert.Equal(t, "signature of the response", res.Signature)
}

/*
Scenario: Catch invalid signature when create account's
	Given a account creation request and an invalid signature
	When I want to create an account
	Then I get an error
*/
func TestAddAccountInvalidSig(t *testing.T) {
	s := NewService("my key", mockClient{}, mockBadSignatureVerifer{})

	req := AccountCreationRequest{
		EncryptedBioData:      "encrypted bio data",
		EncryptedKeychainData: "encrypted wallet data",
		SignatureRequest:      "signature request",
		SignaturesBio: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
		SignaturesKeychain: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
	}

	_, err := s.AddAccount(req)
	assert.Equal(t, err, errors.New("Invalid signature"))
}

/*
Scenario: Register biometric device public key
	Given a encrypted biometric public key
	When I want to store
	Then I the key is stored
*/
func TestRegisterBiod(t *testing.T) {
	s := NewService("my key", mockClient{}, mockGoodSignatureVerifer{})

	req := BiodRegisterRequest{
		EncryptedPublicKey: "pub key",
		Signature:          "sig",
	}

	res, err := s.RegisterBiod(req)
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.PublicKeyHash)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Catch invalid signature when register biometric device
	Given a encrypted public key and a invalid signature
	When I want to store the biometric device
	Then I get an error
*/
func TestRegisterBiodInvalidSig(t *testing.T) {
	s := NewService("my key", mockClient{}, mockBadSignatureVerifer{})

	req := BiodRegisterRequest{
		EncryptedPublicKey: "pub key",
		Signature:          "sig",
	}

	_, err := s.RegisterBiod(req)
	assert.Equal(t, "Invalid signature", err.Error())
}

type mockClient struct{}

func (c mockClient) AddAccount(AccountCreationRequest) (*AccountCreationResult, error) {
	return &AccountCreationResult{
		Transactions: AccountCreationTransactionsResult{
			Biometric: TransactionResult{
				TransactionHash: "transaction hash",
			},
			Keychain: TransactionResult{
				TransactionHash: "transaction hash",
			},
		},
		Signature: "signature of the response",
	}, nil
}

func (mockClient) RegisterBiod(encPubKey string) (*BiodRegisterResponse, error) {
	return &BiodRegisterResponse{
		PublicKeyHash: "hash",
		Signature:     "sig",
	}, nil
}

type mockGoodSignatureVerifer struct{}

func (v mockGoodSignatureVerifer) VerifyAccountCreationRequestSignature(data AccountCreationRequest, pubKey string) error {
	return nil
}
func (v mockGoodSignatureVerifer) VerifyBiodRegisteringRequestSignature(req BiodRegisterRequest, pubKey string) error {
	return nil
}

type mockBadSignatureVerifer struct{}

func (v mockBadSignatureVerifer) VerifyAccountCreationRequestSignature(data AccountCreationRequest, pubKey string) error {
	return errors.New("Invalid signature")
}

func (v mockBadSignatureVerifer) VerifyBiodRegisteringRequestSignature(req BiodRegisterRequest, pubKey string) error {
	return errors.New("Invalid signature")
}
