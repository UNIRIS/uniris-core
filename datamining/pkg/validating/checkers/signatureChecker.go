package checkers

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

//Checker defines methods for transaction data checking
type Checker interface {
	CheckData(interface{}) error
	IsCatchedError(error) bool
}

type sigCheck struct {
	sig Signer
}

//NewSignatureChecker creates a signature checker
func NewSignatureChecker(sig Signer) Checker {
	return sigCheck{sig}
}

func (c sigCheck) IsCatchedError(err error) bool {
	return err == ErrInvalidSignature
}

func (c sigCheck) CheckData(data interface{}) error {

	switch data.(type) {
	case *datamining.WalletData:
		return c.checkWalletData(data.(*datamining.WalletData))
	case *datamining.BioData:
		return c.checkBioData(data.(*datamining.BioData))
	}

	return nil
}

func (c sigCheck) checkWalletData(w *datamining.WalletData) error {
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

func (c sigCheck) checkBioData(b *datamining.BioData) error {
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
