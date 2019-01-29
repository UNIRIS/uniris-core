package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
)

func IsPublicKey(pubKey string) (bool, error) {
	if pubKey == "" {
		return false, errors.New("public key is empty")
	}
	decodedkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, errors.New("public key is not in hexadecimal format")
	}

	pu, err := x509.ParsePKIXPublicKey(decodedkey)
	if err != nil {
		return false, errors.New("public key is not valid")
	}

	switch pu.(type) {
	case *ecdsa.PublicKey:
		return true, nil
	default:
		return false, errors.New("public key is not from an elliptic curve")
	}
}

func IsPrivateKey(pvKey string) (bool, error) {
	if pvKey == "" {
		return false, errors.New("private key is empty")
	}
	decodedkey, err := hex.DecodeString(pvKey)
	if err != nil {
		return false, errors.New("private key is not in hexadecimal format")
	}

	_, err = x509.ParseECPrivateKey(decodedkey)
	if err != nil {
		return false, errors.New("private key is not valid")
	}

	return true, nil
}

func GetPublicKeyFromPrivate(pvKey string) (string, error) {

	if _, err := IsPrivateKey(pvKey); err != nil {
		return "", err
	}

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
