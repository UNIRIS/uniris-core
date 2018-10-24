package leading

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Lead bio validation
	Given a wallet data
	When I want to lead the data validation
	Then I perform POW and returns validated transaction
*/
func TestLeadWalletValidation(t *testing.T) {
	srv := NewService(
		mockSrvPoolFinder{},
		mockSrvPoolDispatcher{},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvTechRepo{},
		"robotPubKey",
		"robotPvKey")

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	w := &datamining.WalletData{
		BiodPubk: hex.EncodeToString(pbKey),
		EmPubk:   hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	endorsement, err := srv.ValidateWallet(w, "hash", nil)
	assert.Nil(t, err)
	assert.Equal(t, "hash", endorsement.TransactionHash())
	assert.Equal(t, "robotPubKey", endorsement.MasterValidation().ProofOfWorkRobotKey())
	assert.Len(t, endorsement.Validations(), 1)
	assert.Equal(t, "validator key", endorsement.Validations()[0].PublicKey())
}

/*
Scenario: Lead bio validation
	Given a bio data
	When I want to lead the data validation
	Then I perform POW and returns validated transaction
*/
func TestLeadBioValidation(t *testing.T) {
	srv := NewService(
		mockSrvPoolFinder{},
		mockSrvPoolDispatcher{},
		mockSrvNotifier{},
		mockSrvSigner{},
		mockSrvTechRepo{},
		"robotPubKey",
		"robotPvKey")

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	bd := &datamining.BioData{
		BiodPubk: hex.EncodeToString(pbKey),
		EmPubk:   hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	endorsement, err := srv.ValidateBio(bd, "hash", nil)
	assert.Nil(t, err)
	assert.Equal(t, "hash", endorsement.TransactionHash())
	assert.Equal(t, "robotPubKey", endorsement.MasterValidation().ProofOfWorkRobotKey())
	assert.Len(t, endorsement.Validations(), 1)
	assert.Equal(t, "validator key", endorsement.Validations()[0].PublicKey())
}

type mockSrvTechRepo struct{}

func (r mockSrvTechRepo) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

type mockSrvSigner struct{}

func (s mockSrvSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSrvSigner) SignMasterValidation(v Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockSrvPoolDispatcher struct{}

func (r mockSrvPoolDispatcher) RequestLock(Pool, string) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestUnlock(Pool, string) error {
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

func (r mockSrvPoolDispatcher) RequestWalletStorage(Pool, *datamining.Wallet) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestBioStorage(Pool, *datamining.BioWallet) error {
	return nil
}

func (r mockSrvPoolDispatcher) RequestLastTx(Pool, string) (string, *datamining.MasterValidation, error) {
	return "last", &datamining.MasterValidation{}, nil
}

type mockSrvPoolFinder struct{}

func (f mockSrvPoolFinder) FindValidationPool() (Pool, error) {
	return Pool{
		Peers: []Peer{
			Peer{
				IP:        net.ParseIP("127.0.0.1"),
				Port:      4000,
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
				Port:      4000,
				PublicKey: "storage key",
			},
		},
	}, nil
}

type mockSrvNotifier struct{}

func (n mockSrvNotifier) NotifyTransactionStatus(txHash string, status TransactionStatus) error {
	return nil
}
