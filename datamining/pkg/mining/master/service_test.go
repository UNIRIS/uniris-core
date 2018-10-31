package master

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/locking"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
	"github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

/*
Scenario: Lead wallet creation
	Given a wallet data
	When I want to lead the wallet creation and validation
	Then I perform POW and process validations and create the wallet
*/
func TestLeadWallet(t *testing.T) {
	repo := mock.NewDatabase()
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
	repo := mock.NewDatabase()
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

type mockSrvPoolDispatcher struct {
	Repo *mock.Databasemock
}

func (r mockSrvPoolDispatcher) RequestLock(pool.PeerGroup, locking.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestUnlock(pool.PeerGroup, locking.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestValidations(sPool pool.PeerGroup, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockSrvPoolDispatcher) RequestStorage(sPool pool.PeerGroup, data interface{}, txType datamining.TransactionType) error {
	switch data.(type) {
	case *datamining.Keychain:
		r.Repo.StoreKeychain(data.(*datamining.Keychain))
	case *datamining.Biometric:
		r.Repo.StoreBiometric(data.(*datamining.Biometric))
	}

	return nil
}

type mockSrvPoolFinder struct{}

func (f mockSrvPoolFinder) FindLastValidationPool(addr string) (pool.PeerGroup, error) {
	return pool.PeerGroup{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindValidationPool() (pool.PeerGroup, error) {
	return pool.PeerGroup{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindStoragePool() (pool.PeerGroup, error) {
	return pool.PeerGroup{
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

func (l mockTxLocker) Lock(txLock locking.TransactionLock) error {
	return nil
}

func (l mockTxLocker) Unlock(txLock locking.TransactionLock) error {
	return nil
}

func (l mockTxLocker) ContainsLock(txLock locking.TransactionLock) bool {
	return false
}

type mockSrvSigner struct{}

func (s mockSrvSigner) SignLock(txLock locking.TransactionLock, pvKey string) (string, error) {
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
