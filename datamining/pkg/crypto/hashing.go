package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

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
