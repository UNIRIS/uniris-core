package internalrpc

import (
	"testing"

	"github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Build a wallet result
	Given a wallet and bio data
	When I want build a wallet result
	Then I get a new object with encrypted aes, wallet and addr given
*/
func TestBuildWalletResult(t *testing.T) {

	bioData := datamining.BioData{
		CipherAESKey: []byte("encrypted aes key"),
	}
	bioWallet := datamining.NewBioWallet(bioData, datamining.Endorsement{})

	walData := datamining.WalletData{
		CipherWallet: []byte("encrypted wallet"),
		WalletAddr:   []byte("encrypted wallet addr"),
	}

	wallet := datamining.NewWallet(walData, datamining.Endorsement{}, nil)

	res := BuildWalletResult(wallet, bioWallet)
	assert.Equal(t, "encrypted aes key", string(res.EncryptedAESkey))
	assert.Equal(t, "encrypted wallet", string(res.EncryptedWallet))
	assert.Equal(t, "encrypted wallet addr", string(res.EncryptedWalletAddress))
}

/*
Scenario: Build a wallet from decoded data
	Given a wallet JSON data decoded
	When I want to build a wallet for the domain
	Then I get a new instance of data wallet
*/
func TestBuildWalletData(t *testing.T) {

	walJSON := WalletDataFromJSON{
		BiodPublicKey:      "pub key",
		EncryptedAddrRobot: "encrypted addr",
		EncryptedWallet:    "encrypted wallet",
		PersonPublicKey:    "pub key",
	}

	w := BuildWalletData(walJSON, &api.Signature{
		Biod:   []byte("bio sig"),
		Person: []byte("em sig"),
	})
	assert.Equal(t, "pub key", string(w.BiodPubk))
	assert.Equal(t, "pub key", string(w.EmPubk))
	assert.Equal(t, "encrypted addr", string(w.CipherAddrRobot))
	assert.Equal(t, "encrypted wallet", string(w.CipherWallet))
	assert.Equal(t, "bio sig", string(w.Sigs.BiodSig))
	assert.Equal(t, "em sig", string(w.Sigs.EmSig))
}

/*
Scenario: Build a bio wallet from decoded data
	Given a bio JSON data decoded
	When I want to build a wallet for the domain
	Then I get a new instance of bio wallet
*/
func TestBuildBioWalletData(t *testing.T) {
	bioJSON := BioDataFromJSON{
		BiodPublicKey:       "pub key",
		EncryptedAddrPerson: "encrypted addr",
		EncryptedAddrRobot:  "encrypted addr",
		EncryptedAESKey:     "encrypted aes key",
		PersonHash:          "person hash",
		PersonPublicKey:     "pub key",
	}

	bw := BuildBioData(bioJSON, &api.Signature{
		Biod:   []byte("bio sig"),
		Person: []byte("em sig"),
	})
	assert.Equal(t, "pub key", string(bw.BiodPubk))
	assert.Equal(t, "pub key", string(bw.EmPubk))
	assert.Equal(t, "encrypted addr", string(bw.CipherAddrRobot))
	assert.Equal(t, "encrypted addr", string(bw.CipherAddrBio))
	assert.Equal(t, "encrypted aes key", string(bw.CipherAESKey))
	assert.Equal(t, "person hash", string(bw.BHash))
	assert.Equal(t, "bio sig", string(bw.Sigs.BiodSig))
	assert.Equal(t, "em sig", string(bw.Sigs.EmSig))
}
