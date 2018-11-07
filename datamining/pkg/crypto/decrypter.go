package crypto

import (
	"crypto/x509"
	"encoding/hex"

	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"

	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

	"github.com/uniris/ecies/pkg"
)

//Decrypter defines methods to handle decryption
type Decrypter interface {
	rpc.Decrypter
	internalrpc.Decrypter
}

type decrypter struct{}

//NewDecrypter create a new decrypter
func NewDecrypter() Decrypter {
	return decrypter{}
}

func (d decrypter) DecryptCipherAddress(addr string, pvKey string) (string, error) {
	return decrypt(pvKey, addr)
}

func (d decrypter) DecryptHashPerson(hash string, pvKey string) (string, error) {
	return decrypt(pvKey, hash)
}

func (d decrypter) DecryptTransactionData(data string, pvKey string) (string, error) {
	return decrypt(pvKey, data)
}

func decrypt(privk string, data string) (string, error) {
	decodeKey, err := hex.DecodeString(privk)
	if err != nil {
		return "", err
	}

	robotKey, err := x509.ParseECPrivateKey(decodeKey)
	if err != nil {
		return "", err
	}

	decodeCipher, err := hex.DecodeString(data)
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
