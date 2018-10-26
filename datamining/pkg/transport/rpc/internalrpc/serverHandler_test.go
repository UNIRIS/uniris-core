package internalrpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/mining/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
	"github.com/uniris/uniris-core/datamining/pkg/mining/transactions"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/uniris/uniris-core/datamining/pkg/listing"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
)

/*
Scenario: Retrieve a wallet from a bio hash
	Given a bio wallet store and a bio hash
	When i want to retrieve the wallet associated
	Then i can retrieve the wallet stored
*/
func TestGetWallet(t *testing.T) {

	repo := &databasemock{
		BioWallets: make([]*datamining.BioWallet, 0),
		Wallets:    make([]*datamining.Wallet, 0),
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	cipherAddr, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "addr")
	cipherBhash, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "hash")

	bdata := &datamining.BioData{
		BHash:           "hash",
		CipherAddrRobot: cipherAddr,
		CipherAESKey:    "encrypted_aes_key",
	}
	endors := datamining.NewEndorsement(
		time.Now(),
		"hello",
		&datamining.MasterValidation{},
		[]datamining.Validation{},
	)
	repo.StoreBioWallet(datamining.NewBioWallet(bdata, endors))

	wdata := &datamining.WalletData{
		WalletAddr:      "addr",
		CipherAddrRobot: cipherAddr,
	}
	repo.StoreWallet(datamining.NewWallet(wdata, endors, ""))

	list := listing.NewService(repo)
	errors := system.DataMininingErrors{}

	leading := mining.NewService(
		list,
		mockPoolFinder{},
		&mockPoolDispatcher{Repo: repo},
		mockTxLocker{},
		mockNotifier{},
		mockSigner{},
		mockHasher{},
		"robotPubKey",
		"robotPvKey",
	)

	h := NewInternalServerHandler(list, leading, hex.EncodeToString(pvKey), errors)
	res, err := h.GetWallet(context.TODO(), &api.WalletSearchRequest{
		EncryptedHashPerson: cipherBhash,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESkey)
}

/*
Scenario: Create a wallet
	Given some wallet encrypted
	When I want to store the wallet
	Then the wallet is stored and the updated wallet hash is updated
*/
func TestCreateWallet(t *testing.T) {
	repo := &databasemock{
		BioWallets: make([]*datamining.BioWallet, 0),
		Wallets:    make([]*datamining.Wallet, 0),
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	list := listing.NewService(repo)
	leading := mining.NewService(
		list,
		mockPoolFinder{},
		&mockPoolDispatcher{Repo: repo},
		mockTxLocker{},
		mockNotifier{},
		mockSigner{},
		mockHasher{},
		"robotPubKey",
		"robotPvKey",
	)
	errors := system.DataMininingErrors{}

	h := NewInternalServerHandler(list, leading, hex.EncodeToString(pvKey), errors)

	cipherAddr, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "encrypted_addr_robot")

	bioData := BioDataFromJSON{
		BiodPublicKey:       hex.EncodeToString(pbKey),
		EncryptedAddrPerson: "encrypted_person_addr",
		EncryptedAddrRobot:  cipherAddr,
		EncryptedAESKey:     "encrypted_aes_key",
		PersonHash:          "person_hash",
		PersonPublicKey:     hex.EncodeToString(pbKey),
	}

	walletData := WalletDataFromJSON{
		BiodPublicKey:      hex.EncodeToString(pbKey),
		EncryptedAddrRobot: cipherAddr,
		EncryptedWallet:    "encrypted_wallet",
		PersonPublicKey:    hex.EncodeToString(pbKey),
	}

	bBData, _ := json.Marshal(bioData)
	cipherBio, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(bBData))
	sigBio, _ := crypto.Sign(hex.EncodeToString(pvKey), string(bBData))

	bWData, _ := json.Marshal(walletData)
	cipherWallet, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(bWData))
	sigWal, _ := crypto.Sign(hex.EncodeToString(pvKey), string(bWData))

	req := &api.WalletCreationRequest{
		EncryptedBioData:    cipherBio,
		EncryptedWalletData: cipherWallet,
		SignatureBioData: &api.Signature{
			Biod:   sigBio,
			Person: sigBio,
		},
		SignatureWalletData: &api.Signature{
			Biod:   sigWal,
			Person: sigWal,
		},
	}

	res, err := h.CreateWallet(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotNil(t, res.BioTransactionHash)
	assert.NotNil(t, res.DataTransactionHash)
}

type databasemock struct {
	BioWallets []*datamining.BioWallet
	Wallets    []*datamining.Wallet
}

func (d *databasemock) FindBioWallet(bh string) (*datamining.BioWallet, error) {
	for _, bw := range d.BioWallets {
		if bw.Bhash() == bh {
			return bw, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindWallet(addr string) (*datamining.Wallet, error) {
	for _, w := range d.Wallets {
		if w.WalletAddr() == addr {
			return w, nil
		}
	}
	return nil, nil
}

func (d *databasemock) ListBiodPubKeys() ([]string, error) {
	return []string{}, nil
}

func (d *databasemock) StoreWallet(w *datamining.Wallet) error {
	d.Wallets = append(d.Wallets, w)
	return nil
}

func (d *databasemock) StoreBioWallet(bw *datamining.BioWallet) error {
	d.BioWallets = append(d.BioWallets, bw)
	return nil
}

type mockTechRepo struct{}

func (r mockTechRepo) ListBiodPubKeys() ([]string, error) {
	return []string{"key1", "key2", "key3"}, nil
}

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

func (s mockSigner) CheckTransactionSignature(pubk string, tx string, der string) error {
	return nil
}

func (s mockSigner) SignLock(lock.TransactionLock, string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignMasterValidation(v mining.Validation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignValidation(v mining.Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockHasher struct{}

func (h mockHasher) HashWallet(w *datamining.Wallet) (string, error) {
	return "hashed wallet", nil
}

type mockPoolDispatcher struct {
	Repo *databasemock
}

func (r mockPoolDispatcher) RequestLock(pool.PeerCluster, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestUnlock(pool.PeerCluster, lock.TransactionLock, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestValidations(sPool pool.PeerCluster, data interface{}, txType transactions.Type) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolDispatcher) RequestStorage(sPool pool.PeerCluster, data interface{}, txType transactions.Type) error {
	switch data.(type) {
	case *datamining.Wallet:
		r.Repo.StoreWallet(data.(*datamining.Wallet))
	case *datamining.BioWallet:
		r.Repo.StoreBioWallet(data.(*datamining.BioWallet))
	}

	return nil
}

type mockPoolFinder struct{}

func (f mockPoolFinder) FindLastValidationPool(addr string) (pool.PeerCluster, error) {
	return pool.PeerCluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockPoolFinder) FindValidationPool() (pool.PeerCluster, error) {
	return pool.PeerCluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockPoolFinder) FindStoragePool() (pool.PeerCluster, error) {
	return pool.PeerCluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "storage key",
			},
		},
	}, nil
}

type mockNotifier struct{}

func (n mockNotifier) NotifyTransactionStatus(txHash string, status mining.TransactionStatus) error {
	return nil
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
