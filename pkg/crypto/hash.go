package crypto

import (
	"crypto"
	"crypto/sha256"
)

//DefaultHashAlgo identifies the default hash algorithm
const DefaultHashAlgo crypto.Hash = crypto.SHA256

type hashFunc func(data []byte) []byte

var supportedHash = map[crypto.Hash]hashFunc{
	crypto.SHA256: sha256Func,
}

//VersionnedHash identifies a digest for a hash marshaled with its algorithm as first byte
type VersionnedHash []byte

//NewVersionnedHash creates a new versionned hash by adding firstly the algorithm used then the digest
func NewVersionnedHash(algo crypto.Hash, digest []byte) VersionnedHash {
	out := make(VersionnedHash, 1+len(digest))
	out[0] = byte(int(algo))
	copy(out[1:], digest)
	return out
}

//Algorithm returns the algorithm used
func (h VersionnedHash) Algorithm() crypto.Hash {
	return crypto.Hash(h[0])
}

//Digest returns the digest
func (h VersionnedHash) Digest() []byte {
	return h[1:]
}

//IsValid checks if it is valid by verifying the emptyness and the size of the hash
func (h VersionnedHash) IsValid() (valid bool) {
	defer func() {
		if r := recover(); r != nil {
			valid = false
			return
		}
	}()
	if len(h) == 0 {
		return false
	}

	return h.Algorithm().Size() == len(h.Digest())
}

//Hash computes a hash checksum using the default hash algorithm and preceed with the algorithm used
func Hash(data []byte) VersionnedHash {
	hash := supportedHash[DefaultHashAlgo](data)
	return NewVersionnedHash(DefaultHashAlgo, hash)
}

func sha256Func(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
