package wallet

import (
	"hash"
	"time"

	common "github.com/uniris/uniris-core/datamining/pkg"
	validation "github.com/uniris/uniris-core/datamining/pkg"
)

//BioHash is the hash describing a biometric identitie
type BioHash hash.Hash64

type BioWallet struct {
	bHash           BioHash
	cipherAddrRobot []byte
	cipherAddrBio   []byte
	timeStamp       time.Time
	emPubk          common.PublicKey
	emSig           common.Signature
	biodPubk        common.PublicKey
	biodSig         common.Signature
	txnHash         hash.Hash64
	masterRobotv    validation.MasterRobotValidation
	robotsv         []validation.RobotValidation
}
