package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	"github.com/uniris/ecies/pkg"
)

//Encrypt encrypt data using a public key
func Encrypt(pubk string, data string) (string, error) {
	decodeKey, err := hex.DecodeString(string(pubk))
	if err != nil {
		return "", err
	}

	pu, err := x509.ParsePKIXPublicKey(decodeKey)
	if err != nil {
		return "", err
	}
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	cipher, err := ecies.Encrypt(rand.Reader, robotEciesKey, []byte(data), nil, nil)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipher), nil
}

//Decrypt decrypt data using a private key
func Decrypt(privk string, edata string) (string, error) {
	decodeKey, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	robotKey, err := x509.ParseECPrivateKey(decodeKey)
	if err != nil {
		return "", err
	}

	decodeCipher, err := hex.DecodeString(edata)
	if err != nil {
		return "", err
	}

	robotEciesKey := ecies.ImportECDSA(robotKey)
	message, err := robotEciesKey.Decrypt(decodeCipher, nil, nil)
	if err != nil {
		return "", err
	}
	return string(message), nil
}
