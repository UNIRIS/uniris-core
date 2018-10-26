package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//Hasher defines methods for hashing
type Hasher interface {
	mining.PreviousDataHasher
}

type hasher struct{}

//NewHasher creates new hasher
func NewHasher() Hasher {
	return hasher{}
}

func (h hasher) HashWallet(w *datamining.Wallet) (string, error) {
	b, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	return HashBytes(b), nil
}

//HashString return the hash of a string data
func HashString(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

//HashBytes return the hash of a bytes data
func HashBytes(data []byte) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
