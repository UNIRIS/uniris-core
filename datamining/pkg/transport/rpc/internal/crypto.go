package internalrpc

import (
	crypto "github.com/uniris/uniris-core/datamining/pkg/crypto"
)

//Decrypter data for the rpc transport layer
type Decrypter struct {
	e        crypto.Encrypter
	robotKey []byte
}

//NewDecrypter creates a decrypted
func NewDecrypter(robotKey []byte) Decrypter {
	return Decrypter{
		e:        crypto.Encrypter{},
		robotKey: robotKey,
	}
}

//Decipher decrypt a data
func (d Decrypter) Decipher(data []byte) ([]byte, error) {
	return d.e.Decrypt(d.robotKey, data)
}
