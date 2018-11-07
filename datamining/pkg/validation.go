package datamining

import (
	"time"
)

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = iota

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = 1
)

//Validation describe a validation of a robot
type Validation interface {
	Status() ValidationStatus
	Timestamp() time.Time
	PublicKey() string
	Signature() string
}

type validation struct {
	status    ValidationStatus
	timestamp time.Time
	pubk      string
	sig       string
}

//NewValidation creates a new validation
func NewValidation(status ValidationStatus, t time.Time, pubKey string, sig string) Validation {
	return validation{
		status:    status,
		timestamp: t,
		pubk:      pubKey,
		sig:       sig,
	}
}

//Status returns the validation's status
func (v validation) Status() ValidationStatus {
	return v.status
}

//Timestamp returns the validation's timestamp
func (v validation) Timestamp() time.Time {
	return v.timestamp
}

//PublicKey returns the validation's public key
func (v validation) PublicKey() string {
	return v.pubk
}

//Signature returns the validation's signature
func (v validation) Signature() string {
	return v.sig
}
