package datamining

import "time"

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = iota

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = 1
)

//Endorsement represents a validation
type Endorsement struct {
	timeStamp        time.Time
	txnHash          string
	masterValidation *MasterValidation
	validations      []Validation
}

//NewEndorsement creates a new endorsement
func NewEndorsement(t time.Time, h string, masterV *MasterValidation, valids []Validation) *Endorsement {
	return &Endorsement{
		timeStamp:        t,
		txnHash:          h,
		masterValidation: masterV,
		validations:      valids,
	}
}

//Timestamp returns the endorsment's timestamp
func (e Endorsement) Timestamp() time.Time {
	return e.timeStamp
}

//TransactionHash returns the endorsment's transaction hash
func (e Endorsement) TransactionHash() string {
	return e.txnHash
}

//MasterValidation returns the endorsment's master validation
func (e Endorsement) MasterValidation() *MasterValidation {
	return e.masterValidation
}

//Validations returns the endorsment's validations
func (e Endorsement) Validations() []Validation {
	return e.validations
}

//MasterValidation describe a validation of an elected master robot
type MasterValidation struct {
	lastTxRvk   []string
	powRobotKey string
	powValid    Validation
}

//NewMasterValidation creates a new master validation
func NewMasterValidation(lastTxRvk []string, powRobotKey string, powValid Validation) *MasterValidation {
	return &MasterValidation{lastTxRvk, powRobotKey, powValid}
}

//ValidatorKeysOfLastTransaction returns the list of public keys which validate the last transaction
func (m MasterValidation) ValidatorKeysOfLastTransaction() []string {
	return m.lastTxRvk
}

//ProofOfWorkRobotKey returns the public key of the robot which perform the PoW
func (m MasterValidation) ProofOfWorkRobotKey() string {
	return m.powRobotKey
}

//ProofOfWorkValidation returns the transaction proceed after the proof of work
func (m MasterValidation) ProofOfWorkValidation() Validation {
	return m.powValid
}

//Validation describe a validation of a robot
type Validation struct {
	status    ValidationStatus
	timestamp time.Time
	pubk      string
	sig       string
}

//NewValidation creates a new validation
func NewValidation(status ValidationStatus, t time.Time, pubKey string) Validation {
	return Validation{
		status:    status,
		timestamp: t,
		pubk:      pubKey,
	}
}

//Status returns the validation's status
func (v Validation) Status() ValidationStatus {
	return v.status
}

//Timestamp returns the validation's timestamp
func (v Validation) Timestamp() time.Time {
	return v.timestamp
}

//PublicKey returns the validation's public key
func (v Validation) PublicKey() string {
	return v.pubk
}

//Signature returns the validation's signature
func (v Validation) Signature() string {
	return v.sig
}
