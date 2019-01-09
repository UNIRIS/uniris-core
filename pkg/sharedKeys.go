package uniris

//SharedKeys describe shared keypair
type SharedKeys struct {
	encPvKey string
	pubKey   string
}

//NewSharedKeyPair creates a new proposed keypair
func NewSharedKeyPair(encPvKey, pubKey string) SharedKeys {
	return SharedKeys{encPvKey, pubKey}
}

//PublicKey returns the public key for the proposed keypair
func (sK SharedKeys) PublicKey() string {
	return sK.pubKey
}

//EncryptedPrivateKey returns the encrypted private key for the proposed keypair
func (sK SharedKeys) EncryptedPrivateKey() string {
	return sK.encPvKey
}
