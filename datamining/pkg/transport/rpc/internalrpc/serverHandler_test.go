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

	"github.com/uniris/uniris-core/datamining/pkg/mining/master"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"

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
		Biometrics: make([]*datamining.Biometric, 0),
		Keychains:  make([]*datamining.Keychain, 0),
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	cipherAddr, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "addr")
	cipherBhash, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "hash")

	bdata := &datamining.BioData{
		PersonHash:      "hash",
		CipherAddrRobot: cipherAddr,
		CipherAESKey:    "encrypted_aes_key",
	}
	endors := datamining.NewEndorsement(
		time.Now(),
		"hello",
		&datamining.MasterValidation{},
		[]datamining.Validation{},
	)
	repo.StoreBiometric(datamining.NewBiometric(bdata, endors))

	kdata := &datamining.KeyChainData{
		WalletAddr:      "addr",
		CipherAddrRobot: cipherAddr,
	}
	repo.StoreKeychain(datamining.NewKeychain(kdata, endors, ""))

	list := listing.NewService(repo)
	errors := system.DataMininingErrors{}

	master := master.NewService(
		mockPoolFinder{},
		&mockPoolDispatcher{Repo: repo},
		mockNotifier{},
		mockSigner{},
		mockHasher{},
		list,
		"robotPubKey",
		"robotPvKey",
	)

	h := NewInternalServerHandler(list, master, hex.EncodeToString(pvKey), errors)
	res, err := h.GetAccount(context.TODO(), &api.AccountSearchRequest{
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
		Biometrics: make([]*datamining.Biometric, 0),
		Keychains:  make([]*datamining.Keychain, 0),
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	list := listing.NewService(repo)
	leading := master.NewService(
		mockPoolFinder{},
		&mockPoolDispatcher{Repo: repo},
		mockNotifier{},
		mockSigner{},
		mockHasher{},
		list,
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
		PersonHash:          "person hash",
		PersonPublicKey:     hex.EncodeToString(pbKey),
	}

	keychainData := KeychainDataFromJSON{
		BiodPublicKey:      hex.EncodeToString(pbKey),
		EncryptedAddrRobot: cipherAddr,
		EncryptedWallet:    "encrypted_wallet",
		PersonPublicKey:    hex.EncodeToString(pbKey),
	}

	bioBytes, _ := json.Marshal(bioData)
	cipherBio, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(bioBytes))
	sigBio, _ := crypto.Sign(hex.EncodeToString(pvKey), string(bioBytes))

	keychainBytes, _ := json.Marshal(keychainData)
	cipherWallet, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(keychainBytes))
	sigKeychain, _ := crypto.Sign(hex.EncodeToString(pvKey), string(keychainBytes))

	req := &api.AccountCreationRequest{
		EncryptedBioData:      cipherBio,
		EncryptedKeychainData: cipherWallet,
		SignatureBioData: &api.Signature{
			Biod:   sigBio,
			Person: sigBio,
		},
		SignatureKeychainData: &api.Signature{
			Biod:   sigKeychain,
			Person: sigKeychain,
		},
	}

	res, err := h.CreateAccount(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotNil(t, res.BioTransactionHash)
	assert.NotNil(t, res.KeychainTransactionHash)
}

type databasemock struct {
	Biometrics []*datamining.Biometric
	Keychains  []*datamining.Keychain
}

func (d *databasemock) FindBiometric(hash string) (*datamining.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == hash {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindKeychain(addr string) (*datamining.Keychain, error) {
	for _, w := range d.Keychains {
		if w.WalletAddr() == addr {
			return w, nil
		}
	}
	return nil, nil
}

func (d *databasemock) ListBiodPubKeys() ([]string, error) {
	return []string{}, nil
}

func (d *databasemock) StoreKeychain(w *datamining.Keychain) error {
	d.Keychains = append(d.Keychains, w)
	return nil
}

func (d *databasemock) StoreBiometric(bw *datamining.Biometric) error {
	d.Biometrics = append(d.Biometrics, bw)
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

func (s mockSigner) SignLock(pool.TransactionLock, string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignMasterValidation(v master.Validation, pvKey string) (string, error) {
	return "sig", nil
}

func (s mockSigner) SignValidation(v slave.Validation, pvKey string) (string, error) {
	return "sig", nil
}

type mockHasher struct{}

func (h mockHasher) HashTransactionData(data interface{}) (string, error) {
	return crypto.NewHasher().HashTransactionData(data)
}

type mockPoolDispatcher struct {
	Repo *databasemock
}

func (r mockPoolDispatcher) RequestLock(pool.Cluster, pool.TransactionLock, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestUnlock(pool.Cluster, pool.TransactionLock, string) error {
	return nil
}

func (r mockPoolDispatcher) RequestValidations(sPool pool.Cluster, data interface{}, txType datamining.TransactionType) ([]datamining.Validation, error) {
	return []datamining.Validation{
		datamining.NewValidation(
			datamining.ValidationOK,
			time.Now(),
			"pubkey",
			"fake sig",
		)}, nil
}

func (r mockPoolDispatcher) RequestStorage(sPool pool.Cluster, data interface{}, txType datamining.TransactionType) error {
	switch data.(type) {
	case *datamining.Keychain:
		r.Repo.StoreKeychain(data.(*datamining.Keychain))
	case *datamining.Biometric:
		r.Repo.StoreBiometric(data.(*datamining.Biometric))
	}

	return nil
}

type mockPoolFinder struct{}

func (f mockPoolFinder) FindLastValidationPool(addr string) (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockPoolFinder) FindValidationPool() (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "validator key",
			},
		},
	}, nil
}

func (f mockPoolFinder) FindStoragePool() (pool.Cluster, error) {
	return pool.Cluster{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "storage key",
			},
		},
	}, nil
}

type mockNotifier struct{}

func (n mockNotifier) NotifyTransactionStatus(txHash string, status master.TransactionStatus) error {
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
