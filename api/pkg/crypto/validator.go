package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/uniris/uniris-core/api/pkg/listing"
)

//RequestValidator implements the request validator
type RequestValidator struct {
	listing.RequestValidator
}

type ecdsaSignature struct {
	R, S *big.Int
}

//CheckRawSignature validate signed hashed data
func (v RequestValidator) CheckRawSignature(hashedData string, pubKey string, sig string) (bool, error) {
	decodeKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, err
	}

	key, err := x509.ParsePKIXPublicKey(decodeKey)
	ecdsaKey := key.(*ecdsa.PublicKey)

	var sigBio ecdsaSignature
	clearSig, err := hex.DecodeString(sig)
	if err != nil {
		return false, err
	}

	asn1.Unmarshal(clearSig, &sigBio)

	return ecdsa.Verify(ecdsaKey, []byte(hashedData), sigBio.R, sigBio.S), nil
}

//CheckDataSignature validates a signature from a data
func (v RequestValidator) CheckDataSignature(data interface{}, pubKey string, sig string) (bool, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	return v.CheckRawSignature(Hash(string(b)), pubKey, sig)
}
