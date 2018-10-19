package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

type ecdsaSignature struct {
	R, S *big.Int
}

type signer struct{}

//NewSigner creates a new signer
func NewSigner() validating.Signer {
	return signer{}
}

//Verify verify a signature and a data using a public key
func (s signer) CheckSignature(pubk string, data interface{}, der string) error {
	var signature ecdsaSignature

	decodedkey, err := hex.DecodeString(pubk)
	if err != nil {
		return err
	}

	decodedsig, err := hex.DecodeString(der)
	if err != nil {
		return err
	}

	pu, err := x509.ParsePKIXPublicKey(decodedkey)
	if err != nil {
		return err
	}

	ecdsaPublic := pu.(*ecdsa.PublicKey)
	asn1.Unmarshal(decodedsig, &signature)

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// log.Print(string(b))

	hash := []byte(HashBytes(b))

	if ecdsa.Verify(ecdsaPublic, hash, signature.R, signature.S) {
		return nil
	}

	return validating.ErrInvalidSignature
}

//Sign data using a privatekey
func Sign(privk string, data string) (string, error) {
	pvDecoded, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	pv, err := x509.ParseECPrivateKey(pvDecoded)
	if err != nil {
		return "", err
	}

	hash := []byte(HashString(data))

	r, s, err := ecdsa.Sign(rand.Reader, pv, hash)
	if err != nil {
		return "", err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil

}
