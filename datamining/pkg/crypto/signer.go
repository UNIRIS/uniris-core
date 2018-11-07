package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
)

type ecdsaSignature struct {
	R, S *big.Int
}

//Signer defines methods to handle signatures
type Signer interface {
	mining.Signer
	accountMining.KeychainSigner
	accountMining.BiometricSigner
}

type signer struct{}

//NewSigner creates a new signer
func NewSigner() Signer {
	return signer{}
}

func (s signer) CheckTransactionSignature(pubk string, txHash string, sig string) error {
	return s.checkSignature(pubk, txHash, sig)
}

func (s signer) CheckBiometricSignature(pubk string, data accountMining.UnsignedBiometricData, sig string) error {
	return s.checkSignature(pubk, data, sig)
}

func (s signer) CheckKeychainSignature(pubk string, data accountMining.UnsignedKeychainData, sig string) error {
	return s.checkSignature(pubk, data, sig)
}

//Verify verify a signature and a data using a public key
func (s signer) checkSignature(pubk string, data interface{}, sig string) error {
	var signature ecdsaSignature

	decodedkey, err := hex.DecodeString(pubk)
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

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	hash := []byte(hashBytes(b))

	if ecdsa.Verify(ecdsaPublic, hash, signature.R, signature.S) {
		return nil
	}

	return errors.New("Invalid signature")
}

func (s signer) SignValidation(v mining.UnsignedValidation, pvKey string) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return sign(pvKey, string(b))
}

func (s signer) SignLock(txLock lock.TransactionLock, pvKey string) (string, error) {
	b, err := json.Marshal(txLock)
	if err != nil {
		return "", err
	}

	return sign(pvKey, string(b))
}

func sign(privk string, data string) (string, error) {
	pvDecoded, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	pv, err := x509.ParseECPrivateKey(pvDecoded)
	if err != nil {
		return "", err
	}

	hash := []byte(hashString(data))

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
