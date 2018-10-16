package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

//Hash a data
func Hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return []byte(hex.EncodeToString(hash.Sum(nil)))
}
