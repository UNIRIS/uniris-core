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
	Sigs            datamining.Signatures
}

//Keychain describes the keychain data with its endorsement
type keychain struct {
	data        *KeyChainData
	oldTxnHash  string
	endorsement datamining.Endorsement
}

//Keychain represents a keychain for an account
type Keychain interface {
	Endorsement() datamining.Endorsement
	CipherWallet() string
	WalletAddr() string
	CipherAddrRobot() string
	PersonPublicKey() string
	BiodPublicKey() string
	Signatures() datamining.Signatures
	OldTransactionHash() string
}

//NewKeychain creates a new keychain
func NewKeychain(data *KeyChainData, endor datamining.Endorsement, oldTxnHash string) Keychain {
	return keychain{data, oldTxnHash, endor}
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
func (k keychain) Signatures() datamining.Signatures {
	return k.data.Sigs
}

//OldTransactionHash returns the hash of the previous transaction
func (k keychain) OldTransactionHash() string {
	return k.oldTxnHash
}
