package master

import (
	"errors"
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Execute the POW
	Given a biod key
	When I execute the POW, and I lookup the tech repository to find it
	Then I get a master validation
*/
func TestExecutePOW(t *testing.T) {

	repo := &mockDatabase{}
	list := listing.NewService(repo)

	pow := NewPOW(list, mockPowSigner{}, "my key", "my pv key")
	lastValidPool := pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{PublicKey: "key"},
		},
	}
	valid, err := pow.Execute("hash", "signature", lastValidPool)
	assert.Nil(t, err)
	assert.NotNil(t, valid)

	assert.Equal(t, "my key", valid.ProofOfWorkRobotKey())
	assert.Equal(t, "my key", valid.ProofOfWorkValidation().PublicKey())
	assert.Equal(t, datamining.ValidationOK, valid.ProofOfWorkValidation().Status())
}

/*
Scenario: Execute the POW and not find a match
	Given a biod key
	When I execute the POW, and I lookup the tech repository to find it
	Then I not find it and return the validation but with KO
*/
func TestExecutePOW_KO(t *testing.T) {

	repo := &mockDatabase{}
	list := listing.NewService(repo)

	pow := NewPOW(list, mockBadPowSigner{}, "my key", "my pv key")
	lastValidPool := pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{PublicKey: "key"},
		},
	}
	valid, err := pow.Execute("hash", "signature", lastValidPool)
	assert.Nil(t, err)
	assert.NotNil(t, valid)

	assert.Equal(t, "my key", valid.ProofOfWorkRobotKey())
	assert.Equal(t, "my key", valid.ProofOfWorkValidation().PublicKey())
	assert.Equal(t, datamining.ValidationKO, valid.ProofOfWorkValidation().Status())
}

type mockDatabase struct {
	Biometrics []*datamining.Biometric
	Keychains  []*datamining.Keychain
}

func (d *mockDatabase) FindBiometric(bh string) (*datamining.Biometric, error) {
	return nil, nil
}

func (d *mockDatabase) FindKeychain(addr string) (*datamining.Keychain, error) {
	return nil, nil
}

func (d *mockDatabase) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

func (d *mockDatabase) StoreKeychain(w *datamining.Keychain) error {
	d.Keychains = append(d.Keychains, w)
	return nil
}

func (d *mockDatabase) StoreBiometric(b *datamining.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	return nil
}

type mockPowSigner struct{}

func (s mockPowSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockPowSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockBadPowSigner struct{}

func (s mockBadPowSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return errors.New("Invalid signature")
}

func (s mockBadPowSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}
