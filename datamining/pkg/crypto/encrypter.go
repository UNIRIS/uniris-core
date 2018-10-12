package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	"github.com/uniris/ecies/pkg"
)

//Encrypter define encryption methods
type Encrypter struct {
}

//Encrypt encrypt data using a public key
func (e Encrypter) Encrypt(pubk []byte, data []byte) ([]byte, error) {
	pu, err := x509.ParsePKIXPublicKey(pubk)
	if err != nil {
		return nil, err
	}
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	cipher, err := ecies.Encrypt(rand.Reader, robotEciesKey, data, nil, nil)
	if err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(cipher)), nil
}

//Decrypt decrypt data using a private key
func (e Encrypter) Decrypt(privk []byte, edata []byte) ([]byte, error) {
	decodeCipher, _ := hex.DecodeString(string(edata))
	robotKey, err := x509.ParseECPrivateKey(privk)
	robotEciesKey := ecies.ImportECDSA(robotKey)
	message, err := robotEciesKey.Decrypt(decodeCipher, nil, nil)
	if err != nil {
		return nil, err
	}
	return message, nil
}
