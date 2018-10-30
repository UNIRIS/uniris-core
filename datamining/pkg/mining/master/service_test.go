package master

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

/*
Scenario: Lead wallet creation
	Given a wallet data
	When I want to lead the wallet creation and validation
	Then I perform POW and process validations and create the wallet
*/
func TestLeadWallet(t *testing.T) {
	repo := &mockSrvDatabase{
		Biometrics: make([]*datamining.Biometric, 0),
		Keychains:  make([]*datamining.Keychain, 0),
	}
	srv := NewService(
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockHasher{},
		listing.NewService(repo),
		"robotPubKey",
		"robotPvKey")

	w := &datamining.KeyChainData{
		BiodPubk:        "pubkey",
		PersonPubk:      "pubkey",
		CipherAddrRobot: "addr",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err := srv.LeadMining("hash", "addr", "fake sig", w, datamining.CreateKeychainTransaction)
	assert.Nil(t, err)
	assert.Len(t, repo.Keychains, 1)
	assert.Equal(t, "hash", repo.Keychains[0].Endorsement().TransactionHash())
	assert.Equal(t, "robotPubKey", repo.Keychains[0].Endorsement().MasterValidation().ProofOfWorkRobotKey())
	assert.Len(t, repo.Keychains[0].Endorsement().Validations(), 1)
	assert.Equal(t, "pubkey", repo.Keychains[0].Endorsement().Validations()[0].PublicKey())
}

/*
Scenario: Lead bio validation
	Given a bio data
	When I want to lead the data validation
	Then I perform POW and returns validated transaction
*/
func TestLeadBio(t *testing.T) {
	repo := &mockSrvDatabase{
		Biometrics: make([]*datamining.Biometric, 0),
		Keychains:  make([]*datamining.Keychain, 0),
	}
	srv := NewService(
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockHasher{},
		listing.NewService(repo),
		"robotPubKey",
		"robotPvKey")

	bd := &datamining.BioData{
		PersonHash:      "bhash",
		BiodPubk:        "pubkey",
		PersonPubk:      "pubkey",
		CipherAddrRobot: "addr",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err := srv.LeadMining("hash", "addr", "fake sig", bd, datamining.CreateBioTransaction)
	assert.Nil(t, err)
	assert.Len(t, repo.Biometrics, 1)
	assert.Equal(t, "bhash", repo.Biometrics[0].PersonHash())
}

type mockSrvDatabase struct {
	Biometrics []*datamining.Biometric
	Keychains  []*datamining.Keychain
}

func (d *mockSrvDatabase) FindBiometric(bh string) (*datamining.Biometric, error) {
	return nil, nil
}

func (d *mockSrvDatabase) FindKeychain(addr string) (*datamining.Keychain, error) {
	return nil, nil
}

func (d *mockSrvDatabase) ListBiodPubKeys() ([]string, error) {
	return []string{}, nil
}

func (d *mockSrvDatabase) StoreKeychain(k *datamining.Keychain) error {
	d.Keychains = append(d.Keychains, k)
	return nil
}

func (d *mockSrvDatabase) StoreBiometric(b *datamining.Biometric) error {
	d.Biometrics = append(d.Biometrics, b)
	return nil
}

type mockSrvPoolDispatcher struct {
	Repo *mockSrvDatabase
}

func (r mockSrvPoolDispatcher) RequestLock(pool.Cluster, pool.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestUnlock(pool.Cluster, pool.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestValidations(sPool pool.Cluster, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockSrvPoolDispatcher) RequestStorage(sPool pool.Cluster, data interface{}, txType datamining.TransactionType) error {
	switch data.(type) {
	case *datamining.Keychain:
		r.Repo.StoreKeychain(data.(*datamining.Keychain))
	case *datamining.Biometric:
		r.Repo.StoreBiometric(data.(*datamining.Biometric))
	}

	return nil
}

type mockSrvPoolFinder struct{}

func (f mockSrvPoolFinder) FindLastValidationPool(addr string) (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindValidationPool() (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindStoragePool() (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "storage key",
			},
		},
	}, nil
}

type mockSrvNotifier struct{}

func (n mockSrvNotifier) NotifyTransactionStatus(txHash string, status TransactionStatus) error {
	return nil
}

type mockTxLocker struct{}

func (l mockTxLocker) Lock(txLock pool.TransactionLock) error {
	return nil
}

func (l mockTxLocker) Unlock(txLock pool.TransactionLock) error {
	return nil
}

func (l mockTxLocker) ContainsLock(txLock pool.TransactionLock) bool {
	return false
}

type mockSrvSigner struct{}

func (s mockSrvSigner) SignLock(txLock pool.TransactionLock, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSrvSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSrvSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSrvSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSrvSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return nil
}

type mockHasher struct{}

func (h mockHasher) HashData(data interface{}) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashTransactionData(data interface{}) (string, error) {
	return "hash", nil
}
