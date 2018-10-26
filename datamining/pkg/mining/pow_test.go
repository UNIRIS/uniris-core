package mining

import (
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"

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

	repo := &mockDatabase{}
	list := listing.NewService(repo)

	pow := NewPOW(list, mockPowSigner{}, "my key", "my pv key")
	lastValidPool := pool.PeerCluster{
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

type mockDatabase struct {
	BioWallets []*datamining.BioWallet
	Wallets    []*datamining.Wallet
}

func (d *mockDatabase) FindBioWallet(bh string) (*datamining.BioWallet, error) {
	return nil, nil
}

func (d *mockDatabase) FindWallet(addr string) (*datamining.Wallet, error) {
	return nil, nil
}

func (d *mockDatabase) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

func (d *mockDatabase) StoreWallet(w *datamining.Wallet) error {
	d.Wallets = append(d.Wallets, w)
	return nil
}

func (d *mockDatabase) StoreBioWallet(bw *datamining.BioWallet) error {
	d.BioWallets = append(d.BioWallets, bw)
	return nil
}

type mockPowSigner struct{}

func (s mockPowSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockPowSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}
