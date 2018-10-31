package slave

import (
	"testing"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave/checks"
)

/*
Scenario: Validate wallet data
	Given a wallet data
	When I want validate it
	Then I get a validated transaction
*/
func TestValidateWallet(t *testing.T) {

	srv := service{
		robotKey: "key",
		sig:      mockSigner{},
		checks: map[datamining.TransactionType][]checks.Handler{
			datamining.CreateKeychainTransaction: []checks.Handler{
				checks.NewSignatureChecker(mockSigner{}),
			},
		},
	}
	k := &datamining.KeyChainData{
		BiodPubk:   "pubKey",
		PersonPubk: "pubKey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	v, err := srv.Validate(k, datamining.CreateKeychainTransaction)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "key", v.PublicKey())
}

/*
Scenario: Validate an invalid transaction
	Given a invalid transaction
	When we validate it
	Then we get a validation with a KO status
*/
func TestValidateWalletWithKO(t *testing.T) {
	srv := service{
		robotKey: "key",
		sig:      mockSigner{},
		checks: map[datamining.TransactionType][]checks.Handler{
			datamining.CreateKeychainTransaction: []checks.Handler{
				checks.NewSignatureChecker(mockBadSigner{}),
			},
		},
	}

	k := &datamining.KeyChainData{
		BiodPubk:   "pubKey",
		PersonPubk: "pubKey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	v, err := srv.Validate(k, datamining.CreateKeychainTransaction)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationKO, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "key", v.PublicKey())
}

/*
Scenario: Validate bio data
	Given a bio data
	When I want validate it
	Then I get a validated transaction
*/
func TestValidatBio(t *testing.T) {
	srv := service{
		robotKey: "key",
		sig:      mockSigner{},
		checks: map[datamining.TransactionType][]checks.Handler{
			datamining.CreateBioTransaction: []checks.Handler{
				checks.NewSignatureChecker(mockSigner{}),
			},
		},
	}

	b := &datamining.BioData{
		BiodPubk:   "pubkey",
		PersonPubk: "pubkey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	v, err := srv.Validate(b, datamining.CreateBioTransaction)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "key", v.PublicKey())
}

type mockSigner struct{}

func (s mockSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return nil
}

type mockBadSigner struct{}

func (s mockBadSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockBadSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return checks.ErrInvalidSignature
}
