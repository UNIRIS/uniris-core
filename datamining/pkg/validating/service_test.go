package validating

import (
	"testing"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Validate wallet data
	Given a wallet data
	When I want validate it
	Then I get a validated transaction
*/
func TestValidateWallet(t *testing.T) {

	srv := NewService(mockSigner{}, "robotKey", "robotPvKey")

	w := &datamining.WalletData{
		BiodPubk: "pubKey",
		EmPubk:   "pubKey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	v, err := srv.ValidateWalletData(w)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "robotKey", v.PublicKey())
}

/*
Scenario: Validate bio data
	Given a bio data
	When I want validate it
	Then I get a validated transaction
*/
func TestValidatBio(t *testing.T) {

	srv := NewService(mockSigner{}, "robotKey", "robotPvKey")

	b := &datamining.BioData{
		BiodPubk: "pubkey",
		EmPubk:   "pubkey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	v, err := srv.ValidateBioData(b)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "robotKey", v.PublicKey())
}

type mockSigner struct{}

func (s mockSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return nil
}
