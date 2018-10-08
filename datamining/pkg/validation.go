package validation

import (
	"time"

	common "github.com/uniris/uniris-core/datamining/pkg"
	wallet "github.com/uniris/uniris-core/datamining/pkg"
)

//MasterRobotValidation describe a validation of an elected master robot
type MasterRobotValidation struct {
	lastTxnv []common.PublicKey
	pow      common.PublicKey
	rbv      RobotValidation
}

//RobotValidation describe a validation of a robot
type RobotValidation struct {
	status    string
	timestamp time.Time
	pubk      common.PublicKey
	sig       common.Signature
}
