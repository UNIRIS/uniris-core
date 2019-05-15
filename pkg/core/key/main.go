package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"

	"golang.org/x/crypto/ed25519"
)

var (

	//Ed25519Curve identifies the Ed25519 elliptic curve
	Ed25519Curve int

	//P256Curve identifies the NIST-P256 elliptic curve
	P256Curve = 1
)

//ECKey represents a key using elliptic curve
type ECKey interface {
	Curve() int
	Bytes() []byte
	Marshal() []byte
}

type key struct {
	bytes []byte
	curve int
}

func (k key) Bytes() []byte {
	return k.bytes
}

func (k key) Curve() int {
	return k.curve
}

func (k key) Marshal() []byte {
	out := make([]byte, 1+len(k.Bytes()))
	out[0] = byte(int(k.Curve()))
	copy(out[1:], k.Bytes())
	return out
}

//PrivateKey represents an EC private key
type PrivateKey interface {
	Sign(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	ECKey
}

type pvKey struct {
	ECKey
}

func (pv pvKey) Sign(data []byte) ([]byte, error) {

	var pluginName string
	switch pv.Curve() {
	case Ed25519Curve:
		pluginName = "ed25519"
		break
	case P256Curve:
		pluginName = "ecdsa"
		break
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), fmt.Sprintf("%s/plugin.so", pluginName)))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("Sign")
	if err != nil {
		return nil, err
	}

	f := sym.(func(key []byte, data []byte) ([]byte, error))
	return f(pv.ECKey.Bytes(), data)
}

func (pv pvKey) Decrypt(data []byte) ([]byte, error) {
	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "ecies/plugin.so"))
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("Decrypt")
	if err != nil {
		return nil, err
	}

	f := sym.(func(data []byte, pvK interface{}) ([]byte, error))
	return f(data, pv)
}

//PublicKey represents an EC public key
type PublicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Encrypt(data []byte) ([]byte, error)
	ECKey
}

type pubKey struct {
	ECKey
}

func (pub pubKey) Verify(data []byte, sig []byte) (bool, error) {

	var pluginName string
	switch pub.Curve() {
	case Ed25519Curve:
		pluginName = "ed25519"
		break
	case P256Curve:
		pluginName = "ecdsa"
		break
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), fmt.Sprintf("%s/plugin.so", pluginName)))
	if err != nil {
		return false, err
	}

	sym, err := p.Lookup("Verify")
	if err != nil {
		return false, err
	}

	f := sym.(func(key []byte, data []byte, sig []byte) (bool, error))
	return f(pub.ECKey.Bytes(), data, sig)
}

func (pub pubKey) Encrypt(data []byte) ([]byte, error) {
	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "ecies/plugin.so"))
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("Encrypt")
	if err != nil {
		return nil, err
	}

	f := sym.(func(data []byte, pubK interface{}) ([]byte, error))
	return f(data, pub)
}

//GenerateKeys creates a new elliptic keypair from a given specific curve
func GenerateKeys(curve int, src io.Reader) (interface{}, interface{}, error) {
	switch curve {
	case P256Curve:

		p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "ecdsa/plugin.so"))
		if err != nil {
			return nil, nil, err
		}

		sym, err := p.Lookup("GenerateKeys")
		if err != nil {
			return nil, nil, err
		}

		f := sym.(func(io.Reader, elliptic.Curve) ([]byte, []byte, error))
		pv, pub, err := f(src, elliptic.P256())
		if err != nil {
			return nil, nil, err
		}
		return pvKey{
				ECKey: key{
					bytes: pv,
					curve: P256Curve,
				},
			}, pubKey{
				ECKey: key{
					bytes: pub,
					curve: P256Curve,
				},
			}, nil
	case Ed25519Curve:

		p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "ed25519/plugin.so"))
		if err != nil {
			return nil, nil, err
		}

		sym, err := p.Lookup("GenerateKeys")
		if err != nil {
			return nil, nil, err
		}

		f := sym.(func(io.Reader) ([]byte, []byte, error))
		pv, pub, err := f(src)
		if err != nil {
			return nil, nil, err
		}
		return pvKey{
				ECKey: key{
					bytes: pv,
					curve: Ed25519Curve,
				},
			}, pubKey{
				ECKey: key{
					bytes: pub,
					curve: Ed25519Curve,
				},
			}, nil
	default:
		return nil, nil, errors.New("unsupported EC curve")
	}
}

//ParsePublicKey converts a marshalled versionned public key into a PublicKey
func ParsePublicKey(k []byte) (interface{}, error) {
	if len(k) == 0 {
		return nil, errors.New("cannot parse an empty public key")
	}
	switch int(k[0]) {
	case P256Curve:
		pub, err := x509.ParsePKIXPublicKey(k[1:])
		if err != nil {
			return nil, err
		}
		if _, ok := pub.(*ecdsa.PublicKey); !ok {
			return nil, errors.New("invalid ECDSA public key")
		}

		return pubKey{
			ECKey: key{
				bytes: k[1:],
				curve: P256Curve,
			},
		}, nil
	case Ed25519Curve:
		if len(k[1:]) != ed25519.PublicKeySize {
			return nil, errors.New("invalid ed25519 private key")
		}
		return pubKey{
			ECKey: key{
				bytes: k[1:],
				curve: Ed25519Curve,
			},
		}, nil
	default:
		return nil, errors.New("unsupported curve")
	}
}

//ParsePrivateKey converts a marshalled versionned private key into a Private key
func ParsePrivateKey(k []byte) (interface{}, error) {
	if len(k) == 0 {
		return nil, errors.New("cannot parse an empty private key")
	}
	switch int(k[0]) {
	case P256Curve:
		_, err := x509.ParseECPrivateKey(k[1:])
		if err != nil {
			return nil, err
		}
		return pvKey{
			ECKey: key{
				bytes: k[1:],
				curve: P256Curve,
			},
		}, nil
	case Ed25519Curve:
		if len(k[1:]) != ed25519.PrivateKeySize {
			return nil, errors.New("invalid ed25519 private key")
		}
		return pvKey{
			ECKey: key{
				bytes: k[1:],
				curve: Ed25519Curve,
			},
		}, nil
	default:
		return nil, errors.New("unsupported curve")
	}
}
