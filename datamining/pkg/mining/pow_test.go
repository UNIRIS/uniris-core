package mining

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/emitter"
	emlisting "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
)

/*
Scenario: Execute the POW
	Given a biod key
	When I execute the POW, and I lookup the tech repository to find it
	Then I get a master validation
*/
func TestExecutePOW(t *testing.T) {

	repo := &mockDatabase{}
	emLister := emlisting.NewService(repo)

	lastValidPool := datamining.NewPool(datamining.Peer{PublicKey: "key"})

	pow := pow{
		lastVPool:   lastValidPool,
		emLister:    emLister,
		robotPubKey: "my key",
		robotPvKey:  "my key",
		signer:      mockPowSigner{},
		txEmSig:     "signature",
		txData:      "data",
		txType:      KeychainTransaction,
	}

	valid, err := pow.execute()
	assert.Nil(t, err)
	assert.NotNil(t, valid)

	assert.Equal(t, "key1", valid.ProofOfWorkKey())
	assert.Equal(t, "my key", valid.ProofOfWorkValidation().PublicKey())
	assert.Equal(t, ValidationOK, valid.ProofOfWorkValidation().Status())
}

/*
Scenario: Execute the POW and not find a match
	Given a biod key
	When I execute the POW, and I lookup the tech repository to find it
	Then I not find it and return the validation but with KO
*/
func TestExecutePOW_KO(t *testing.T) {

	repo := &mockDatabase{}
	emLister := emlisting.NewService(repo)

	lastValidPool := datamining.NewPool(datamining.Peer{PublicKey: "key"})

	pow := pow{
		lastVPool:   lastValidPool,
		emLister:    emLister,
		robotPubKey: "my key",
		robotPvKey:  "my key",
		signer:      mockBadPowSigner{},
		txEmSig:     "signature",
		txData:      "data",
		txType:      KeychainTransaction,
	}

	valid, err := pow.execute()
	assert.Nil(t, err)
	assert.NotNil(t, valid)

	assert.Equal(t, "", valid.ProofOfWorkKey())
	assert.Equal(t, "my key", valid.ProofOfWorkValidation().PublicKey())
	assert.Equal(t, ValidationKO, valid.ProofOfWorkValidation().Status())
}

type mockDatabase struct {
}

func (d *mockDatabase) ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error) {
	return []emitter.SharedKeyPair{
		emitter.SharedKeyPair{
			PublicKey: "key1",
		}, emitter.SharedKeyPair{
			PublicKey: "key2",
		}, emitter.SharedKeyPair{
			PublicKey: "key3",
		}}, nil
}

type mockPowSigner struct{}

func (s mockPowSigner) VerifyTransactionDataSignature(txType TransactionType, pubk string, data interface{}, der string) error {
	return nil
}

func (s mockPowSigner) SignValidation(v Validation, pvKey string) (Validation, error) {
	return NewValidation(v.Status(), v.Timestamp(), v.PublicKey(), "sig"), nil
}

type mockBadPowSigner struct{}

func (s mockBadPowSigner) VerifyTransactionDataSignature(txType TransactionType, pubk string, data interface{}, der string) error {
	return errors.New("Invalid signature")
}

func (s mockBadPowSigner) SignValidation(v Validation, pvKey string) (Validation, error) {
	return NewValidation(v.Status(), v.Timestamp(), v.PublicKey(), "sig"), nil
}
