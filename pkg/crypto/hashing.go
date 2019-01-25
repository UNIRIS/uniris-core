package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func IsHash(hash string) (bool, error) {
	if hash == "" {
		return false, errors.New("Hash is empty")
	}

	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return false, errors.New("Hash is not in hexadecimal format")
	}

	size := sha256.New().Size()
	if len(bytes) != size {
		return false, errors.New("Hash is not valid")
	}

	return true, nil
}
