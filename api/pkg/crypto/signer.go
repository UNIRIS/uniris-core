package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
)

//SignData creates signature from a struct
func SignData(privk string, data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return HashAndSign(privk, string(b))
}

//SignRaw creates a signature from hashed data
func SignRaw(privk string, hash string) (string, error) {
	pvDecoded, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	pv, err := x509.ParseECPrivateKey(pvDecoded)
	if err != nil {
		return "", err
	}

	r, s, err := ecdsa.Sign(rand.Reader, pv, []byte(hash))
	if err != nil {
		return "", err
	}

	sig, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}

//HashAndSign hash data and sign it
func HashAndSign(privk string, data string) (string, error) {
	return SignRaw(privk, Hash(data))
}
