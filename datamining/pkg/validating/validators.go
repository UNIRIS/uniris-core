package validating

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg"
)

//DataValidator defines methods to valid wallet and bio data
type DataValidator interface {
	ValidWallet(*datamining.WalletData) error
	ValidBioWallet(*datamining.BioData) error
}

//Signer defines methods to handle signatures
type Signer interface {
	CheckSignature(pubKey string, data interface{}, sig string) error
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

func (v signatureValidator) ValidWallet(w *datamining.WalletData) error {
	wValid := WalletData{
		BIODPublicKey:      w.BiodPubk,
		EncryptedAddrRobot: w.CipherAddrRobot,
		EncryptedWallet:    w.CipherWallet,
		PersonPublicKey:    w.EmPubk,
	}

	if err := v.sig.CheckSignature(w.BiodPubk, wValid, w.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := v.sig.CheckSignature(w.EmPubk, wValid, w.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	return nil
}

func (v signatureValidator) ValidBioWallet(b *datamining.BioData) error {
	bValid := BioData{
		BIODPublicKey:       b.BiodPubk,
		EncryptedAddrPerson: b.CipherAddrBio,
		EncryptedAddrRobot:  b.CipherAddrRobot,
		EncryptedAESKey:     b.CipherAESKey,
		PersonHash:          b.BHash,
		PersonPublicKey:     b.EmPubk,
	}

	if err := v.sig.CheckSignature(b.BiodPubk, bValid, b.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := v.sig.CheckSignature(b.EmPubk, bValid, b.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	return nil
}
