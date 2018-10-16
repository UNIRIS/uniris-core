package datamining

type ValidationStatus int

const (
	ValidationOK ValidationStatus = iota
	ValidationKO ValidationStatus = 1
)

//Endorsement represents a validation
type Endorsement struct {
	timeStamp        Timestamp
	txnHash          Hash
	masterValidation MasterValidation
	validations      []Validation
}

//NewEndorsement creates a new endorsement
func NewEndorsement(t Timestamp, h Hash, masterV MasterValidation, valids []Validation) Endorsement {
	return Endorsement{
		timeStamp:        t,
		txnHash:          h,
		masterValidation: masterV,
		validations:      valids,
	}
}

func (e Endorsement) Timestamp() Timestamp {
	return e.timeStamp
}

func (e Endorsement) TransactionHash() Hash {
	return e.txnHash
}

func (e Endorsement) MasterValidation() MasterValidation {
	return e.masterValidation
}

func (e Endorsement) Validations() []Validation {
	return e.validations
}

//MasterValidation describe a validation of an elected master robot
type MasterValidation struct {
	lastTxRvk   []PublicKey
	powRobotKey PublicKey
	powValid    Validation
}

func NewMasterValidation(lastTxRvk []PublicKey, powRobotKey PublicKey, powValid Validation) MasterValidation {
	return MasterValidation{lastTxRvk, powRobotKey, powValid}
}

//ValidatorKeysOfLastTransaction returns the list of public keys which validate the last transaction
func (m MasterValidation) ValidatorKeysOfLastTransaction() []PublicKey {
	return m.lastTxRvk
}

//ProofOfWorkRobotKey returns the public key of the robot which perform the PoW
func (m MasterValidation) ProofOfWorkRobotKey() PublicKey {
	return m.powRobotKey
}

//ProofOfWorkValidation returns the transaction proceed after the proof of work
func (m MasterValidation) ProofOfWorkValidation() Validation {
	return m.powValid
}

//Validation describe a validation of a robot
type Validation struct {
	status    ValidationStatus
	timestamp Timestamp
	pubk      PublicKey
	sig       DERSignature
}

//NewValidation creates a new validation
func NewValidation(status ValidationStatus, t Timestamp, pubKey PublicKey) Validation {
	return Validation{
		status:    status,
		timestamp: t,
		pubk:      pubKey,
	}
}

func (v *Validation) AddSignature(sig DERSignature) {
	v.sig = sig
}

//Status ..
func (v Validation) Status() ValidationStatus {
	return v.status
}

//Timestamp ...
func (v Validation) Timestamp() Timestamp {
	return v.timestamp
}

//Pubk ...
func (v Validation) Pubk() PublicKey {
	return v.pubk
}

//Sig ...
func (v Validation) Sig() DERSignature {
	return v.sig
}
