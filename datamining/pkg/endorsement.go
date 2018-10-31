package datamining

import (
	"encoding/json"
	"time"
)

//Endorsement represents a validation
type Endorsement struct {
	timeStamp        time.Time
	txnHash          string
	masterValidation *MasterValidation
	validations      []Validation
}

func (e Endorsement) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp        time.Time         `json:"timestamp"`
		TransactionHash  string            `json:"transaction_hash"`
		MasterValidation *MasterValidation `json:"master_validation"`
		Validations      []Validation      `json:"validations"`
	}{
		Timestamp:        e.timeStamp,
		TransactionHash:  e.txnHash,
		MasterValidation: e.masterValidation,
		Validations:      e.validations,
	})
}

func (e *Endorsement) UnmarshalJSON(b []byte) error {
	data := struct {
		Timestamp        time.Time         `json:"timestamp"`
		TransactionHash  string            `json:"transaction_hash"`
		MasterValidation *MasterValidation `json:"master_validation"`
		Validations      []Validation      `json:"validations"`
	}{}

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	e.timeStamp = data.Timestamp
	e.txnHash = data.TransactionHash
	e.masterValidation = data.MasterValidation
	e.validations = data.Validations
	return nil
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
