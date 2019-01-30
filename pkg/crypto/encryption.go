package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	ecies "github.com/uniris/ecies/pkg"
)

func Encrypt(data string, pubKey string) (string, error) {
	decodeKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}

	key, err := x509.ParsePKIXPublicKey(decodeKey)
	if err != nil {
		return "", err
	}

	eciesKey := ecies.ImportECDSAPublic(key.(*ecdsa.PublicKey))
	cipher, err := ecies.Encrypt(rand.Reader, eciesKey, []byte(data), nil, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(cipher), nil
}

//Decrypt use ECIES to decrypt a data
func Decrypt(data string, pvKey string) (string, error) {

	decodeKey, err := hex.DecodeString(pvKey)
	if err != nil {
		return "", err
	}

	key, err := x509.ParseECPrivateKey(decodeKey)
	if err != nil {
		return "", err
	}

	decodeCipher, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	eciesKey := ecies.ImportECDSA(key)
	message, err := eciesKey.Decrypt(decodeCipher, nil, nil)
	if err != nil {
		return "", err
	}
	return string(message), nil
}
