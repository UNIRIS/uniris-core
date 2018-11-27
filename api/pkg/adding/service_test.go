package adding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

/*
Scenario: Enroll an user
	Given a encrypted public key and a signature
	When I want to get the account details
	Then I can get the encrypted data from the roboto
*/
func TestAddAccount(t *testing.T) {
	c := mockClient{}
	sig := mockSigVerifier{}
	l := listing.NewService(c, sig)
	s := NewService(l, c, sig)
	req := &AccountCreationRequest{
		EncryptedID:       "encrypted ID",
		EncryptedKeychain: "encrypted keychain",
		Signature:         "signature request",
	}

	res, err := s.AddAccount(req)
	assert.Nil(t, err)
	assert.Equal(t, "transaction hash", res.Transactions.ID.TransactionHash)
	assert.Equal(t, "transaction hash", res.Transactions.Keychain.TransactionHash)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Catch invalid signature when get account's details from the robot
	Given a encrypted public key and a invalid signature
	When I want to get the account details
	Then I get an error
*/
func TestAddAccountInvalidSig(t *testing.T) {
	c := mockClient{}
	sig := mockSigVerifier{isInvalid: true}
	l := listing.NewService(c, sig)
	s := NewService(l, c, sig)

	req := &AccountCreationRequest{
		EncryptedID:       "encrypted bio data",
		EncryptedKeychain: "encrypted wallet data",
		Signature:         "signature",
	}

	_, err := s.AddAccount(req)
	assert.Equal(t, err, errors.New("Invalid signature"))
}

type mockClient struct{}

func (c mockClient) AddAccount(*AccountCreationRequest) (*AccountCreationResult, error) {
	return &AccountCreationResult{
		Transactions: AccountCreationTransactionsResult{
			ID: TransactionResult{
				TransactionHash: "transaction hash",
			},
			Keychain: TransactionResult{
				TransactionHash: "transaction hash",
			},
		},
		Signature: "sig",
	}, nil
}

func (c mockClient) GetAccount(encIDHash string) (*listing.AccountResult, error) {
	return &listing.AccountResult{
		EncryptedAESKey:  "encrypted_aes_key",
		EncryptedAddress: "encrypted_address",
		EncryptedWallet:  "encrypted_wallet",
		Signature:        "sig",
	}, nil
}

func (c mockClient) IsEmitterAuthorized(emPubKey string) error {
	if emPubKey == "em pub key" {
		return nil
	}
	return listing.ErrUnauthorized
}

func (c mockClient) GetSharedKeys() (*listing.SharedKeysResult, error) {
	return &listing.SharedKeysResult{
		RobotPublicKey: "robot key",
		EmitterKeys: []listing.SharedKeyPair{
			listing.SharedKeyPair{
				EncryptedPrivateKey: "enc pv key",
				PublicKey:           "pub key",
			},
		},
	}, nil
}

type mockSigVerifier struct {
	isInvalid bool
}

func (v mockSigVerifier) VerifyAccountCreationRequestSignature(req *AccountCreationRequest, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyAccountCreationResultSignature(req *AccountCreationResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyAccountResultSignature(res *listing.AccountResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) SignAccountCreationResult(res *AccountCreationResult, pvKey string) error {
	res.Signature = "sig"
	return nil
}

func (v mockSigVerifier) VerifyHashSignature(data string, pubKey string, sig string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyCreationTransactionResultSignature(res TransactionResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}
