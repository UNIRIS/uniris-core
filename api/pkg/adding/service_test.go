package adding

import (
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
		val:          mockGoodRequestValidator{},
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
		val:          mockBadRequestValidator{},
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
	assert.Equal(t, err, ErrInvalidSignature)
}

type mockClient struct{}

func (c mockClient) AddAccount(AccountCreationRequest) (*AccountCreationResult, error) {
	return &AccountCreationResult{
		Transactions: AccountCreationTransactions{
			Biometric: CreationTransaction{
				TransactionHash: "transaction hash",
			},
			Keychain: CreationTransaction{
				TransactionHash: "transaction hash",
			},
		},
		Signature: "signature of the response",
	}, nil
}

type mockGoodRequestValidator struct{}

func (v mockGoodRequestValidator) CheckDataSignature(data interface{}, pubKey string, sig string) (bool, error) {
	return true, nil
}

type mockBadRequestValidator struct{}

func (v mockBadRequestValidator) CheckDataSignature(data interface{}, pubKey string, sig string) (bool, error) {
	return false, nil
}
