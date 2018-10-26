package validations

import (
	"errors"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//ErrInvalidSignature is returned when a invalid signature is provided
var ErrInvalidSignature = errors.New("Invalid signature")

//Signer defines methods to handle signatures
type Signer interface {
	CheckSignature(pubKey string, data interface{}, sig string) error
}

type sigValid struct {
	sig Signer
}

//NewSignatureValidation creates a signature validation
func NewSignatureValidation(sig Signer) Handler {
	return sigValid{sig}
}

func (c sigValid) IsCatchedError(err error) bool {
	return err == ErrInvalidSignature
}

func (c sigValid) CheckData(data interface{}) error {

	switch data.(type) {
	case *datamining.WalletData:
		return c.checkWalletData(data.(*datamining.WalletData))
	case *datamining.BioData:
		return c.checkBioData(data.(*datamining.BioData))
	}

	return nil
}

func (c sigValid) checkWalletData(w *datamining.WalletData) error {
	wValid := WalletData{
		BIODPublicKey:      w.BiodPubk,
		EncryptedAddrRobot: w.CipherAddrRobot,
		EncryptedWallet:    w.CipherWallet,
		PersonPublicKey:    w.EmPubk,
	}

	if err := c.sig.CheckSignature(w.BiodPubk, wValid, w.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := c.sig.CheckSignature(w.EmPubk, wValid, w.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}
	return nil
}

func (c sigValid) checkBioData(b *datamining.BioData) error {
	bValid := BioData{
		BIODPublicKey:       b.BiodPubk,
		EncryptedAddrPerson: b.CipherAddrBio,
		EncryptedAddrRobot:  b.CipherAddrRobot,
		EncryptedAESKey:     b.CipherAESKey,
		PersonHash:          b.BHash,
		PersonPublicKey:     b.EmPubk,
	}

	if err := c.sig.CheckSignature(b.BiodPubk, bValid, b.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := c.sig.CheckSignature(b.EmPubk, bValid, b.Sigs.EmSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	return nil
}
