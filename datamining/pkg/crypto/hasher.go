package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
)

//Hasher defines methods for hashing
type Hasher interface {
	accountMining.KeychainHasher
	accountMining.BiometricHasher
	internalrpc.Hasher
}

type hasher struct{}

//NewHasher creates new hasher
func NewHasher() Hasher {
	return hasher{}
}

func (h hasher) HashBiometricJSON(data *internalrpc.BioDataFromJSON) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashKeychainJSON(data *internalrpc.KeychainDataFromJSON) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashUnsignedKeychainData(data accountMining.UnsignedKeychainData) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func (h hasher) HashUnsignedBiometricData(data accountMining.UnsignedBiometricData) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return hashBytes(b), nil
}

func hashString(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func hashBytes(data []byte) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
