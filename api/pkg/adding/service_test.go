package adding

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/api/pkg"
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
	req := NewAccountCreationRequest("encrypted ID", "encrypted keychain", "sig")

	res, err := s.AddAccount(req)
	assert.Nil(t, err)
	assert.Equal(t, "transaction hash", res.ResultTransactions().ID().TransactionHash())
	assert.Equal(t, "transaction hash", res.ResultTransactions().Keychain().TransactionHash())
	assert.Equal(t, "sig", res.Signature())
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

	req := NewAccountCreationRequest("encrypted ID", "encrypted keychain", "sig")

	_, err := s.AddAccount(req)
	assert.Equal(t, err, errors.New("Invalid signature"))
}

type mockClient struct{}

func (c mockClient) AddAccount(AccountCreationRequest) (AccountCreationResult, error) {
	txID := api.NewTransactionResult("transaction hash", "", "")
	txKeychain := api.NewTransactionResult("transaction hash", "", "")

	res := NewAccountCreationTransactionResult(txID, txKeychain)
	return NewAccountCreationResult(res, "sig"), nil
}

func (c mockClient) GetAccount(encIDHash string) (listing.AccountResult, error) {
	return listing.NewAccountResult("encrypted_aes_key", "encrypted_wallet", "encrypted_address", "sig"), nil
}

func (c mockClient) IsEmitterAuthorized(emPubKey string) error {
	if emPubKey == "em pub key" {
		return nil
	}
	return listing.ErrUnauthorized
}

func (c mockClient) AddSmartContract(ContractCreationRequest) (api.TransactionResult, error) {
	return api.NewTransactionResult("transaction hash", "", ""), nil
}

func (c mockClient) AddContractMessage(ContractMessageCreationRequest) (api.TransactionResult, error) {
	return api.NewTransactionResult("transctionhash", "", ""), nil
}

func (c mockClient) GetSharedKeys() (listing.SharedKeys, error) {
	return listing.NewSharedKeys(
		"robot pv key",
		"robot pub key",
		[]listing.SharedKeyPair{
			listing.NewSharedKeyPair("enc pv key", "pub key"),
		}), nil
}

func (c mockClient) GetTransactionStatus(addr string, txHash string) (listing.TransactionStatus, error) {
	return listing.TransactionSuccess, nil
}

func (c mockClient) GetContractState(addr string) (listing.ContractState, error) {
	return listing.NewContractState("", "sig"), nil
}

type mockSigVerifier struct {
	isInvalid bool
}

func (v mockSigVerifier) VerifyAccountCreationRequestSignature(req AccountCreationRequest, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyAccountCreationResultSignature(req AccountCreationResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyAccountResultSignature(res listing.AccountResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyContractStateSignature(state listing.ContractState, pubKey string) error {
	return nil
}

func (v mockSigVerifier) SignAccountCreationResult(res AccountCreationResult, pvKey string) (AccountCreationResult, error) {
	return NewAccountCreationResult(
		NewAccountCreationTransactionResult(
			api.NewTransactionResult("transaction hash", "ip", "sig"),
			api.NewTransactionResult("transaction hash", "ip", "sig"),
		), "sig",
	), nil
}

func (v mockSigVerifier) VerifyHashSignature(data string, pubKey string, sig string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyCreationTransactionResultSignature(res api.TransactionResult, pubKey string) error {
	if v.isInvalid {
		return errors.New("Invalid signature")
	}
	return nil
}

func (v mockSigVerifier) VerifyContractCreationRequestSignature(req ContractCreationRequest, key string) error {
	return nil
}

func (v mockSigVerifier) VerifyContractMessageCreationRequestSignature(req ContractMessageCreationRequest, key string) error {
	return nil
}
