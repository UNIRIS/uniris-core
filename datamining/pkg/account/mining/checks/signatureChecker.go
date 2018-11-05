package checks

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//ErrInvalidSignature is returned when a invalid signature is provided
var ErrInvalidSignature = errors.New("Invalid signature")

type sigCheck struct {
	sig Signer
}

//NewSignatureChecker creates a signature validation
func NewSignatureChecker(sig Signer) Handler {
	return sigCheck{sig}
}

func (c sigCheck) CheckData(data interface{}, txHash string) error {
	switch data.(type) {
	case *account.KeyChainData:
		return c.checkKeychainData(data.(*account.KeyChainData))
	case *account.BioData:
		return c.checkBioData(data.(*account.BioData))
	}

	return errors.New("Unsupported data")
}

func (c sigCheck) checkKeychainData(kc *account.KeyChainData) error {
	kcValid := rawKeychainData{
		BiodPublicKey:      kc.BiodPubk,
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

func (c sigCheck) checkBioData(b *account.BioData) error {
	bValid := rawBiometricData{
		BiodPublicKey:       b.BiodPubk,
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
