package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
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

func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return []byte(hex.EncodeToString(hash.Sum(nil)))
}

type ecdsaSignature struct {
	R, S *big.Int
}

//CheckSignature validates a signature from a data
func (v RequestValidator) CheckSignature(data interface{}, pubKey []byte, sig []byte) (bool, error) {
	var sigBio ecdsaSignature
	clearSig, err := hex.DecodeString(string(sig))
	if err != nil {
		return false, err
	}

	asn1.Unmarshal(clearSig, &sigBio)
	b, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	decodeKey, err := hex.DecodeString(string(pubKey))
	if err != nil {
		return false, err
	}

	key, err := x509.ParsePKIXPublicKey(decodeKey)
	ecdsaKey := key.(*ecdsa.PublicKey)

	if err != nil {
		return false, err
	}

	return ecdsa.Verify(ecdsaKey, hash(b), sigBio.R, sigBio.S), nil
}
