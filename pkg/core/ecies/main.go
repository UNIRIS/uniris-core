package main

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"io"
	"os"
	"path/filepath"
	"plugin"

	"golang.org/x/crypto/hkdf"
)

var hashPlugin = filepath.Join(os.Getenv("PLUGINS_DIR"), "hash/plugin.so")
var aesPlugin = filepath.Join(os.Getenv("PLUGINS_DIR"), "aes/plugin.so")
var keyPlugin = filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so")

func curvePlugins() (map[int]string, error) {
	p, err := plugin.Open(keyPlugin)
	if err != nil {
		return nil, err
	}

	p256Sym, err := p.Lookup("P256Curve")
	if err != nil {
		return nil, err
	}
	ed25519Sym, err := p.Lookup("Ed25519Curve")
	if err != nil {
		return nil, err
	}

	return map[int]string{
		*p256Sym.(*int):    filepath.Join(os.Getenv("PLUGINS_DIR"), "ecdsa/plugin.so"),
		*ed25519Sym.(*int): filepath.Join(os.Getenv("PLUGINS_DIR"), "ed25519/plugin.so"),
	}, nil
}

type key interface {
	Bytes() []byte
	Curve() int
}

//Encrypt a message using ECIES to a given public key
//It generates a ephemer keypair and a common shared key based on the epehmer key and the public key.
//The shared key is used to derive two keys: one for AES encryption and other to an HMAC authentication
//The data is finally encrypted with the AES derivate key and authenticated with the HMAC key
//Finally the cipher message is composed from the [Ephemer public key][Encrypted data][HMAC Authentication Tag]
func Encrypt(data []byte, pubK interface{}) ([]byte, error) {

	pubKey, ok := pubK.(key)
	if !ok {
		return nil, errors.New("invalid key")
	}

	hashPlugin, err := plugin.Open(hashPlugin)
	if err != nil {
		return nil, err
	}

	defHashSym, err := hashPlugin.Lookup("DefaultHashAlgo")
	if err != nil {
		return nil, err
	}
	defHashAlgo := defHashSym.(*crypto.Hash)

	aesPlugin, err := plugin.Open(aesPlugin)
	if err != nil {
		return nil, err
	}

	aesEncryptSym, err := aesPlugin.Lookup("Encrypt")
	if err != nil {
		return nil, err
	}
	aesEncrypt := aesEncryptSym.(func(key []byte, msg []byte) ([]byte, error))

	keyPlugin, err := plugin.Open(keyPlugin)
	if err != nil {
		return nil, err
	}

	keyGenerateSym, err := keyPlugin.Lookup("GenerateKeys")
	if err != nil {
		return nil, err
	}
	keyGenerate := keyGenerateSym.(func(curve int, src io.Reader) (interface{}, interface{}, error))
	rPv, rPub, err := keyGenerate(pubKey.Curve(), rand.Reader)
	if err != nil {
		return nil, err
	}

	curves, err := curvePlugins()
	if err != nil {
		return nil, err
	}

	pluginPath, exist := curves[pubKey.Curve()]
	if !exist {
		return nil, errors.New("unsupported curve")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	genSharedSecretSym, err := p.Lookup("GenerateSharedSecret")
	if err != nil {
		return nil, err
	}

	genSharedSecret := genSharedSecretSym.(func(pub []byte, pv []byte) ([]byte, error))
	secret, err := genSharedSecret(pubKey.Bytes(), rPv.(key).Bytes())
	if err != nil {
		return nil, err
	}

	kdfKeys, err := derivateKeys(*defHashAlgo, secret, 2)
	if err != nil {
		return nil, err
	}

	em, err := aesEncrypt(kdfKeys[0], data)
	if err != nil {
		return nil, err
	}

	tag := authenticateMessage(*defHashAlgo, kdfKeys[1], em)

	return newEncodedCipher(rPub.(key), em, tag), nil
}

//Decrypt a message using ECIES to a given private key
//It retrieve the random epehmer public key from the message and re-generate a common shared key
//The shared key is used to derive two keys: one for AES decryption and other to an HMAC authentication
//An authentication check if done by the HMAC authentication key
//And finally the the data is decrypted using the AES key
func Decrypt(cipher []byte, pvK interface{}) ([]byte, error) {

	pvKey, ok := pvK.(key)
	if !ok {
		return nil, errors.New("invalid key")
	}

	if len(cipher) == 0 {
		return nil, errors.New("invalid message")
	}

	hashPlugin, err := plugin.Open(hashPlugin)
	if err != nil {
		return nil, err
	}

	defHashSym, err := hashPlugin.Lookup("DefaultHashAlgo")
	if err != nil {
		return nil, err
	}
	defHashAlgo := defHashSym.(*crypto.Hash)

	aesPlugin, err := plugin.Open(aesPlugin)
	if err != nil {
		return nil, err
	}

	aesDecryptSym, err := aesPlugin.Lookup("Decrypt")
	if err != nil {
		return nil, err
	}
	aesDecrypt := aesDecryptSym.(func(key []byte, msg []byte) ([]byte, error))

	rPub, encMsg, tag, err := decodeCipher(cipher, *defHashAlgo, pvKey.Curve())
	if err != nil {
		return nil, err
	}

	curves, err := curvePlugins()
	if err != nil {
		return nil, err
	}

	pluginPath, exist := curves[pvKey.Curve()]
	if !exist {
		return nil, errors.New("unsupported curve")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}

	genSharedSecretSym, err := p.Lookup("GenerateSharedSecret")
	if err != nil {
		return nil, err
	}

	genSharedSecret := genSharedSecretSym.(func(pub []byte, pv []byte) ([]byte, error))
	secret, err := genSharedSecret(rPub, pvKey.Bytes())
	if err != nil {
		return nil, err
	}

	kdfKeys, err := derivateKeys(*defHashAlgo, secret, 2)
	if err != nil {
		return nil, err
	}

	ackTag := authenticateMessage(*defHashAlgo, kdfKeys[1], encMsg)
	if subtle.ConstantTimeCompare(ackTag, tag) != 1 {
		return nil, errors.New("invalid message")
	}

	data, err := aesDecrypt(kdfKeys[0], encMsg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func derivateKeys(hashAlgo crypto.Hash, secret []byte, nbKeys int) (keys [][]byte, err error) {

	hkdf := hkdf.New(hashAlgo.New, secret, nil, nil)

	for i := 0; i < nbKeys; i++ {
		key := make([]byte, 16)
		if _, err := io.ReadFull(hkdf, key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return
}

func authenticateMessage(hashAlgo crypto.Hash, key []byte, msg []byte) []byte {
	mac := hmac.New(hashAlgo.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func decodeCipher(cipher []byte, hashAlgo crypto.Hash, curve int) (randKey []byte, encMsg []byte, tag []byte, err error) {

	curves, err := curvePlugins()
	if err != nil {
		return nil, nil, nil, err
	}

	var randKeyLength int

	pluginPath, exist := curves[curve]
	if !exist {
		return nil, nil, nil, errors.New("unsupported curve")
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, nil, nil, err
	}
	sym, err := p.Lookup("ExtractMessagePublicKey")
	if err != nil {
		return nil, nil, nil, err
	}
	f := sym.(func([]byte) ([]byte, int, error))
	randKey, randKeyLength, err = f(cipher)
	if err != nil {
		return nil, nil, nil, err
	}

	hLen := hashAlgo.Size()
	encMsgEnd := len(cipher) - hLen
	encMsg = cipher[randKeyLength:encMsgEnd]
	tag = cipher[encMsgEnd:]

	return
}

func newEncodedCipher(randKey key, em []byte, tag []byte) []byte {
	randKeyBytes := randKey.Bytes()
	out := make([]byte, len(randKeyBytes)+len(em)+len(tag))
	copy(out, randKeyBytes)
	copy(out[len(randKeyBytes):], em)
	copy(out[len(randKeyBytes)+len(em):], tag)
	return out
}
