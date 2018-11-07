package account

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//KeyChainData describe a decoded keychain data
type KeyChainData struct {
	WalletAddr      string
	CipherAddrRobot string
	CipherWallet    string
	PersonPubk      string
	BiodPubk        string
	Sigs            Signatures
}

//Keychain describes the keychain data with its endorsement
type keychain struct {
	data        *KeyChainData
	endorsement datamining.Endorsement
}

//Keychain represents a keychain for an account
type Keychain interface {
	CipherWallet() string
	WalletAddr() string
	CipherAddrRobot() string
	PersonPublicKey() string
	BiodPublicKey() string
	Signatures() Signatures
	Endorsement() datamining.Endorsement
}

//NewKeychain creates a new keychain
func NewKeychain(data *KeyChainData, endor datamining.Endorsement) Keychain {
	return keychain{data, endor}
}

//Endorsement returns the wallet endorsement
func (k keychain) Endorsement() datamining.Endorsement {
	return k.endorsement
}

//CipherWallet returns the encrypted wallet
func (k keychain) CipherWallet() string {
	return k.data.CipherWallet
}

//WalletAddr returns address of the encrypted wallet
func (k keychain) WalletAddr() string {
	return k.data.WalletAddr
}

//CipherAddrRobot get the wallet address encrypted for roboto
func (k keychain) CipherAddrRobot() string {
	return k.data.CipherAddrRobot
}

//PersonPublicKey returns the wallet personal public key
func (k keychain) PersonPublicKey() string {
	return k.data.PersonPubk
}

//BiodPublicKey returns the wallet biod public key
func (k keychain) BiodPublicKey() string {
	return k.data.BiodPubk
}

//Signatures return the wallet signatures
func (k keychain) Signatures() Signatures {
	return k.data.Sigs
}
