package leading

import (
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Execute the POW
	Given a biod key
	When I execute the POW, and I lookup the tech repository to find it
	Then I get a master validation
*/
func TestExecutePOW(t *testing.T) {

	repo := mockTechRepo{}
	repo.BioKeys = append(repo.BioKeys, "bioKey1")
	repo.BioKeys = append(repo.BioKeys, "bioKey2")
	repo.BioKeys = append(repo.BioKeys, "bioKey3")

	pow := NewPOW(repo, mockPowSigner{}, "my key", "my pv key")
	valid, err := pow.Execute("hash", "signature", nil)
	assert.Nil(t, err)
	assert.NotNil(t, valid)

	assert.Equal(t, "my key", valid.ProofOfWorkRobotKey())
	assert.Equal(t, "my key", valid.ProofOfWorkValidation().PublicKey())
	assert.Equal(t, datamining.ValidationOK, valid.ProofOfWorkValidation().Status())
}

type mockTechRepo struct {
	BioKeys []string
}

func (r mockTechRepo) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

type mockPowSigner struct{}

func (s mockPowSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockPowSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}
