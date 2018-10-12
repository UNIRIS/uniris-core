package datamining

import "time"

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

//MasterValidation describe a validation of an elected master robot
type MasterValidation struct {
	lastTxnv []PublicKey
	pow      PublicKey
	rbv      Validation
}

//LastTxnv ...
func (m MasterValidation) LastTxnv() []PublicKey {
	return m.lastTxnv
}

//Pow ...
func (m MasterValidation) Pow() PublicKey {
	return m.pow
}

//Rbv ...
func (m MasterValidation) Rbv() Validation {
	return m.rbv
}

//Validation describe a validation of a robot
type Validation struct {
	status    string
	timestamp time.Time
	pubk      PublicKey
	sig       DERSignature
}

//Status ..
func (m Validation) Status() string {
	return m.status
}

//Timestamp ...
func (m Validation) Timestamp() time.Time {
	return m.timestamp
}

//Pubk ...
func (m Validation) Pubk() PublicKey {
	return m.pubk
}

//Sig ...
func (m Validation) Sig() DERSignature {
	return m.sig
}
