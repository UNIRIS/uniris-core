package mining

//MasterValidation describe a validation of an elected master robot
type MasterValidation interface {

	//LastTransactionMiners returns the miners list of the last transaction
	LastTransactionMiners() []string

	//ProofOfWorkRobotKey returns the public key from the robot which performs the Proof of work
	ProofOfWorkRobotKey() string

	//ProofOfWorkRobotKey returns the validation performed during the Proof of work
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

func (m masterValidation) LastTransactionMiners() []string {
	return m.lastTxRvk
}

func (m masterValidation) ProofOfWorkRobotKey() string {
	return m.powRobotKey
}

func (m masterValidation) ProofOfWorkValidation() Validation {
	return m.powValid
}
