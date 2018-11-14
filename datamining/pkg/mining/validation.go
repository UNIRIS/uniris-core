package mining

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

	//Status returns the validation's status
	Status() ValidationStatus

	//Timestamp returns the validation's timestamp
	Timestamp() time.Time

	//PublicKey returns the validation's public key
	PublicKey() string

	//Signature returns the validation's signature
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

func (v validation) Status() ValidationStatus {
	return v.status
}

func (v validation) Timestamp() time.Time {
	return v.timestamp
}

func (v validation) PublicKey() string {
	return v.pubk
}

func (v validation) Signature() string {
	return v.sig
}
