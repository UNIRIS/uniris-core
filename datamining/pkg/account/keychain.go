package account

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//Keychain describe a keychain
type Keychain interface {

	//EncryptedAddrByRobot returns encrypted address by the shared robot key
	EncryptedAddrByRobot() string

	//EncryptedWallet returns encrypted wallet by the person AES key
	EncryptedWallet() string

	//IDPublicKey returns the ID public key
	IDPublicKey() string

	//IDSignature returns the signature provided by the ID
	IDSignature() string

	//EmitterSignature returns the signature provided by the emitter's device
	EmitterSignature() string

	//Proposal returns the proposal for this transaction
	Proposal() datamining.Proposal
}

type keychain struct {
	encAddr   string
	encWallet string
	idPubk    string
	idSig     string
	emSig     string
	prop      datamining.Proposal
}

//NewKeychain creates a new keychain
func NewKeychain(encAddrRobot, encWallet, idPubk, idSig, emSig string, prop datamining.Proposal) Keychain {
	return keychain{
		encAddr:   encAddrRobot,
		encWallet: encWallet,
		idPubk:    idPubk,
		idSig:     idSig,
		emSig:     emSig,
		prop:      prop,
	}
}

func (k keychain) EncryptedAddrByRobot() string {
	return k.encAddr
}

func (k keychain) EncryptedWallet() string {
	return k.encWallet
}

func (k keychain) IDPublicKey() string {
	return k.idPubk
}

func (k keychain) IDSignature() string {
	return k.idSig
}

func (k keychain) EmitterSignature() string {
	return k.emSig
}

func (k keychain) Proposal() datamining.Proposal {
	return k.prop
}

//EndorsedKeychain aggregates keychain and it's endorsement
type EndorsedKeychain interface {
	Keychain

	//Address returns the keychain address
	Address() string

	//Endorsement returns the keychain data endorsement
	Endorsement() mining.Endorsement
}

type endorsedKeychain struct {
	address     string
	k           Keychain
	endorsement mining.Endorsement
}

//NewEndorsedKeychain creates a new keychain endorsed
func NewEndorsedKeychain(address string, k Keychain, endor mining.Endorsement) EndorsedKeychain {
	return endorsedKeychain{address, k, endor}
}

func (eK endorsedKeychain) Address() string {
	return eK.address
}

func (eK endorsedKeychain) EncryptedAddrByRobot() string {
	return eK.k.EncryptedAddrByRobot()
}

func (eK endorsedKeychain) EncryptedWallet() string {
	return eK.k.EncryptedWallet()
}

func (eK endorsedKeychain) IDPublicKey() string {
	return eK.k.IDPublicKey()
}

func (eK endorsedKeychain) IDSignature() string {
	return eK.k.IDSignature()
}

func (eK endorsedKeychain) EmitterSignature() string {
	return eK.k.EmitterSignature()
}

func (eK endorsedKeychain) Proposal() datamining.Proposal {
	return eK.k.Proposal()
}

func (eK endorsedKeychain) Endorsement() mining.Endorsement {
	return eK.endorsement
}
