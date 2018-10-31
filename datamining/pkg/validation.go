package datamining

import (
	"encoding/json"
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
type Validation struct {
	status    ValidationStatus
	timestamp time.Time
	pubk      string
	sig       string
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

//LastTransactionMiners returns the list of public keys which validate the last transaction
func (m MasterValidation) LastTransactionMiners() []string {
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

//NewValidation creates a new validation
func NewValidation(status ValidationStatus, t time.Time, pubKey string, sig string) Validation {
	return Validation{
		status:    status,
		timestamp: t,
		pubk:      pubKey,
		sig:       sig,
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

func (v Validation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status    ValidationStatus `json:"status"`
		Timestamp time.Time        `json:"timestamp"`
		PublicKey string           `json:"public_key"`
		Signature string           `json:"signature"`
	}{
		Status:    v.status,
		Timestamp: v.timestamp,
		PublicKey: v.pubk,
		Signature: v.sig,
	})
}

func (m *Validation) UnmarshalJSON(b []byte) error {
	data := struct {
		Status    ValidationStatus `json:"status"`
		Timestamp time.Time        `json:"timestamp"`
		PublicKey string           `json:"public_key"`
		Signature string           `json:"signature"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	m.status = data.Status
	m.timestamp = data.Timestamp
	m.pubk = data.PublicKey
	m.sig = data.Signature

	return nil
}
