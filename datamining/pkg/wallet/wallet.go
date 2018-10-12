package wallet

import (
	"hash"
	"time"

	robot "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validation"
	formater "github.com/uniris/uniris-core/datamining/pkg/walletformating"
)

//Database provides access to the local repository
type Database interface {
	GetBioWallet(bh robot.BioHash) (BioWallet, error)
	GetWallet(addr []byte) (Wallet, error)
	AddWallet(w Wallet) error
	AddBioWallet(bw BioWallet) error
}

//Wallet describe a secure Wallet
type Wallet struct {
	fw           formater.FormatedWallet
	timeStamp    time.Time
	oldTxnHash   hash.Hash64
	txnHash      hash.Hash64
	masterRobotv validation.MasterRobotValidation
	robotsv      []validation.RobotValidation
}

//CWallet returns the encrypted wallet
func (w Wallet) CWallet() robot.CipherWallet {
	return w.fw.CWallet
}

//WalletAddr returns address of the encrypted wallet
func (w Wallet) WalletAddr() []byte {
	return w.fw.WalletAddr
}

//BioWallet describe the data adressing biometric imprint and the wallet address
type BioWallet struct {
	fbw          formater.FormatedBioWallet
	timeStamp    time.Time
	txnHash      hash.Hash64
	masterRobotv validation.MasterRobotValidation
	robotsv      []validation.RobotValidation
}

//Bhash returns the biometric hash
func (b BioWallet) Bhash() robot.BioHash {
	return b.fbw.BHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b BioWallet) CipherAddrRobot() []byte {
	return b.fbw.CipherAddrRobot
}

//CipherAddrBio returns the address of the wallet encrypted with person keys
func (b BioWallet) CipherAddrBio() []byte {
	return b.fbw.CipherAddrBio
}

//CipherAesKey returns the AES key encrypted with person keys
func (b BioWallet) CipherAesKey() []byte {
	return b.fbw.CipherAesKey
}
