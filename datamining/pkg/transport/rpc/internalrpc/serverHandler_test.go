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

	cipher, _ := crypto.Encrypt([]byte(hex.EncodeToString(pbKey)), []byte("biohash"))

	data := datamining.BioData{
		BHash:           []byte("hash"),
		CipherAddrRobot: datamining.WalletAddr(cipher),
		CipherAESKey:    []byte("encrypted_aes_key"),
	}
	endors := datamining.NewEndorsement(
		datamining.Timestamp(time.Now()),
		[]byte("hello"),
		datamining.MasterValidation{},
		[]datamining.Validation{},
	)
	repo.StoreBioWallet(datamining.NewBioWallet(data, endors))

	list := listing.NewService(repo)
	valid := validating.NewService(mockSigner{}, mockValiationRequester{})
	adding := adding.NewService(repo, valid)

	h := NewInternalServerHandler(list, adding, []byte(hex.EncodeToString(pbKey)), []byte(hex.EncodeToString(pvKey)))
	res, err := h.GetWallet(context.TODO(), &api.WalletRequest{
		EncryptedHashPerson: []byte("hash"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, []byte("encrypted_aes_key"), res.EncryptedAESkey)
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

	h := NewInternalServerHandler(list, adding, []byte(hex.EncodeToString(pbKey)), []byte(hex.EncodeToString(pvKey)))

	cipherAddr, _ := crypto.Encrypt([]byte(hex.EncodeToString(pbKey)), []byte("encrypted_addr_robot"))

	bioData := BioDataFromJSON{
		BiodPublicKey:       hex.EncodeToString(pbKey),
		EncryptedAddrPerson: "encrypted_person_addr",
		EncryptedAddrRobot:  string(cipherAddr),
		EncryptedAESKey:     "encrypted_aes_key",
		PersonHash:          "person_hash",
		PersonPublicKey:     hex.EncodeToString(pbKey),
	}

	walletData := WalletDataFromJSON{
		BiodPublicKey:      hex.EncodeToString(pbKey),
		EncryptedAddrRobot: "encrypted_addr_robot",
		EncryptedWallet:    "encrypted_wallet",
		PersonPublicKey:    hex.EncodeToString(pbKey),
	}

	bBData, _ := json.Marshal(bioData)
	cipherBio, _ := crypto.Encrypt([]byte(hex.EncodeToString(pbKey)), bBData)
	sigBio, _ := crypto.Sign([]byte(hex.EncodeToString(pvKey)), bBData)

	bWData, _ := json.Marshal(walletData)
	cipherWallet, _ := crypto.Encrypt([]byte(hex.EncodeToString(pbKey)), bWData)
	sigWal, _ := crypto.Sign([]byte(hex.EncodeToString(pvKey)), bWData)

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

	assert.NotNil(t, res.HashUpdatedWallet)
}

type databasemock struct {
	BioWallets []datamining.BioWallet
	Wallets    []datamining.Wallet
}

func (d *databasemock) FindBioWallet(bh datamining.BioHash) (b datamining.BioWallet, err error) {
	for _, bw := range d.BioWallets {
		if string(bw.Bhash()) == string(bh) {
			b = bw
			break
		}
	}
	return
}

func (d *databasemock) FindWallet(addr datamining.WalletAddr) (b datamining.Wallet, err error) {
	for _, b := range d.Wallets {
		if string(b.WalletAddr()) == string(addr) {
			return b, nil
		}
	}
	return
}

func (d *databasemock) StoreWallet(w datamining.Wallet) error {
	d.Wallets = append(d.Wallets, w)
	return nil
}

func (d *databasemock) StoreBioWallet(bw datamining.BioWallet) error {
	d.BioWallets = append(d.BioWallets, bw)
	return nil
}

type mockSigner struct{}

func (mockSigner) CheckSignature(pubk []byte, data interface{}, der []byte) error {
	return nil
}

type mockValiationRequester struct{}

func (v mockValiationRequester) RequestWalletValidation(validating.Peer, datamining.WalletData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}

func (v mockValiationRequester) RequestBioValidation(validating.Peer, datamining.BioData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}
