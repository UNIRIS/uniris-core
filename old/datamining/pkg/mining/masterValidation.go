package mining

//MasterValidation describe a validation of an elected master robot
type MasterValidation interface {

	//LastTransactionMiners returns the miners list of the last transaction
	LastTransactionMiners() []string

	//ProofOfWorkKey returns the key that validated the proof of work
	ProofOfWorkKey() string

	//Validation returns the validation for the proof of work
	ProofOfWorkValidation() Validation
}

type masterValidation struct {
	lastTxRvk []string
	powKey    string
	powValid  Validation
}

//NewMasterValidation creates a new master validation
func NewMasterValidation(lastTxRvk []string, powKey string, powValid Validation) MasterValidation {
	return masterValidation{lastTxRvk, powKey, powValid}
}

func (m masterValidation) LastTransactionMiners() []string {
	return m.lastTxRvk
}

func (m masterValidation) ProofOfWorkKey() string {
	return m.powKey
}

func (m masterValidation) ProofOfWorkValidation() Validation {
	return m.powValid
}
