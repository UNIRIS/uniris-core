package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"

	"github.com/uniris/ecies/pkg"
	robot "github.com/uniris/uniris-core/datamining/pkg"
)

type encrypter struct {
}

//Newencrypter creates a new encrypter service
func Newencrypter() (robot.Encrypter, error) {
	return encrypter{}, nil
}

//Decrypt decrypt data using a private key
func (e encrypter) Decrypt(privk []byte, edata []byte) ([]byte, error) {
	robotKey, err := x509.ParseECPrivateKey(privk)
	robotEciesKey := ecies.ImportECDSA(robotKey)
	message, err := robotEciesKey.Decrypt(edata, nil, nil)
	if err != nil {
		return nil, err
	}
	return message, nil
}

//Encrypt encrypt data using a public key
func (e encrypter) Ecrypt(pubk []byte, data []byte) ([]byte, error) {
	pu, err := x509.ParsePKIXPublicKey(pubk)
	if err != nil {
		return nil, err
	}
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	cipher, err := ecies.Encrypt(rand.Reader, robotEciesKey, data, nil, nil)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}
