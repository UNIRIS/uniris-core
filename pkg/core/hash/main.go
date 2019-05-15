package main

import (
	"crypto"
	"crypto/sha256"
)

//DefaultHashAlgo identifies the default hash algorithm
var DefaultHashAlgo = crypto.SHA256

type hashFunc func(data []byte) []byte

var supportedHash = map[crypto.Hash]hashFunc{
	crypto.SHA256: sha256Func,
}

func versionHash(algo crypto.Hash, digest []byte) []byte {
	out := make([]byte, 1+len(digest))
	out[0] = byte(int(algo))
	copy(out[1:], digest)
	return out
}

//Hash computes a hash checksum using the default hash algorithm and preceed with the algorithm used
func Hash(data []byte) []byte {
	hash := supportedHash[DefaultHashAlgo](data)
	return versionHash(DefaultHashAlgo, hash)
}

//HashAlgorithm returns the algorithm of the versionned hash
func HashAlgorithm(hash []byte) crypto.Hash {
	return crypto.Hash(hash[0])
}

//HashDigest returns the digest of a versionned hash
func HashDigest(hash []byte) []byte {
	return hash[1:]
}

//IsValidHash checks if it is valid by verifying the emptyness and the size of the hash
func IsValidHash(hash []byte) (valid bool) {
	defer func() {
		if r := recover(); r != nil {
			valid = false
			return
		}
	}()
	if len(hash) == 0 {
		return false
	}

	return HashAlgorithm(hash).Size() == len(HashDigest(hash))
}

func sha256Func(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
