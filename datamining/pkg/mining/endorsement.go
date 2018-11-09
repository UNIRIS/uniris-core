package mining

//Endorsement represents a validation
type Endorsement interface {
	//TransactionHash returns the endorsment's transaction hash
	TransactionHash() string

	//LastTransactionHash returns the previous transaction hash
	LastTransactionHash() string

	//MasterValidation returns the endorsment's master validation
	MasterValidation() MasterValidation

	//Validations returns the endorsment's validations
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

func (e endorsement) LastTransactionHash() string {
	return e.lastTxHash
}

func (e endorsement) TransactionHash() string {
	return e.txHash
}

func (e endorsement) MasterValidation() MasterValidation {
	return e.masterValidation
}

func (e endorsement) Validations() []Validation {
	return e.validations
}
