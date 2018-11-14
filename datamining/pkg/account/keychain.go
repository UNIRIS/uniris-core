package account

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//KeychainData describe a keychain data
type KeychainData interface {

	//CipherAddrRobot returns encrypted address by the shared robot key
	CipherAddrRobot() string

	//CipherWallet returns encrypted wallet by the person AES key
	CipherWallet() string

	//PersonPublicKey returns the person public key
	PersonPublicKey() string

	//Signature returns the signatures of the keychain data
	Signatures() Signatures
}

type keyChainData struct {
	cipherAddr   string
	cipherWallet string
	personPubk   string
	sigs         Signatures
}

//NewKeychainData creates a new keychain data
func NewKeychainData(cipherAddr, cipherWallet, personPubk string, sigs Signatures) KeychainData {
	return keyChainData{cipherAddr, cipherWallet, personPubk, sigs}
}

func (k keyChainData) CipherAddrRobot() string {
	return k.cipherAddr
}

func (k keyChainData) CipherWallet() string {
	return k.cipherWallet
}

func (k keyChainData) PersonPublicKey() string {
	return k.personPubk
}

func (k keyChainData) Signatures() Signatures {
	return k.sigs
}

//Keychain aggregates keychain data and it's endorsement
type Keychain interface {
	KeychainData

	//Address returns the keychain address
	Address() string

	//Endorsement returns the keychain data endorsement
	Endorsement() mining.Endorsement
}

type keychain struct {
	address     string
	data        KeychainData
	endorsement mining.Endorsement
}

//NewKeychain creates a new keychain aggregate
func NewKeychain(address string, data KeychainData, endor mining.Endorsement) Keychain {
	return keychain{address, data, endor}
}

func (k keychain) Address() string {
	return k.address
}

func (k keychain) CipherAddrRobot() string {
	return k.data.CipherAddrRobot()
}

func (k keychain) CipherWallet() string {
	return k.data.CipherWallet()
}

func (k keychain) PersonPublicKey() string {
	return k.data.PersonPublicKey()
}

func (k keychain) Signatures() Signatures {
	return k.data.Signatures()
}

func (k keychain) Endorsement() mining.Endorsement {
	return k.endorsement
}
