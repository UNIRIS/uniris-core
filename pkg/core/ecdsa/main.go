package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"io"
	"math/big"
)

type ecdsaSignature struct {
	R, S *big.Int
}

//Sign creates a signature for a data and a given private key . The data will be hashed
func Sign(key []byte, data []byte) ([]byte, error) {
	pvKey, err := extractPrivateKey(key)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, pvKey, hash[:])
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ecdsaSignature{r, s})
}

//Verify checks the signature of a given data and given public key. The data will be hashed
func Verify(key []byte, data []byte, sig []byte) (bool, error) {
	var ecdsaSig ecdsaSignature
	if _, err := asn1.Unmarshal(sig, &ecdsaSig); err != nil {
		return false, err
	}

	ecdsaPub, err := extractPublicKey(key)
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)

	return ecdsa.Verify(ecdsaPub, hash[:], ecdsaSig.R, ecdsaSig.S), nil
}

//GenerateKeys creates new ECDSA keypair using a specific curve
func GenerateKeys(src io.Reader, c elliptic.Curve) ([]byte, []byte, error) {

	priv, err := ecdsa.GenerateKey(c, src)
	if err != nil {
		return nil, nil, err
	}

	pvBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return nil, nil, err
	}

	return pvBytes, pubBytes, nil
}

//GenerateSharedSecret creates a shared secret by using the scalar multiplication on the Elliptic curve
func GenerateSharedSecret(pub []byte, pv []byte) ([]byte, error) {
	pubKey, err := extractPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pvKey, err := extractPrivateKey(pv)
	if err != nil {
		return nil, err
	}

	x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, pvKey.D.Bytes())
	return x.Bytes(), nil
}

func extractPublicKey(key []byte) (*ecdsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	switch pub.(type) {
	case *ecdsa.PublicKey:
		return pub.(*ecdsa.PublicKey), nil
	default:
		return nil, errors.New("invalid ecdsa public key")
	}
}

func extractPrivateKey(key []byte) (*ecdsa.PrivateKey, error) {
	pvKey, err := x509.ParseECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pvKey, nil
}

//ExtractMessagePublicKey finds a public key in a message
//It using the x509 Parse with ASN to determinate the public key
func ExtractMessagePublicKey(cipher []byte) ([]byte, int, error) {

	if len(cipher) < 91 {
		return nil, 0, errors.New("invalid ECIES cipher for ECDSA")
	}

	//X509.MarshalPKIXPublicKey returns 91 bytes for the P256 curve
	//So we need to extract it to avoid ASN trailing data error
	sample256 := cipher[:91]
	pub, err := x509.ParsePKIXPublicKey(sample256)
	if err != nil {
		return nil, 0, errors.New("unsupported ECDSA curve")
	}
	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, 0, errors.New("invalid ECDSA public key")
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(ecdsaPub)
	if err != nil {
		return nil, 0, err
	}

	return pubBytes, len(pubBytes), nil

}
