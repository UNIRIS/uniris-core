package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

//HashString computes a sha256 hash from a string value
func HashString(data string) string {
	return HashBytes([]byte(data))
}

//HashBytes computes a sha256 hash from a byte slice
func HashBytes(data []byte) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
