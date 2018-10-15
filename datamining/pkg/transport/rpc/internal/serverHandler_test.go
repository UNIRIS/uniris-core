package internalrpc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/uniris/ecies/pkg"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
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

	pu, err := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, []byte("biohash"), nil, nil)

	data := datamining.BioData{
		BHash:           []byte("hash"),
		CipherAddrRobot: datamining.WalletAddr(hex.EncodeToString(cipher)),
		CipherAESKey:    []byte("encrypted_aes_key"),
	}
	endors := datamining.NewEndorsement(
		datamining.Timestamp(time.Now()),
		[]byte("hello"),
		datamining.MasterValidation{},
		[]datamining.Validation{},
	)
	repo.StoreBioWallet(datamining.NewBioWallet(data, endors))

	h := NewInternalServerHandler(repo, repo, []byte(hex.EncodeToString(pvKey)))
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

	h := NewInternalServerHandler(repo, repo, []byte(hex.EncodeToString(pvKey)))

	bioData := BioDataFromJSON{
		BiodPublicKey:       "pub key",
		EncryptedAddrPerson: "encrypted_addr_person",
		EncryptedAddrRobot:  "encrypted_addr_robot",
		EncryptedAESKey:     "encrypted_aes_key",
		PersonHash:          "person_hash",
		PersonPublicKey:     "pub key",
		Sigs: Signatures{
			Biod:   "biod_sig",
			Person: "biod_pers",
		},
	}

	walletData := WalletDataFromJSON{
		BiodPublicKey:      "pub key",
		EncryptedAddrRobot: "encrypted_addr_robot",
		EncryptedWallet:    "encrypted_wallet",
		PersonPublicKey:    "pub key",
		Sigs: Signatures{
			Biod:   "biod_sig",
			Person: "biod_pers",
		},
	}

	pu, err := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	bBData, _ := json.Marshal(bioData)
	cipherBio, _ := ecies.Encrypt(rand.Reader, robotEciesKey, []byte(hex.EncodeToString(bBData)), nil, nil)

	bWData, _ := json.Marshal(walletData)
	cipherWallet, _ := ecies.Encrypt(rand.Reader, robotEciesKey, []byte(hex.EncodeToString(bWData)), nil, nil)

	req := &api.Wallet{
		EncryptedBioData:    []byte(hex.EncodeToString(cipherBio)),
		EncryptedWalletData: []byte(hex.EncodeToString(cipherWallet)),
		SignatureBioData: &api.Signature{
			Biod:   []byte(bioData.Sigs.Biod),
			Person: []byte(bioData.Sigs.Person),
		},
		SignatureWalletData: &api.Signature{
			Biod:   []byte(walletData.Sigs.Biod),
			Person: []byte(walletData.Sigs.Person),
		},
	}

	res, err := h.StoreWallet(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Len(t, repo.Wallets, 1)
	assert.Len(t, repo.BioWallets, 1)

	// TODO when mining: assert.NotNil(t, res.HashUpdatedWallet)
}

type databasemock struct {
	BioWallets []datamining.BioWallet
	Wallets    []datamining.Wallet
}

func (d *databasemock) FindBioWallet(bh datamining.BioHash) (b datamining.BioWallet, err error) {
	for _, b := range d.BioWallets {
		if string(b.Bhash()) == string(bh) {
			return b, nil
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
