package datamining

//Endorsement represents a validation
type Endorsement interface {
	TransactionHash() string
	LastTransactionHash() string
	MasterValidation() MasterValidation
	Validations() []Validation
}

type endorsement struct {
	lastTxHash       string
	txHash           string
	masterValidation MasterValidation
	validations      []Validation
}

//NewEndorsement creates a new endorsement
func NewEndorsement(lastTxHash, txHash string, masterV MasterValidation, valids []Validation) Endorsement {
	return endorsement{
		lastTxHash:       lastTxHash,
		txHash:           txHash,
		masterValidation: masterV,
		validations:      valids,
	}
}

//LastTransactionHash returns the previous transaction hash
func (e endorsement) LastTransactionHash() string {
	return e.lastTxHash
}

//TransactionHash returns the endorsment's transaction hash
func (e endorsement) TransactionHash() string {
	return e.txHash
}

//MasterValidation returns the endorsment's master validation
func (e endorsement) MasterValidation() MasterValidation {
	return e.masterValidation
}

//Validations returns the endorsment's validations
func (e endorsement) Validations() []Validation {
	return e.validations
}
