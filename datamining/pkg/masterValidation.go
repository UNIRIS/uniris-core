package datamining

//MasterValidation describe a validation of an elected master robot
type MasterValidation interface {
	LastTransactionMiners() []string
	ProofOfWorkRobotKey() string
	ProofOfWorkValidation() Validation
}

type masterValidation struct {
	lastTxRvk   []string
	powRobotKey string
	powValid    Validation
}

//NewMasterValidation creates a new master validation
func NewMasterValidation(lastTxRvk []string, powRobotKey string, powValid Validation) MasterValidation {
	return masterValidation{lastTxRvk, powRobotKey, powValid}
}

//LastTransactionMiners returns the list of public keys which validate the last transaction
func (m masterValidation) LastTransactionMiners() []string {
	return m.lastTxRvk
}

//ProofOfWorkRobotKey returns the public key of the robot which perform the PoW
func (m masterValidation) ProofOfWorkRobotKey() string {
	return m.powRobotKey
}

//ProofOfWorkValidation returns the transaction proceed after the proof of work
func (m masterValidation) ProofOfWorkValidation() Validation {
	return m.powValid
}
