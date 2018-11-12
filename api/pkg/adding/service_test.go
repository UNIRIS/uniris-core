package adding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Enroll an user
	Given a encrypted public key and a signature
	When I want to get the account details
	Then I can get the encrypted data from the roboto
*/
func TestAddAccount(t *testing.T) {
	s := service{
		client:       mockClient{},
		sigVerif:     mockGoodSignatureVerifer{},
		sharedBioPub: "my key",
	}

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
Scenario: Catch invalid signature when get account's details from the robot
	Given a encrypted public key and a invalid signature
	When I want to get the account details
	Then I get an error
*/
func TestAddAccountInvalidSig(t *testing.T) {
	s := service{
		client:       mockClient{},
		sigVerif:     mockBadSignatureVerifer{},
		sharedBioPub: "my key",
	}

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

type mockClient struct{}

func (c mockClient) AddAccount(AccountCreationRequest) (*AccountCreationResult, error) {
	return &AccountCreationResult{
		Transactions: AccountCreationTransactions{
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

type mockGoodSignatureVerifer struct{}

func (v mockGoodSignatureVerifer) VerifyAccountCreationRequestSignature(data AccountCreationRequest, pubKey string) error {
	return nil
}

type mockBadSignatureVerifer struct{}

func (v mockBadSignatureVerifer) VerifyAccountCreationRequestSignature(data AccountCreationRequest, pubKey string) error {
	return errors.New("Invalid signature")
}
