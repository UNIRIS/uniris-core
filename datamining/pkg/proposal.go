package datamining

//Proposal describe a proposal for a transaction
type Proposal interface {

	//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
	SharedEmitterKeyPair() ProposedKeyPair
}

type prop struct {
	sharedEmitterKP ProposedKeyPair
}

//NewProposal create a new proposal for a transaction
func NewProposal(shdEmitterKP ProposedKeyPair) Proposal {
	return prop{
		sharedEmitterKP: shdEmitterKP,
	}
}

func (p prop) SharedEmitterKeyPair() ProposedKeyPair {
	return p.sharedEmitterKP
}

//ProposedKeyPair describe proposed keypair
type ProposedKeyPair interface {

	//PublicKey returns the public key for the proposed keypair
	PublicKey() string

	//EncryptedPrivateKey returns the encrypted private key for the proposed keypair
	EncryptedPrivateKey() string
}

type propKP struct {
	encPvKey string
	pubKey   string
}

//NewProposedKeyPair creates a new proposed keypair
func NewProposedKeyPair(encPvKey, pubKey string) ProposedKeyPair {
	return propKP{encPvKey, pubKey}
}

func (p propKP) PublicKey() string {
	return p.pubKey
}

func (p propKP) EncryptedPrivateKey() string {
	return p.encPvKey
}
