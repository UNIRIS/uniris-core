package robot

import ()

//PrivateKey describe a Private Key
type PrivateKey []byte

//Signature describe a digital signature
type Signature []byte

//PublicKey describe a Public key
type PublicKey []byte

//CipherWallet describe the encrypted wallet
type CipherWallet []byte

//BioHash is the hash describing a biometric identitie
type BioHash []byte

//ECDSAKeyPair represent ECDSA key pair
type ECDSAKeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

//NewECDSAPair creates a new pair
func NewECDSAPair(pv []byte, pb []byte) ECDSAKeyPair {
	return ECDSAKeyPair{pv, pb}
}

//KeyReader describes methods to read in keys
type KeyReader interface {
	SharedRobotPrivateKey() (PrivateKey, error)
	SharedRobotPublicKey() (PublicKey, error)
	SharedBiodPublicKey() (PublicKey, error)
}

//Encrypter describes methods to encrypt or decrypt data
type Encrypter interface {
	Decrypt(PrivateKey, data []byte) ([]byte, error)
	Encrypt(PublicKey, data []byte) ([]byte, error)
}

//Signer describes methods to sign and verify data
type Signer interface {
	Verify(PublicKey, data []byte, hash []byte) error
	Sign(PrivateKey, data []byte) ([]byte, error)
}
