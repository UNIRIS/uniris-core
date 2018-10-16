package validating

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

	val := signatureValidator{
		sig: mockSigner{},
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	w := datamining.WalletData{
		BiodPubk: []byte(hex.EncodeToString(pbKey)),
		EmPubk:   []byte(hex.EncodeToString(pbKey)),
		Sigs: datamining.Signatures{
			BiodSig: []byte("fake sig"),
			EmSig:   []byte("fake sig"),
		},
	}

	status, err := val.ValidWallet(w)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, status)

	bw := datamining.BioData{
		BiodPubk: []byte(hex.EncodeToString(pbKey)),
		EmPubk:   []byte(hex.EncodeToString(pbKey)),
		Sigs: datamining.Signatures{
			BiodSig: []byte("fake sig"),
			EmSig:   []byte("fake sig"),
		},
	}

	status, err = val.ValidBioWallet(bw)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationOK, status)
}

func TestKOSignatureValidator(t *testing.T) {

	val := signatureValidator{
		sig: badMockSigner{},
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	w := datamining.WalletData{
		BiodPubk: []byte(hex.EncodeToString(pbKey)),
		EmPubk:   []byte(hex.EncodeToString(pbKey)),
		Sigs: datamining.Signatures{
			BiodSig: []byte("fake sig"),
			EmSig:   []byte("fake sig"),
		},
	}

	status, err := val.ValidWallet(w)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationKO, status)

	bw := datamining.BioData{
		BiodPubk: []byte(hex.EncodeToString(pbKey)),
		EmPubk:   []byte(hex.EncodeToString(pbKey)),
		Sigs: datamining.Signatures{
			BiodSig: []byte("fake sig"),
			EmSig:   []byte("fake sig"),
		},
	}

	status, err = val.ValidBioWallet(bw)
	assert.Nil(t, err)
	assert.Equal(t, datamining.ValidationKO, status)
}

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk []byte, data interface{}, der []byte) error {
	return nil
}

type badMockSigner struct{}

func (s badMockSigner) CheckSignature(pubk []byte, data interface{}, der []byte) error {
	return ErrInvalidSignature
}
