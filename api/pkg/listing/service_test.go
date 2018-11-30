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
	s := NewService(mockClient{}, mockSigVerifier{})

	res, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Nil(t, err)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESKey())
	assert.Equal(t, "encrypted_wallet", res.EncryptedWallet())
	assert.Equal(t, "encrypted_address", res.EncryptedAddress())
}

/*
Scenario: Catch invalid signature when get account's details from the robot
	Given a encrypted public key and a invalid signature
	When I want to get the account details
	Then I get an error
*/
func TestGetAccountInvalidSig(t *testing.T) {
	s := NewService(mockClient{}, mockSigVerifier{isInvalid: true})
	_, err := s.GetAccount("encrypted person pub key", "sig")
	assert.Equal(t, err, errors.New("Invalid signature"))
}

/*
Scenario: Get the shared keys
	Given emitter publi ckey and signature
	When I want to get the shared keys
	Then I get the shared keys
*/
func TestGetSharedKeys(t *testing.T) {
	s := NewService(mockClient{}, mockSigVerifier{})
	res, err := s.GetSharedKeys("em pub key", "sig")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "robot pub key", res.RobotPublicKey())
	assert.Equal(t, "enc pv key", res.EmitterKeyPairs()[0].EncryptedPrivateKey())
}

/*
Scenario: Catch invalid signature when get the shared keys
	Given emitter public key and invalid signature
	When I want to get the shared keys
	Then I get an error
*/
func TestInvalidSigGetSharedKeys(t *testing.T) {
	s := NewService(mockClient{}, mockSigVerifier{isInvalid: true})
	_, err := s.GetSharedKeys("em key", "invalid sig")
	assert.Equal(t, "Invalid signature", err.Error())
}

/*
Scenario: Get the shared keys using unauthorized public key
	Given an unauthorized emitter public key and signature
	When I want to get the shared keys
	Then I get an error
*/
func TestGetSharedKeysWithUnauthorized(t *testing.T) {
	s := NewService(mockClient{}, mockSigVerifier{})
	_, err := s.GetSharedKeys("invalid key", "sig")
	assert.Equal(t, ErrUnauthorized, err)
}

type mockClient struct{}

func (c mockClient) GetAccount(encIDHash string) (AccountResult, error) {
	return NewAccountResult("encrypted_aes_key", "encrypted_wallet", "encrypted_address", "sig"), nil
}

func (c mockClient) IsEmitterAuthorized(emPubKey string) error {
	if emPubKey == "em pub key" {
		return nil
	}
	return ErrUnauthorized
}

func (c mockClient) GetSharedKeys() (SharedKeys, error) {
	return NewSharedKeys(
		"robot pv key",
		"robot pub key",
		[]SharedKeyPair{
			NewSharedKeyPair("enc pv key", "pub key"),
		}), nil
}

type mockSigVerifier struct {
	isInvalid bool
}

func (v mockSigVerifier) VerifyHashSignature(data string, pubKey string, sig string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyAccountResultSignature(res AccountResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}
