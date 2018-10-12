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

//LastTxnv ...
func (m MasterRobotValidation) LastTxnv() []robot.PublicKey {
	return m.lastTxnv
}

//Pow ...
func (m MasterRobotValidation) Pow() robot.PublicKey {
	return m.pow
}

//Rbv ...
func (m MasterRobotValidation) Rbv() RobotValidation {
	return m.rbv
}

//RobotValidation describe a validation of a robot
type RobotValidation struct {
	status    string
	timestamp time.Time
	pubk      robot.PublicKey
	sig       robot.Signature
}

//Status ..
func (m RobotValidation) Status() string {
	return m.status
}

//Timestamp ...
func (m RobotValidation) Timestamp() time.Time {
	return m.timestamp
}

//Pubk ...
func (m RobotValidation) Pubk() robot.PublicKey {
	return m.pubk
}

//Sig ...
func (m RobotValidation) Sig() robot.Signature {
	return m.sig
}
