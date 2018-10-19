package internalrpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
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
	valid := validating.NewService(mockSigner{}, mockValiationRequester{})
	adding := adding.NewService(repo, valid)
	errors := system.DataMininingErrors{}

	h := NewInternalServerHandler(list, adding, hex.EncodeToString(pvKey), errors)
	res, err := h.GetWallet(context.TODO(), &api.WalletRequest{
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
	valid := validating.NewService(mockSigner{}, mockValiationRequester{})
	adding := adding.NewService(repo, valid)
	errors := system.DataMininingErrors{}

	h := NewInternalServerHandler(list, adding, hex.EncodeToString(pvKey), errors)

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

	req := &api.Wallet{
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

	assert.Len(t, repo.Wallets, 1)
	assert.Len(t, repo.BioWallets, 1)

	assert.NotNil(t, res.TransactionHash)
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

type mockSigner struct{}

func (mockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

type mockValiationRequester struct{}

func (v mockValiationRequester) RequestWalletValidation(validating.Peer, *datamining.WalletData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}

func (v mockValiationRequester) RequestBioValidation(validating.Peer, *datamining.BioData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}
