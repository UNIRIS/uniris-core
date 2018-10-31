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

	sigCheck := NewSignatureChecker(mockSigCheckSigner{})

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	w := &datamining.KeyChainData{
		BiodPubk:   hex.EncodeToString(pbKey),
		PersonPubk: hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err := sigCheck.CheckData(w)
	assert.Nil(t, err)

	bd := &datamining.BioData{
		BiodPubk:   hex.EncodeToString(pbKey),
		PersonPubk: hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err = sigCheck.CheckData(bd)
	assert.Nil(t, err)
}

func TestKOSignatureValidator(t *testing.T) {

	sigCheck := NewSignatureChecker(mockBadSigCheckSigner{})

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	w := &datamining.KeyChainData{
		BiodPubk:   hex.EncodeToString(pbKey),
		PersonPubk: hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err := sigCheck.CheckData(w)
	assert.Equal(t, err, ErrInvalidSignature)

	bd := &datamining.BioData{
		BiodPubk:   hex.EncodeToString(pbKey),
		PersonPubk: hex.EncodeToString(pbKey),
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	err = sigCheck.CheckData(bd)
	assert.Equal(t, err, ErrInvalidSignature)
}

type mockSigCheckSigner struct{}

func (s mockSigCheckSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return nil
}

type mockBadSigCheckSigner struct{}

func (s mockBadSigCheckSigner) CheckSignature(pubk string, data interface{}, der string) error {
	return ErrInvalidSignature
}
