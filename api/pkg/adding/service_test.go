package adding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/api/pkg/system"
)

/*
Scenario: Enroll an user
	Given a encrypted public key and a signature
	When I want to get the account details
	Then I can get the encrypted data from the roboto
*/
func TestAddAccount(t *testing.T) {
	s := service{
		client:   mockClient{},
		sigVerif: mockGoodSignatureVerifer{},
		conf: system.UnirisConfig{
			SharedKeys: system.SharedKeys{
				EmKeys: []system.KeyPair{
					system.KeyPair{
						PublicKey: "my key",
					},
				},
			},
		},
	}

	req := AccountCreationRequest{
		EncryptedID:       "encrypted ID",
		EncryptedKeychain: "encrypted keychain",
		Signature:         "signature request",
	}

	res, err := s.AddAccount(req)
	assert.Nil(t, err)
	assert.Equal(t, "transaction hash", res.Transactions.ID.TransactionHash)
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
		client:   mockClient{},
		sigVerif: mockBadSignatureVerifer{},
		conf: system.UnirisConfig{
			SharedKeys: system.SharedKeys{
				EmKeys: []system.KeyPair{
					system.KeyPair{
						PublicKey: "my key",
					},
				},
			},
		},
	}

	req := AccountCreationRequest{
		EncryptedID:       "encrypted bio data",
		EncryptedKeychain: "encrypted wallet data",
		Signature:         "signature",
	}

	_, err := s.AddAccount(req)
	assert.Equal(t, err, errors.New("Invalid signature"))
}

type mockClient struct{}

func (c mockClient) AddAccount(AccountCreationRequest) (*AccountCreationResult, error) {
	return &AccountCreationResult{
		Transactions: AccountCreationTransactionsResult{
			ID: TransactionResult{
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
