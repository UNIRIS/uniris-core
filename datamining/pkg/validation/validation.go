package validation

import (
	"time"

	robot "github.com/uniris/uniris-core/datamining/pkg"
)

//MasterRobotValidation describe a validation of an elected master robot
type MasterRobotValidation struct {
	lastTxnv []robot.PublicKey
	pow      robot.PublicKey
	rbv      RobotValidation
}

//RobotValidation describe a validation of a robot
type RobotValidation struct {
	status    string
	timestamp time.Time
	pubk      robot.PublicKey
	sig       robot.Signature
}
