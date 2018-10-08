package wallet

import (
	"hash"
	"time"

	common "github.com/uniris/uniris-core/datamining/pkg"
	validation "github.com/uniris/uniris-core/datamining/pkg"
)

//Repository provides access to the local repository
type database interface {
	GetWalletAddr(BioHash) ([]byte, error)
	GetWallet(addr []byte) (CipherWallet, error)
}

//CipherWallet describe the encrypted wallet
type CipherWallet []byte

//Wallet describe a secure Wallet
type Wallet struct {
	walletAddr   []byte
	cWallet      CipherWallet
	timeStamp    time.Time
	emPubk       common.PublicKey
	emSig        common.Signature
	biodPubk     common.PublicKey
	biodSig      common.Signature
	oldTxnHash   hash.Hash64
	txnHash      hash.Hash64
	masterRobotv validation.MasterRobotValidation
	robotsv      []validation.RobotValidation
}
