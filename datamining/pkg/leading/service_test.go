package leading

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/validating"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Lead bio validation
	Given a wallet data
	When I want to lead the data validation
	Then I perform POW and returns validated transaction
*/
func TestLeadWallet(t *testing.T) {
	repo := &mockSrvDatabase{
		BioWallets: make([]*datamining.BioWallet, 0),
		Wallets:    make([]*datamining.Wallet, 0),
	}
	srv := NewService(
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvHasher{},
		mockSrvTechRepo{},
		"robotPubKey",
		"robotPvKey")

	w := &datamining.WalletData{
		BiodPubk: "pubkey",
		EmPubk:   "pubkey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err := srv.LeadWalletTransaction(w, "hash")
	assert.Nil(t, err)
	assert.Len(t, repo.Wallets, 1)
	assert.Equal(t, "hash", repo.Wallets[0].Endorsement().TransactionHash())
	assert.Equal(t, "robotPubKey", repo.Wallets[0].Endorsement().MasterValidation().ProofOfWorkRobotKey())
	assert.Len(t, repo.Wallets[0].Endorsement().Validations(), 1)
	assert.Equal(t, "validator key", repo.Wallets[0].Endorsement().Validations()[0].PublicKey())
}

/*
Scenario: Lead bio validation
	Given a bio data
	When I want to lead the data validation
	Then I perform POW and returns validated transaction
*/
func TestLeadBio(t *testing.T) {
	repo := &mockSrvDatabase{
		BioWallets: make([]*datamining.BioWallet, 0),
		Wallets:    make([]*datamining.Wallet, 0),
	}
	srv := NewService(
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvHasher{},
		mockSrvTechRepo{},
		"robotPubKey",
		"robotPvKey")

	bd := &datamining.BioData{
		BHash:    "bhash",
		BiodPubk: "pubkey",
		EmPubk:   "pubkey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err := srv.LeadBioTransaction(bd, "hash")
	assert.Nil(t, err)
	assert.Len(t, repo.BioWallets, 1)
	assert.Equal(t, "bhash", repo.BioWallets[0].Bhash())
}

type mockSigner struct{}

func (s mockSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return nil
}

type mockSrvTechRepo struct{}

func (r mockSrvTechRepo) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

type mockSrvSigner struct{}

func (s mockSrvSigner) SignLock(txLock validating.TransactionLock, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockSrvSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSrvSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockSrvDatabase struct {
	BioWallets []*datamining.BioWallet
	Wallets    []*datamining.Wallet
}

func (d *mockSrvDatabase) StoreWallet(w *datamining.Wallet) error {
	d.Wallets = append(d.Wallets, w)
	return nil
}

func (d *mockSrvDatabase) StoreBioWallet(bw *datamining.BioWallet) error {
	d.BioWallets = append(d.BioWallets, bw)
	return nil
}

type mockSrvPoolDispatcher struct {
	Repo *mockSrvDatabase
}

func (r mockSrvPoolDispatcher) RequestLock(Pool, validating.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestUnlock(Pool, validating.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestWalletValidation(Pool, *datamining.WalletData) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(datamining.ValidationOK, time.Now(), "validator key", "sig"),
	}, nil
}
func (r mockSrvPoolDispatcher) RequestBioValidation(Pool, *datamining.BioData) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(datamining.ValidationOK, time.Now(), "validator key", "sig"),
	}, nil
}

func (r *mockSrvPoolDispatcher) RequestWalletStorage(p Pool, w *datamining.Wallet) error {
	return r.Repo.StoreWallet(w)
}

func (r *mockSrvPoolDispatcher) RequestBioStorage(p Pool, w *datamining.BioWallet) error {
	return r.Repo.StoreBioWallet(w)
}

func (r mockSrvPoolDispatcher) RequestLastWallet(Pool, string) ([]*datamining.Wallet, error) {
	return []*datamining.Wallet{
		datamining.NewWallet(
			&datamining.WalletData{
				BiodPubk:        "pub",
				EmPubk:          "pub",
				CipherAddrRobot: "addr",
				CipherWallet:    "addr",
			},
			datamining.NewEndorsement(time.Now(), "hash", nil, nil),
			"hash"),
	}, nil
}

type mockSrvPoolFinder struct{}

func (f mockSrvPoolFinder) FindLastValidationPool(addr string) (Pool, error) {
	return Pool{
		Peers: []Peer{
			Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindValidationPool() (Pool, error) {
	return Pool{
		Peers: []Peer{
			Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindStoragePool() (Pool, error) {
	return Pool{
		Peers: []Peer{
			Peer{
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

type mockSrvHasher struct{}

func (h mockSrvHasher) HashWallet(*datamining.Wallet) (string, error) {
	return "hash", nil
}
