package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
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

func IsPublicKey(pubKey string) (bool, error) {
	if pubKey == "" {
		return false, errors.New("Public key is empty")
	}
	decodedkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, errors.New("Public key is not in hexadecimal format")
	}

	pu, err := x509.ParsePKIXPublicKey(decodedkey)
	if err != nil {
		return false, errors.New("Public key is not valid")
	}

	ecdsaPublic := pu.(*ecdsa.PublicKey)
	if ecdsaPublic == nil {
		return false, errors.New("Public key is not valid")
	}

	return true, nil
}
