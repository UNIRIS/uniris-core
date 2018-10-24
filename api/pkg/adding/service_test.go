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

	req := EnrollmentRequest{
		EncryptedBioData:    "encrypted bio data",
		EncryptedWalletData: "encrypted wallet data",
		SignatureRequest:    "signature request",
		SignaturesBio: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
		SignaturesWallet: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
	}

	res, err := s.AddAccount(req)
	assert.Nil(t, err)
	assert.Equal(t, "transaction hash", res.Transactions.Biod)
	assert.Equal(t, "transaction hash", res.Transactions.Data)
	assert.Equal(t, "signature of the response", res.SignatureRequest)
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

	req := EnrollmentRequest{
		EncryptedBioData:    "encrypted bio data",
		EncryptedWalletData: "encrypted wallet data",
		SignatureRequest:    "signature request",
		SignaturesBio: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
		SignaturesWallet: Signatures{
			BiodSig:   "biod signature",
			PersonSig: "person sig",
		},
	}

	_, err := s.AddAccount(req)
	assert.Equal(t, err, ErrInvalidSignature)
}

type mockClient struct{}

func (c mockClient) AddAccount(EnrollmentRequest) (*EnrollmentResult, error) {
	return &EnrollmentResult{
		Transactions: EnrollmentTransactions{
			Biod: "transaction hash",
			Data: "transaction hash",
		},
		SignatureRequest: "signature of the response",
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
