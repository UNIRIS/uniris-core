package crypto

import (
	"errors"
	"io"

	"github.com/agl/ed25519/extra25519"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

type ed25519PrivateKey struct {
	priv ed25519.PrivateKey
}

func (edPriv ed25519PrivateKey) bytes() []byte {
	return edPriv.priv
}

func (edPriv ed25519PrivateKey) Marshal() (VersionnedKey, error) {
	return versionKey(Ed25519Curve, edPriv.priv), nil
}

func (edPriv ed25519PrivateKey) curve() Curve {
	return Ed25519Curve
}

func (edPriv ed25519PrivateKey) Sign(data []byte) (sig []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()
	sig = ed25519.Sign(edPriv.priv, data)
	return
}

func (edPriv ed25519PrivateKey) Decrypt(c Cipher) ([]byte, error) {
	return eciesDecrypt(c, edPriv, ed25519GenerateShared, ed25519ExtractRandomPublicKey)
}

type ed25519PublicKey struct {
	pub ed25519.PublicKey
}

func (edPub ed25519PublicKey) bytes() []byte {
	return edPub.pub
}

func (edPub ed25519PublicKey) Marshal() (VersionnedKey, error) {
	return versionKey(Ed25519Curve, edPub.pub), nil
}

func (edPub ed25519PublicKey) curve() Curve {
	return Ed25519Curve
}

func (edPub ed25519PublicKey) Verify(data []byte, sig []byte) (valid bool) {
	defer func() {
		if r := recover(); r != nil {
			valid = false
			return
		}
	}()

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	return ed25519.Verify(edPub.pub, data, sig)
}

func (edPub ed25519PublicKey) Encrypt(data []byte) (Cipher, error) {
	return eciesEncrypt(data, edPub, ed25519GenerateShared)
}

func generateEd25519Keys(src io.Reader) (PrivateKey, PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}
	return ed25519PrivateKey{priv}, ed25519PublicKey{pub}, nil
}

func ed25519GenerateShared(pub PublicKey, priv PrivateKey) ([]byte, error) {
	ePub := pub.(ed25519PublicKey)
	curve25519Public, err := convertPublicEd25519ToCurve25519(ePub.pub)
	if err != nil {
		return nil, err
	}

	edpriv := priv.(ed25519PrivateKey)
	curve25519priv := convertPrivateEd25519ToCurve25519(edpriv.priv)

	var s [ed25519.PublicKeySize]byte
	curve25519.ScalarMult(&s, &curve25519priv, &curve25519Public)
	return s[:], nil
}

func ed25519ExtractRandomPublicKey(cipherData []byte) (PublicKey, int, error) {
	if len(cipherData) < ed25519.PublicKeySize {
		return nil, 0, errors.New("invalid message")
	}
	pubBytes := cipherData[0:ed25519.PublicKeySize]
	return ed25519PublicKey{pub: pubBytes}, ed25519.PublicKeySize, nil
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
