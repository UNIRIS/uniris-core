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
	"sync"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"github.com/uniris/uniris-core/datamining/pkg/validating"

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

	repo := &databasemock{}

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
	valid := validating.NewService(
		mockTechRepo{},
		mockTransactionFetcher{},
		mockPoolDispatcher{},
		mockPoolFinder{},
		mockHasher{},
		mockSigner{},
		mockNotifier{},
		"robotPubKey",
		"robotPvKey",
	)
	d := &mockPoolDispatcher{}
	adding := adding.NewService(valid, mockPoolFinder{}, d)
	errors := system.DataMininingErrors{}

	h := NewInternalServerHandler(list, adding, hex.EncodeToString(pvKey), errors, nil)
	res, err := h.GetWallet(context.TODO(), &api.WalletSearchRequest{
		EncryptedHashPerson: cipherBhash,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESkey)
}

/*
Scenario: Store a wallet
	Given some wallet encrypted
	When I want to store the wallet
	Then the wallet is stored and the updated wallet hash is updated
*/
func TestStoreWallet(t *testing.T) {
	repo := &databasemock{}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	list := listing.NewService(repo)
	valid := validating.NewService(
		mockTechRepo{},
		mockTransactionFetcher{},
		mockPoolDispatcher{},
		mockPoolFinder{},
		mockHasher{},
		mockSigner{},
		mockNotifier{},
		"robotPubKey",
		"robotPvKey",
	)
	d := &mockPoolDispatcher{
		Repo: repo,
	}
	adding := adding.NewService(valid, mockPoolFinder{}, d)
	errors := system.DataMininingErrors{}

	storeChan := make(chan string)
	h := NewInternalServerHandler(list, adding, hex.EncodeToString(pvKey), errors, storeChan)

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

	req := &api.WalletStorageRequest{
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

	res, err := h.StoreWallet(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotNil(t, res.BioTransactionHash)
	assert.NotNil(t, res.DataTransactionHash)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		for range storeChan {
			wg.Done()
		}
		close(storeChan)
	}()

	wg.Wait()

	assert.Len(t, repo.Wallets, 1)
	assert.Len(t, repo.BioWallets, 1)

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

type mockTransactionFetcher struct{}

func (f mockTransactionFetcher) FindPreviousTransactionMiners(addr string) ([]string, error) {
	return []string{
		"minerKey1",
		"minerKey2",
	}, nil
}
func (f mockTransactionFetcher) FindLastWalletTx(addr string) (string, error) {
	return "last tx", nil
}

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

func (s mockSigner) SignValidation(v validating.Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockHasher struct{}

func (h mockHasher) HashMasterValidation(v *datamining.MasterValidation) (string, error) {
	return "hashed validation", nil
}

type mockPoolDispatcher struct {
	Repo *databasemock
}

func (r mockPoolDispatcher) RequestLastTx(pool adding.Pool, txHash string) (oldTxHash string, validation *datamining.MasterValidation, err error) {
	return "old tx hash", datamining.NewMasterValidation([]string{"key", "key2"}, "key", datamining.Validation{}), nil
}

func (r *mockPoolDispatcher) RequestWalletStorage(p adding.Pool, w *datamining.Wallet) error {
	return r.Repo.StoreWallet(w)
}

func (r *mockPoolDispatcher) RequestBioStorage(p adding.Pool, w *datamining.BioWallet) error {
	return r.Repo.StoreBioWallet(w)
}

func (r mockPoolDispatcher) RequestLock(validating.Pool, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestUnlock(validating.Pool, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestWalletValidation(validating.Pool, *datamining.WalletData) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(datamining.ValidationOK, time.Now(), "validator key", "sig"),
	}, nil
}
func (r mockPoolDispatcher) RequestBioValidation(validating.Pool, *datamining.BioData) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(datamining.ValidationOK, time.Now(), "validator key", "sig"),
	}, nil
}

type mockPoolFinder struct{}

func (f mockPoolFinder) FindValidationPool() (validating.Pool, error) {
	return validating.Pool{
		Peers: []validating.Peer{
			validating.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				Port:      4000,
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockPoolFinder) FindStoragePool() (adding.Pool, error) {
	return adding.Pool{
		Peers: []adding.Peer{
			adding.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				Port:      4000,
				PublicKey: "validator key",
			},
		},
	}, nil
}

type mockNotifier struct{}

func (n mockNotifier) NotifyTransactionStatus(txHash string, status validating.TransactionStatus) error {
	return nil
}
