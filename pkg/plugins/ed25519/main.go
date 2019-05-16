package main

import (
	"errors"
	"io"

	"github.com/uniris/ed25519/extra25519"
	"golang.org/x/crypto/curve25519"
	ed25519 "golang.org/x/crypto/ed25519"
)

//Sign creates a signature for a data and a given private key . The data will be hashed
func Sign(key []byte, data []byte) (sig []byte, err error) {
	if len(key) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid ed25519 private key")
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	sig = ed25519.Sign(key, data)
	return
}

//Verify checks the signature of a given data and given public key. The data will be hashed
func Verify(key []byte, data []byte, sig []byte) (valid bool, err error) {

	if len(key) != ed25519.PublicKeySize {
		return false, errors.New("invalid ed25519 public key")
	}

	defer func() {
		if r := recover(); r != nil {
			valid = false
			return
		}
	}()

	if len(sig) != ed25519.SignatureSize {
		return false, nil
	}

	return ed25519.Verify(key, data, sig), nil
}

//GenerateKeys creates new Ed25519 keypair
func GenerateKeys(src io.Reader) ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}

//GenerateSharedSecret creates a shared secret by using the scalar multiplication on the Elliptic curve
func GenerateSharedSecret(pub []byte, pv []byte) ([]byte, error) {
	curve25519Public, err := convertPublicEd25519ToCurve25519(pub)
	if err != nil {
		return nil, err
	}

	curve25519priv := convertPrivateEd25519ToCurve25519(pv)

	var s [ed25519.PublicKeySize]byte
	curve25519.ScalarMult(&s, &curve25519priv, &curve25519Public)
	return s[:], nil
}

//ExtractMessagePublicKey finds a public key in a message
//It using the ed25519 fixed size to identify the key
func ExtractMessagePublicKey(cipher []byte) ([]byte, int, error) {

	if len(cipher) < ed25519.PublicKeySize {
		return nil, 0, errors.New("invalid message")
	}
	pubBytes := cipher[0:ed25519.PublicKeySize]
	return pubBytes, ed25519.PublicKeySize, nil
}

func convertPrivateEd25519ToCurve25519(privKey ed25519.PrivateKey) [32]byte {
	var tmpprivKey [ed25519.PrivateKeySize]byte
	copy(tmpprivKey[:], privKey[:])
	var curve25519priv [ed25519.PublicKeySize]byte
	extra25519.PrivateKeyToCurve25519(&curve25519priv, &tmpprivKey)
	return curve25519priv
}

func convertPublicEd25519ToCurve25519(pubKey ed25519.PublicKey) ([ed25519.PublicKeySize]byte, error) {
	var tmpPubKey, curve25519Public [ed25519.PublicKeySize]byte
	copy(tmpPubKey[:], pubKey[:])
	if tmpPubKey == [ed25519.PublicKeySize]byte{} {
		return [ed25519.PublicKeySize]byte{}, errors.New("invalid public key")
	}
	if !extra25519.PublicKeyToCurve25519(&curve25519Public, &tmpPubKey) {
		return [ed25519.PublicKeySize]byte{}, errors.New("cannot convert ed25519 public key to curve25519 public key")
	}
	return curve25519Public, nil
}
