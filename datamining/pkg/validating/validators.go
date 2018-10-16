package validating

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg"
)

//DataValidator defines methods to valid wallet and bio data
type DataValidator interface {
	ValidWallet(datamining.WalletData) (datamining.ValidationStatus, error)
	ValidBioWallet(datamining.BioData) (datamining.ValidationStatus, error)
}

//Signer defines methods to handle signatures
type Signer interface {
	CheckSignature(pubKey []byte, data interface{}, sig []byte) error
}

//ErrInvalidSignature is returned when a invalid signature is provided
var ErrInvalidSignature = errors.New("Invalid signature")

type signatureValidator struct {
	sig Signer
}

//NewSignatureValidator creates a signature validator
func NewSignatureValidator(sig Signer) DataValidator {
	return signatureValidator{sig}
}

func (v signatureValidator) ValidWallet(w datamining.WalletData) (datamining.ValidationStatus, error) {
	if err := v.sig.CheckSignature(w.BiodPubk, w, w.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return datamining.ValidationKO, nil
		}
		return datamining.ValidationKO, err
	}

	if err := v.sig.CheckSignature(w.EmPubk, w, w.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return datamining.ValidationKO, nil
		}
		return datamining.ValidationKO, err
	}

	return datamining.ValidationOK, nil
}

func (v signatureValidator) ValidBioWallet(b datamining.BioData) (datamining.ValidationStatus, error) {
	if err := v.sig.CheckSignature(b.BiodPubk, b, b.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return datamining.ValidationKO, nil
		}
		return datamining.ValidationKO, err
	}

	if err := v.sig.CheckSignature(b.EmPubk, b, b.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return datamining.ValidationKO, nil
		}
		return datamining.ValidationKO, err
	}

	return datamining.ValidationOK, nil
}
