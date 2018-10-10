package wallet

import (
	"hash"
	"time"

	robot "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validation"
)

//BioHash is the hash describing a biometric identitie
type BioHash hash.Hash64

//BioWallet describe the data adressing biometric imprint and the wallet address
type BioWallet struct {
	bHash           BioHash
	cipherAddrRobot []byte
	cipherAddrBio   []byte
	timeStamp       time.Time
	emPubk          robot.PublicKey
	emSig           robot.Signature
	biodPubk        robot.PublicKey
	biodSig         robot.Signature
	txnHash         hash.Hash64
	masterRobotv    validation.MasterRobotValidation
	robotsv         []validation.RobotValidation
}

//Bhash returns the biometric hash
func (b BioWallet) Bhash() BioHash {
	return b.bHash
}

//CipherAddrRobot returns the address of the wallet encrypted with shared robot publickey
func (b BioWallet) CipherAddrRobot() []byte {
	return b.cipherAddrRobot
}
