package listing

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Get account's details from the robot
	Given a encrypted public key and a signature
	When I want to get the account details
	Then I can get the encrypted data from the roboto
*/
func TestGetAccount(t *testing.T) {
	s := service{
		client:       mockClient{},
		sigVerif:     mockGoodSignatureVerif{},
		sharedBioPub: "my key",
	}

	res, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Nil(t, err)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESKey)
	assert.Equal(t, "encrypted_wallet", res.EncryptedWallet)
	assert.Equal(t, "encrypted_address", res.EncryptedAddress)
}

/*
Scenario: Catch invalid signature when get account's details from the robot
	Given a encrypted public key and a invalid signature
	When I want to get the account details
	Then I get an error
*/
func TestGetAccountInvalidSig(t *testing.T) {
	s := service{
		client:       mockClient{},
		sigVerif:     mockBadSignatureVerif{},
		sharedBioPub: "my key",
	}

	_, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Equal(t, err, errors.New("Invalid signature"))
}

type mockClient struct{}

func (c mockClient) GetAccount(encHash string) (*AccountResult, error) {
	return &AccountResult{
		EncryptedAESKey:  "encrypted_aes_key",
		EncryptedAddress: "encrypted_address",
		EncryptedWallet:  "encrypted_wallet",
		Signature:        "sig",
	}, nil
}

type mockGoodSignatureVerif struct{}

func (v mockGoodSignatureVerif) VerifyHashSignature(data string, pubKey string, sig string) error {
	return nil
}

type mockBadSignatureVerif struct{}

func (v mockBadSignatureVerif) VerifyHashSignature(data string, pubKey string, sig string) error {
	return errors.New("Invalid signature")
}
