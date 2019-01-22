package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"math/big"
)

type ecdsaSignature struct {
	R, S *big.Int
}

var ErrInvalidSignature = errors.New("Invalid signature")

//Sign creates a signature from a given data
func Sign(data string, pvKey string) (string, error) {
	pvDecoded, err := hex.DecodeString(pvKey)
	if err != nil {
		return "", err
	}

	pv, err := x509.ParseECPrivateKey(pvDecoded)
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, pv, []byte(HashString(data)))
	if err != nil {
		return "", err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}

//VerifySignature checks if the signature is valid by providing the public key and the data
func VerifySignature(data string, pubKey string, sig string) error {
	var signature ecdsaSignature

	decodedkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return err
	}

	decodedsig, err := hex.DecodeString(sig)
	if err != nil {
		return err
	}

	pu, err := x509.ParsePKIXPublicKey(decodedkey)
	if err != nil {
		return err
	}

	ecdsaPublic := pu.(*ecdsa.PublicKey)
	asn1.Unmarshal(decodedsig, &signature)

	if ecdsa.Verify(ecdsaPublic, []byte(HashString(data)), signature.R, signature.S) {
		return nil
	}

	return ErrInvalidSignature
}
