package checks

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg"
)

func TestOkSignatureValidator(t *testing.T) {

	val := sigCheck{
		sig: mockSigner{},
	}

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

	err := val.CheckWalletData(w)
	assert.Nil(t, err)

	bd := &datamining.BioData{
		BiodPubk: hex.EncodeToString(pbKey),
		EmPubk:   hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err = val.CheckBioData(bd)
	assert.Nil(t, err)
}

func TestKOSignatureValidator(t *testing.T) {

	val := sigCheck{
		sig: badMockSigner{},
	}

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

	err := val.CheckWalletData(w)
	assert.Equal(t, err, ErrInvalidSignature)

	bd := &datamining.BioData{
		BiodPubk: hex.EncodeToString(pbKey),
		EmPubk:   hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig: "fake sig",
			EmSig:   "fake sig",
		},
	}

	err = val.CheckBioData(bd)
	assert.Equal(t, err, ErrInvalidSignature)
}

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

type badMockSigner struct{}

func (s badMockSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return ErrInvalidSignature
}
