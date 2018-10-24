package internalrpc

import (
	"testing"

	"github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Build a wallet search result
	Given a wallet and bio data
	When I want build a wallet search result
	Then I get a new object with encrypted aes, wallet and addr given
*/
func TestBuildWalletSearchResult(t *testing.T) {

	bioData := &datamining.BioData{
		CipherAESKey:  "encrypted aes key",
		CipherAddrBio: "encrypted wallet addr",
	}
	bioWallet := datamining.NewBioWallet(bioData, &datamining.Endorsement{})

	walData := &datamining.WalletData{
		CipherWallet: "encrypted wallet",
	}

	wallet := datamining.NewWallet(walData, &datamining.Endorsement{}, "")

	res := BuildWalletSearchResult(wallet, bioWallet)
	assert.Equal(t, "encrypted aes key", res.EncryptedAESkey)
	assert.Equal(t, "encrypted wallet", res.EncryptedWallet)
	assert.Equal(t, "encrypted wallet addr", res.EncryptedWalletAddress)
}

/*
Scenario: Build a wallet from decoded data
	Given a wallet JSON data decoded
	When I want to build a wallet for the domain
	Then I get a new instance of data wallet
*/
func TestBuildWalletData(t *testing.T) {

	walJSON := &WalletDataFromJSON{
		BiodPublicKey:      "pub key",
		EncryptedAddrRobot: "encrypted addr",
		EncryptedWallet:    "encrypted wallet",
		PersonPublicKey:    "pub key",
	}

	w := BuildWalletData(walJSON, &api.Signature{
		Biod:   "bio sig",
		Person: "em sig",
	}, "addr")
	assert.Equal(t, "pub key", w.BiodPubk)
	assert.Equal(t, "pub key", w.EmPubk)
	assert.Equal(t, "encrypted addr", w.CipherAddrRobot)
	assert.Equal(t, "encrypted wallet", w.CipherWallet)
	assert.Equal(t, "bio sig", w.Sigs.BiodSig)
	assert.Equal(t, "em sig", w.Sigs.EmSig)
	assert.Equal(t, "addr", w.WalletAddr)
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
		Biod:   "bio sig",
		Person: "em sig",
	})
	assert.Equal(t, "pub key", bw.BiodPubk)
	assert.Equal(t, "pub key", bw.EmPubk)
	assert.Equal(t, "encrypted addr", bw.CipherAddrRobot)
	assert.Equal(t, "encrypted addr", bw.CipherAddrBio)
	assert.Equal(t, "encrypted aes key", bw.CipherAESKey)
	assert.Equal(t, "person hash", bw.BHash)
	assert.Equal(t, "bio sig", bw.Sigs.BiodSig)
	assert.Equal(t, "em sig", bw.Sigs.EmSig)
}
