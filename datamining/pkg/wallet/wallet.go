package wallet

import (
	"hash"
	"time"

	robot "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validation"
)

//Database provides access to the local repository
type Database interface {
	GetEncWalletAddr(bh BioHash) ([]byte, error)
	GetEncWallet(addr []byte) (CipherWallet, error)
}

//CipherWallet describe the encrypted wallet
type CipherWallet []byte

//Wallet describe a secure Wallet
type Wallet struct {
	walletAddr   []byte
	cWallet      CipherWallet
	timeStamp    time.Time
	emPubk       robot.PublicKey
	emSig        robot.Signature
	biodPubk     robot.PublicKey
	biodSig      robot.Signature
	oldTxnHash   hash.Hash64
	txnHash      hash.Hash64
	masterRobotv validation.MasterRobotValidation
	robotsv      []validation.RobotValidation
}

//CWallet returns the encrypted wallet
func (w Wallet) CWallet() CipherWallet {
	return w.cWallet
}

//WalletAddr returns address of the encrypted wallet
func (w Wallet) WalletAddr() []byte {
	return w.walletAddr
}
