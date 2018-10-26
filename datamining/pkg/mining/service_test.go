package mining

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/transactions"
	"github.com/uniris/uniris-core/datamining/pkg/mining/validations"

	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Lead wallet creation
	Given a wallet data
	When I want to lead the wallet creation and validation
	Then I perform POW and process validations and create the wallet
*/
func TestLeadWallet(t *testing.T) {
	repo := &mockSrvDatabase{
		BioWallets: make([]*datamining.BioWallet, 0),
		Wallets:    make([]*datamining.Wallet, 0),
	}
	srv := NewService(
		listing.NewService(repo),
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockTxLocker{},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvHasher{},
		"robotPubKey",
		"robotPvKey")

	w := &datamining.WalletData{
		BiodPubk:        "pubkey",
		EmPubk:          "pubkey",
		CipherAddrRobot: "addr",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err := srv.Lead("hash", "addr", "fake sig", w, transactions.CreateWallet)
	assert.Nil(t, err)
	assert.Len(t, repo.Wallets, 1)
	assert.Equal(t, "hash", repo.Wallets[0].Endorsement().TransactionHash())
	assert.Equal(t, "robotPubKey", repo.Wallets[0].Endorsement().MasterValidation().ProofOfWorkRobotKey())
	assert.Len(t, repo.Wallets[0].Endorsement().Validations(), 1)
	assert.Equal(t, "pubkey", repo.Wallets[0].Endorsement().Validations()[0].PublicKey())
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
		listing.NewService(repo),
		mockSrvPoolFinder{},
		&mockSrvPoolDispatcher{Repo: repo},
		mockTxLocker{},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvHasher{},
		"robotPubKey",
		"robotPvKey")

	bd := &datamining.BioData{
		BHash:           "bhash",
		BiodPubk:        "pubkey",
		EmPubk:          "pubkey",
		CipherAddrRobot: "addr",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err := srv.Lead("hash", "addr", "fake sig", bd, transactions.CreateBio)
	assert.Nil(t, err)
	assert.Len(t, repo.BioWallets, 1)
	assert.Equal(t, "bhash", repo.BioWallets[0].Bhash())
}

/*
Scenario: Validate wallet data
	Given a wallet data
	When I want validate it
	Then I get a validated transaction
*/
func TestValidateWallet(t *testing.T) {

	srv := service{
		robotKey: "key",
		sig:      mockSrvSigner{},
		checks: map[transactions.Type][]validations.Handler{
			transactions.CreateWallet: []validations.Handler{
				validations.NewSignatureValidation(mockSrvSigner{}),
			},
		},
	}
	w := &datamining.WalletData{
		BiodPubk: "pubKey",
		EmPubk:   "pubKey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	v, err := srv.Validate(w, transactions.CreateWallet)
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
		sig:      mockSrvSigner{},
		checks: map[transactions.Type][]validations.Handler{
			transactions.CreateWallet: []validations.Handler{
				validations.NewSignatureValidation(mockBadSigner{}),
			},
		},
	}

	w := &datamining.WalletData{
		BiodPubk: "pubKey",
		EmPubk:   "pubKey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	v, err := srv.Validate(w, transactions.CreateWallet)
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
		sig:      mockSrvSigner{},
		checks: map[transactions.Type][]validations.Handler{
			transactions.CreateWallet: []validations.Handler{
				validations.NewSignatureValidation(mockSrvSigner{}),
			},
		},
	}

	b := &datamining.BioData{
		BiodPubk: "pubkey",
		EmPubk:   "pubkey",
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	v, err := srv.Validate(b, transactions.CreateBio)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, v.Status())
	assert.Equal(t, "signature", v.Signature())
	assert.Equal(t, "key", v.PublicKey())
}

type mockBadSigner struct{}

func (s mockBadSigner) SignValidation(v Validation, pvKey string) (string, error) {
	return "signature", nil
}

func (s mockBadSigner) CheckSignature(pubKey string, data interface{}, sig string) error {
	return validations.ErrInvalidSignature
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

func (s mockSrvSigner) SignLock(txLock lock.TransactionLock, pvKey string) (string, error) {
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

type mockSrvDatabase struct {
	BioWallets []*datamining.BioWallet
	Wallets    []*datamining.Wallet
}

func (d *mockSrvDatabase) FindBioWallet(bh string) (*datamining.BioWallet, error) {
	return nil, nil
}

func (d *mockSrvDatabase) FindWallet(addr string) (*datamining.Wallet, error) {
	return nil, nil
}

func (d *mockSrvDatabase) ListBiodPubKeys() ([]string, error) {
	return []string{}, nil
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

func (r mockSrvPoolDispatcher) RequestLock(pool.PeerCluster, lock.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestUnlock(pool.PeerCluster, lock.TransactionLock, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestValidations(sPool pool.PeerCluster, data interface{}, txType transactions.Type) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockSrvPoolDispatcher) RequestStorage(sPool pool.PeerCluster, data interface{}, txType transactions.Type) error {
	switch data.(type) {
	case *datamining.Wallet:
		r.Repo.StoreWallet(data.(*datamining.Wallet))
	case *datamining.BioWallet:
		r.Repo.StoreBioWallet(data.(*datamining.BioWallet))
	}

	return nil
}

type mockSrvPoolFinder struct{}

func (f mockSrvPoolFinder) FindLastValidationPool(addr string) (pool.PeerCluster, error) {
	return pool.PeerCluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindValidationPool() (pool.PeerCluster, error) {
	return pool.PeerCluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockSrvPoolFinder) FindStoragePool() (pool.PeerCluster, error) {
	return pool.PeerCluster{
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

type mockSrvHasher struct{}

func (h mockSrvHasher) HashWallet(*datamining.Wallet) (string, error) {
	return "hash", nil
}

type mockTxLocker struct{}

func (l mockTxLocker) Lock(txLock lock.TransactionLock) error {
	return nil
}

func (l mockTxLocker) Unlock(txLock lock.TransactionLock) error {
	return nil
}

func (l mockTxLocker) ContainsLock(txLock lock.TransactionLock) bool {
	return false
}
