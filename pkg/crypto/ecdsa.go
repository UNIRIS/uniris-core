package crypto

import (
	"bytes"
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

type ecdsaPrivateKey struct {
	priv *ecdsa.PrivateKey
}

func (ePriv ecdsaPrivateKey) Marshal() (VersionnedKey, error) {
	m, err := x509.MarshalECPrivateKey(ePriv.priv)
	if err != nil {
		return nil, err
	}
	return versionKey(ePriv.curve(), m), nil
}

func (ePriv ecdsaPrivateKey) Bytes() []byte {
	return elliptic.Marshal(ePriv.priv.Curve, ePriv.priv.X, ePriv.priv.Y)
}

func (ePriv ecdsaPrivateKey) curve() Curve {
	return getEcdsaCurve(ePriv.priv.Curve)
}

func (ePriv ecdsaPrivateKey) Decrypt(c Cipher) ([]byte, error) {
	return eciesDecrypt(c, ePriv, ecdsaGenerateShared, ecdsaExtractRandomPublicKey(ePriv.priv.Curve))
}

func (ePriv ecdsaPrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, ePriv.priv, hash[:])
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ecdsaSignature{r, s})
}

type ecdsaPublicKey struct {
	pub *ecdsa.PublicKey
}

func (ePub ecdsaPublicKey) Bytes() []byte {
	return elliptic.Marshal(ePub.pub.Curve, ePub.pub.X, ePub.pub.Y)
}

func (ePub ecdsaPublicKey) Marshal() (VersionnedKey, error) {
	m, err := x509.MarshalPKIXPublicKey(ePub.pub)
	if err != nil {
		return nil, err
	}
	return versionKey(ePub.curve(), m), nil
}

func (ePub ecdsaPublicKey) curve() Curve {
	return getEcdsaCurve(ePub.pub.Curve)
}

func (ePub ecdsaPublicKey) Encrypt(data []byte) (Cipher, error) {
	return eciesEncrypt(data, ePub, ecdsaGenerateShared)
}

func (ePub ecdsaPublicKey) Verify(data []byte, sig []byte) bool {
	var ecdsaSig ecdsaSignature
	if _, err := asn1.Unmarshal(sig, &ecdsaSig); err != nil {
		return false
	}

	hash := sha256.Sum256(data)
	return ecdsa.Verify(ePub.pub, hash[:], ecdsaSig.R, ecdsaSig.S)
}

func (ePub ecdsaPublicKey) Equals(otherPub PublicKey) bool {
	otherBytes := otherPub.Bytes()
	pubBytes := ePub.Bytes()
	return bytes.Equal(otherBytes, pubBytes)
}

func generateECDSAKeys(src io.Reader, c elliptic.Curve) (PrivateKey, PublicKey, error) {
	priv, err := ecdsa.GenerateKey(c, src)
	if err != nil {
		return nil, nil, err
	}
	return ecdsaPrivateKey{priv}, ecdsaPublicKey{&priv.PublicKey}, nil
}

func ecdsaGenerateShared(pub PublicKey, pv PrivateKey) ([]byte, error) {
	ePub := pub.(ecdsaPublicKey)
	ePv := pv.(ecdsaPrivateKey)
	s, _ := ePub.pub.Curve.ScalarMult(ePub.pub.X, ePub.pub.Y, ePv.priv.D.Bytes())
	return s.Bytes(), nil
}

func ecdsaExtractRandomPublicKey(c elliptic.Curve) func([]byte) (PublicKey, int, error) {
	return func(cipher []byte) (PublicKey, int, error) {
		var rLen int

		//Find the public key length
		switch cipher[0] {
		case 2, 3, 4:
			rLen = ((c.Params().BitSize + 7) / 4)
		default:
			return nil, 0, errors.New("invalid public key")
		}

		//Reshape the random public key
		R := new(ecdsa.PublicKey)
		R.Curve = c
		R.X, R.Y = elliptic.Unmarshal(c, cipher[:rLen])
		if R.X == nil || R.Y == nil {
			return nil, 0, errors.New("invalid public key")
		}
		if !R.Curve.IsOnCurve(R.X, R.Y) {
			return nil, 0, errors.New("invalid curve")
		}

		return ecdsaPublicKey{R}, rLen, nil
	}

}

func getEcdsaCurve(c elliptic.Curve) Curve {
	switch c.Params().Name {
	case "P-256":
		return P256Curve
	}
	return -1
}
