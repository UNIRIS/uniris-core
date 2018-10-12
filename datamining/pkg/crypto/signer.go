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

//ErrBadSignature define a bad signature error
var ErrBadSignature = errors.New("Error: Bad Signature")

//Signer defines signatures methods
type Signer struct{}

type ecdsaSignature struct {
	R, S *big.Int
}

//Verify verify a signature and a data using a public key
func (si Signer) Verify(pubk []byte, der []byte, hash []byte) error {
	var signature ecdsaSignature

	decodedsig, _ := hex.DecodeString(string(der))

	pu, err := x509.ParsePKIXPublicKey(pubk)
	if err != nil {
		return err
	}

	ecdsaPublic := pu.(*ecdsa.PublicKey)
	asn1.Unmarshal(decodedsig, &signature)

	if ecdsa.Verify(ecdsaPublic, hash, signature.R, signature.S) {
		return nil
	}

	return ErrBadSignature
}

//Sign data using a privatekey
func (si Signer) Sign(privk []byte, data []byte) ([]byte, error) {
	pv, err := x509.ParseECPrivateKey(privk)
	if err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, pv, data)
	if err != nil {
		return nil, err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(sig)), nil

}
