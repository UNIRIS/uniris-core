package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"errors"
	"io"
)

//Curve identifies a supported elliptic curve used for keys
type Curve int

const (

	//Ed25519Curve identifies the Ed25519 elliptic curve
	Ed25519Curve Curve = 0

	//P256Curve identifies the NIST-P256 elliptic curve
	P256Curve Curve = 1
)

type generateSharedFunc func(pub PublicKey, pv PrivateKey) (secret []byte, err error)

//VersionnedKey represents an elliptic key versionned by its curve
type VersionnedKey []byte

//Curve returns the elliptic curve of the key
func (k VersionnedKey) Curve() Curve {
	return Curve(k[0])
}

//Marshalling returns the marshalled key
func (k VersionnedKey) Marshalling() []byte {
	return k[1:]
}

type key interface {

	//Marshal exports the key using marshalling and versionning based on the elliptic curve
	//First byte identify the curve and the rest the key marshalling
	Marshal() (VersionnedKey, error)

	bytes() []byte
	curve() Curve
}

//PrivateKey represents an elliptic private key
type PrivateKey interface {
	key

	//Sign creates a new signature from a given data
	//The signing algorithm will be choose depending from the private key elliptic curve
	Sign(data []byte) ([]byte, error)

	//Decrypt uses ECIES to decipher the encrypted data
	//An extract of the ECDH shared key, the encrypted message and the message authentication code is made
	//Finally the authentication code is check and the encrypted message is decrypted with AES
	Decrypt(Cipher) ([]byte, error)
}

//PublicKey represents an elliptic public key
type PublicKey interface {
	key

	//Verifies check the signature for the given data
	//The verfying algorithm will be choose depending from the public key elliptic curve
	Verify(data []byte, sig []byte) bool

	//Encrypt uses ECIES to cipher the given data
	//ECDH shared key is generated, the data is encrypted with AES and an message authentication code is made
	//Finally the all these fields are encoded
	Encrypt(data []byte) (Cipher, error)

	//Equals compares the current key with another
	Equals(pub PublicKey) bool
}

func versionKey(c Curve, marshalling []byte) VersionnedKey {
	out := make(VersionnedKey, 1+len(marshalling))
	out[0] = byte(int(c))
	copy(out[1:], marshalling)
	return out
}

//GenerateECKeyPair creates a new elliptic keypair from a given specific curve
func GenerateECKeyPair(c Curve, src io.Reader) (PrivateKey, PublicKey, error) {
	switch c {
	case P256Curve:
		return generateECDSAKeys(src, elliptic.P256())
	case Ed25519Curve:
		return generateEd25519Keys(src)
	default:
		return nil, nil, errors.New("unsupported EC curve")
	}
}

//ParsePublicKey converts a marshalled versionned public key into a PublicKey
func ParsePublicKey(key VersionnedKey) (PublicKey, error) {
	switch key.Curve() {
	case P256Curve:
		pub, err := x509.ParsePKIXPublicKey(key.Marshalling())
		if err != nil {
			return nil, err
		}
		ecdsaPub := pub.(*ecdsa.PublicKey)
		return ecdsaPublicKey{ecdsaPub}, nil
	case Ed25519Curve:
		//TODO: need a way to check if it's valid
		return ed25519PublicKey{key.Marshalling()}, nil
	default:
		return nil, errors.New("unsupported curve")
	}
}

//ParsePrivateKey converts a marshalled versionned private key into a Private key
func ParsePrivateKey(key VersionnedKey) (PrivateKey, error) {
	switch key.Curve() {
	case P256Curve:
		ecdsaPriv, err := x509.ParseECPrivateKey(key.Marshalling())
		if err != nil {
			return nil, err
		}
		return ecdsaPrivateKey{ecdsaPriv}, nil
	case Ed25519Curve:
		return ed25519PrivateKey{key.Marshalling()}, nil
	default:
		return nil, errors.New("unsupported curve")
	}
}
