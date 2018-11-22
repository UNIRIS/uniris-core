package account

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//ID describe ID data
type ID interface {

	//Hash returns the ID hash
	Hash() string

	//EncryptedAddrByRobot returns the account's address encrypted with shared robot publickey
	EncryptedAddrByRobot() string

	//EncryptedAddrByID returns the account's address encrypted with ID public key
	EncryptedAddrByID() string

	//EncryptedAESKey returns the AES key encrypted with ID public key
	EncryptedAESKey() string

	//PublicKey returns ID public key
	PublicKey() string

	//IDSignature returns the signature provided by the ID
	IDSignature() string

	//EmitterSignature returns the signature provided by the emitter's device
	EmitterSignature() string

	//Proposal returns the proposal for this transaction
	Proposal() datamining.Proposal
}

type id struct {
	hash         string
	encAddrRobot string
	encAddrID    string
	encAESKey    string
	pubk         string
	idSig        string
	emSig        string
	prop         datamining.Proposal
}

//NewID create new ID
func NewID(hash, encAddrRobot, encAddrID, encAESKey, pubk, idSig, emSig string, prop datamining.Proposal) ID {
	return id{
		hash:         hash,
		encAddrRobot: encAddrRobot,
		encAddrID:    encAddrID,
		encAESKey:    encAESKey,
		pubk:         pubk,
		idSig:        idSig,
		emSig:        emSig,
		prop:         prop,
	}
}

func (id id) Hash() string {
	return id.hash
}

func (id id) EncryptedAddrByRobot() string {
	return id.encAddrRobot
}

func (id id) EncryptedAddrByID() string {
	return id.encAddrID
}

func (id id) EncryptedAESKey() string {
	return id.encAESKey
}

func (id id) PublicKey() string {
	return id.pubk
}

func (id id) IDSignature() string {
	return id.idSig
}

func (id id) EmitterSignature() string {
	return id.emSig
}

func (id id) Proposal() datamining.Proposal {
	return id.prop
}

//EndorsedID aggregates ID and its endorsement
type EndorsedID interface {
	ID

	//Endorsement returns the id's endorsement
	Endorsement() mining.Endorsement
}

type endorsedID struct {
	id          ID
	endorsement mining.Endorsement
}

//NewEndorsedID creates a new id's endorsed
func NewEndorsedID(id ID, endor mining.Endorsement) EndorsedID {
	return endorsedID{id, endor}
}

func (eID endorsedID) Hash() string {
	return eID.id.Hash()
}

func (eID endorsedID) EncryptedAddrByRobot() string {
	return eID.id.EncryptedAddrByRobot()
}

func (eID endorsedID) EncryptedAddrByID() string {
	return eID.id.EncryptedAddrByID()
}

func (eID endorsedID) EncryptedAESKey() string {
	return eID.id.EncryptedAESKey()
}

func (eID endorsedID) PublicKey() string {
	return eID.id.PublicKey()
}

func (eID endorsedID) IDSignature() string {
	return eID.id.IDSignature()
}

func (eID endorsedID) EmitterSignature() string {
	return eID.id.EmitterSignature()
}

func (eID endorsedID) Proposal() datamining.Proposal {
	return eID.id.Proposal()
}

func (eID endorsedID) Endorsement() mining.Endorsement {
	return eID.endorsement
}
