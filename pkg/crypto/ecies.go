package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
)

type extractPubKeyFunc func(cipherData []byte) (PublicKey, int, error)

//Cipher identifies an encoded cipher message with ECIES
type Cipher []byte

func (c Cipher) decode(extPub extractPubKeyFunc) (PublicKey, []byte, []byte, error) {
	rPub, rPubEnd, err := extPub(c)
	if err != nil {
		return nil, nil, nil, err
	}

	hLen := DefaultHashAlgo.Size()
	encMsgEnd := len(c) - hLen
	encMsg := c[rPubEnd:encMsgEnd]
	tag := c[encMsgEnd:]

	return rPub, encMsg, tag, nil
}

func newEncodedCipher(rPub PublicKey, em []byte, tag []byte) Cipher {
	rPubBytes := rPub.bytes()
	out := make(Cipher, len(rPubBytes)+len(em)+len(tag))
	copy(out, rPubBytes)
	copy(out[len(rPubBytes):], em)
	copy(out[len(rPubBytes)+len(em):], tag)
	return out
}

func eciesEncrypt(data []byte, pubKey PublicKey, sharedKey generateSharedFunc) (Cipher, error) {

	rPv, rPub, err := GenerateECKeyPair(pubKey.curve(), rand.Reader)
	if err != nil {
		return nil, err
	}

	secret, err := sharedKey(pubKey, rPv)
	if err != nil {
		return nil, err
	}

	kdfKeys, err := derivateKeys(secret, 2)
	if err != nil {
		return nil, err
	}

	em, err := AESEncrypt(kdfKeys[0], data)
	if err != nil {
		return nil, err
	}

	tag := authenticateMessage(kdfKeys[1], em)

	return newEncodedCipher(rPub, em, tag), nil
}

func eciesDecrypt(c Cipher, pvKey PrivateKey, sharedKey generateSharedFunc, extPub extractPubKeyFunc) ([]byte, error) {
	if len(c) == 0 {
		return nil, errors.New("invalid message")
	}

	rPub, encMsg, tag, err := c.decode(extPub)
	if err != nil {
		return nil, err
	}

	secret, err := sharedKey(rPub, pvKey)
	if err != nil {
		return nil, err
	}

	kdfKeys, err := derivateKeys(secret, 2)
	if err != nil {
		return nil, err
	}

	ackTag := authenticateMessage(kdfKeys[1], encMsg)
	if subtle.ConstantTimeCompare(ackTag, tag) != 1 {
		return nil, errors.New("invalid message")
	}

	data, err := AESDecrypt(kdfKeys[0], encMsg)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func derivateKeys(secret []byte, nbKeys int) (keys [][]byte, err error) {
	hash := DefaultHashAlgo.New
	hkdf := hkdf.New(hash, secret, nil, nil)

	for i := 0; i < nbKeys; i++ {
		key := make([]byte, 16)
		if _, err := io.ReadFull(hkdf, key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return
}

func authenticateMessage(key []byte, msg []byte) []byte {
	mac := hmac.New(DefaultHashAlgo.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}
