package listing

import (
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
		val:          mockGoodRequestValidator{},
		sharedBioPub: []byte("my key"),
	}

	res, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Nil(t, err)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESKey)
	assert.Equal(t, "encrypted_wallet", res.EncryptedWallet)
	assert.Equal(t, "addr_wallet_person", res.EncryptedAddrPerson)
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
		val:          mockBadRequestValidator{},
		sharedBioPub: []byte("my key"),
	}

	_, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Equal(t, err, ErrInvalidSignature)
}

type mockClient struct{}

func (c mockClient) GetAccount(AccountRequest) (AccountResult, error) {
	return AccountResult{
		EncryptedAESKey:     "encrypted_aes_key",
		EncryptedAddrPerson: "addr_wallet_person",
		EncryptedWallet:     "encrypted_wallet",
	}, nil
}

type mockGoodRequestValidator struct{}

func (v mockGoodRequestValidator) CheckSignature(data interface{}, pubKey []byte, sig []byte) (bool, error) {
	return true, nil
}

type mockBadRequestValidator struct{}

func (v mockBadRequestValidator) CheckSignature(data interface{}, pubKey []byte, sig []byte) (bool, error) {
	return false, nil
}
