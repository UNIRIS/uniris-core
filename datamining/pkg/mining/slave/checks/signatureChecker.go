package checks

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

type sigCheck struct {
	sig Signer
}

//NewSignatureChecker creates a signature validation
func NewSignatureChecker(sig Signer) Handler {
	return sigCheck{sig}
}

func (c sigCheck) IsCatchedError(err error) bool {
	return err == ErrInvalidSignature
}

func (c sigCheck) CheckData(data interface{}) error {

	switch data.(type) {
	case *datamining.KeyChainData:
		return c.checkKeychainData(data.(*datamining.KeyChainData))
	case *datamining.BioData:
		return c.checkBioData(data.(*datamining.BioData))
	}

	return nil
}

func (c sigCheck) checkKeychainData(kc *datamining.KeyChainData) error {
	kcValid := KeychainData{
		BIODPublicKey:      kc.BiodPubk,
		EncryptedAddrRobot: kc.CipherAddrRobot,
		EncryptedWallet:    kc.CipherWallet,
		PersonPublicKey:    kc.PersonPubk,
	}

	if err := c.sig.CheckSignature(kc.BiodPubk, kcValid, kc.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := c.sig.CheckSignature(kc.PersonPubk, kcValid, kc.Sigs.PersonSig); err != nil {
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
		PersonHash:          b.PersonHash,
		PersonPublicKey:     b.PersonPubk,
	}

	if err := c.sig.CheckSignature(b.BiodPubk, bValid, b.Sigs.BiodSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	if err := c.sig.CheckSignature(b.PersonPubk, bValid, b.Sigs.PersonSig); err != nil {
		if err == ErrInvalidSignature {
			return ErrInvalidSignature
		}
		return err
	}

	return nil
}
