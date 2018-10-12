package datamining

import (
	"hash"
	"time"
)

//BioHash is the hash describing a biometric identitie
type BioHash []byte

//Hash define
type Hash hash.Hash64

//Timestamp describe timestamp
type Timestamp time.Time

//CipherWallet describe the encrypted wallet
type CipherWallet []byte

//PublicKey describe a Public key
type PublicKey []byte

//DERSignature describe an encoded signature
type DERSignature []byte

//WalletAddr describes the wallet address
type WalletAddr []byte

//Signatures describe differnet needed signatures
type Signatures struct {
	EmSig   DERSignature
	BiodSig DERSignature
}

//WalletData describe a decoded wallet data from the wallet request
type WalletData struct {
	WalletAddr      []byte
	CipherAddrRobot []byte
	CipherWallet    CipherWallet
	EmPubk          PublicKey
	BiodPubk        PublicKey
	Sigs            Signatures
}

//BioData describe a decoded biometric data from the wallet request
type BioData struct {
	BHash           BioHash
	CipherAddrRobot WalletAddr
	CipherAddrBio   WalletAddr
	CipherAESKey    []byte
	EmPubk          PublicKey
	BiodPubk        PublicKey
	Sigs            Signatures
}

//Wallet describes the stored wallet with its endorsement
type Wallet struct {
	data       WalletData
	oldTxnHash Hash
	endorsment Endorsement
}

//NewWallet creates a new wallet
func NewWallet(data WalletData, endor Endorsement, oldTxnHash Hash) Wallet {
	return Wallet{data, oldTxnHash, endor}
}

//CipherWallet returns the encrypted wallet
func (w Wallet) CipherWallet() CipherWallet {
	return w.data.CipherWallet
}

//WalletAddr returns address of the encrypted wallet
func (w Wallet) WalletAddr() WalletAddr {
	return w.data.WalletAddr
}

//BioWallet describes the stored biometric wallet with its endorsement
type BioWallet struct {
	data       BioData
	endorsment Endorsement
}

//NewBioWallet creates a new bio wallet
func NewBioWallet(data BioData, endor Endorsement) BioWallet {
	return BioWallet{data, endor}
}

//Bhash returns the biometric hash
func (b BioWallet) Bhash() BioHash {
	return b.data.BHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b BioWallet) CipherAddrRobot() []byte {
	return b.data.CipherAddrRobot
}

//CipherAddrBio returns the address of the wallet encrypted with person keys
func (b BioWallet) CipherAddrBio() []byte {
	return b.data.CipherAddrBio
}

//CipherAESKey returns the AES key encrypted with person keys
func (b BioWallet) CipherAESKey() []byte {
	return b.data.CipherAESKey
}
