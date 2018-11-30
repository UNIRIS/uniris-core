package listing

//SharedKeys describes the shared keys
type SharedKeys interface {

	//RobotPublicKey returns the shared robot public key
	RobotPublicKey() string

	//RobotPrivateKey  returns the shared robot private key (used only on memory | no sent to the client)
	RobotPrivateKey() string

	//EmitterKeyPairs returns the list of shared emitter keys
	EmitterKeyPairs() []SharedKeyPair

	//RequestPublicKey returns the public key for emitter request
	RequestPublicKey() string
}

type sharedKeys struct {
	rPubKey string
	rPvKey  string
	emKP    []SharedKeyPair
}

//NewSharedKeys creates a new shared keys list
func NewSharedKeys(rPv, rPub string, emKP []SharedKeyPair) SharedKeys {
	return sharedKeys{
		rPubKey: rPub,
		rPvKey:  rPv,
		emKP:    emKP,
	}
}

func (sk sharedKeys) RobotPublicKey() string {
	return sk.rPubKey
}

func (sk sharedKeys) RobotPrivateKey() string {
	return sk.rPvKey
}

func (sk sharedKeys) EmitterKeyPairs() []SharedKeyPair {
	return sk.emKP
}

func (sk sharedKeys) RequestPublicKey() string {
	return sk.emKP[0].PublicKey()
}

//SharedKeyPair represent a shared keypair
type SharedKeyPair interface {
	EncryptedPrivateKey() string
	PublicKey() string
}

type sharedKeyPair struct {
	encPvKey string
	pubKey   string
}

//NewSharedKeyPair create a new shared keypair
func NewSharedKeyPair(encPv, pub string) SharedKeyPair {
	return sharedKeyPair{
		encPvKey: encPv,
		pubKey:   pub,
	}
}

func (kp sharedKeyPair) EncryptedPrivateKey() string {
	return kp.encPvKey
}

func (kp sharedKeyPair) PublicKey() string {
	return kp.pubKey
}
