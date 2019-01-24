package crypto

import (
	"crypto/x509"
	"encoding/hex"
)

func GetPublicKey(pvKey string) (string, error) {

	decodeKey, err := hex.DecodeString(pvKey)
	if err != nil {
		return "", err
	}

	key, err := x509.ParseECPrivateKey(decodeKey)
	if err != nil {
		return "", err
	}

	bPub, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bPub), nil
}
