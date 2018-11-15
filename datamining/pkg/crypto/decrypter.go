package crypto

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

	"github.com/uniris/ecies/pkg"
)

//Decrypter defines methods to handle decryption
type Decrypter interface {
	rpc.Decrypter
}

type decrypter struct{}

//NewDecrypter create a new decrypter
func NewDecrypter() Decrypter {
	return decrypter{}
}

func (d decrypter) DecryptHash(hash string, pvKey string) (string, error) {
	return decrypt(pvKey, hash)
}

func (d decrypter) DecryptBiometricData(data string, pvKey string) (account.BiometricData, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var bio biometricRaw
	err = json.Unmarshal([]byte(clear), &bio)
	if err != nil {
		return nil, err
	}
	return account.NewBiometricData(
		bio.PersonHash,
		bio.EncryptedAddrRobot,
		bio.EncryptedAddrPerson,
		bio.EncryptedAESKey,
		bio.PersonPublicKey, nil), nil
}

func (d decrypter) DecryptKeychainData(data string, pvKey string) (account.KeychainData, error) {
	clear, err := decrypt(pvKey, data)
	if err != nil {
		return nil, err
	}
	var keychain keychainRaw
	err = json.Unmarshal([]byte(clear), &keychain)
	if err != nil {
		return nil, err
	}
	return account.NewKeychainData(
		keychain.EncryptedAddrRobot,
		keychain.EncryptedWallet,
		keychain.PersonPublicKey, nil), nil
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
