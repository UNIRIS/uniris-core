package uniris

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//MinerValidation represents a Transaction validation made by a miner
type MinerValidation struct {
	status    ValidationStatus
	timestamp time.Time
	minerPubk string
	minerSig  string
}

//NewMinerValidation creates a new miner validation
func NewMinerValidation(status ValidationStatus, t time.Time, minerPubk string, minerSig string) MinerValidation {
	return MinerValidation{
		status:    status,
		timestamp: t,
		minerPubk: minerPubk,
		minerSig:  minerSig,
	}
}

//Status return the validation status
func (v MinerValidation) Status() ValidationStatus {
	return v.status
}

//Timestamp return the validation timestamp
func (v MinerValidation) Timestamp() time.Time {
	return v.timestamp
}

//MinerPublicKey return the miner's public key performed this validation
func (v MinerValidation) MinerPublicKey() string {
	return v.minerPubk
}

//MinerSignature returne the miner's signature which performed this validation
func (v MinerValidation) MinerSignature() string {
	return v.minerSig
}

//CheckValidation insures the validation signature is correct
func (v MinerValidation) CheckValidation() error {
	vBytes, err := json.Marshal(struct {
		Status    ValidationStatus `json:"status"`
		PublicKey string           `json:"public_key"`
		Timestamp time.Time        `json:"timestamp"`
	}{
		Status:    v.Status(),
		PublicKey: v.MinerPublicKey(),
		Timestamp: v.Timestamp(),
	})
	if err != nil {
		return err
	}
	err = crypto.VerifySignature(string(vBytes), v.MinerPublicKey(), v.MinerSignature())
	if err == crypto.ErrInvalidSignature {
		return errors.New("Invalid validation signature")
	}
	return err
}

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = iota

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = 1
)

//MasterValidation describe the master Transaction validation
type MasterValidation struct {
	prevMiners []string
	pow        string
	validation MinerValidation
}

//NewMasterValidation creates a new master Transaction validation
func NewMasterValidation(prevMiners []string, pow string, valid MinerValidation) MasterValidation {
	return MasterValidation{
		prevMiners: prevMiners,
		pow:        pow,
		validation: valid,
	}
}

//PreviousTransactionMiners returns the miners for the previous Transaction
func (mv MasterValidation) PreviousTransactionMiners() []string {
	return mv.prevMiners
}

//ProofOfWork returns the Transaction proof of work (emitter public key) validated the emitter signature
func (mv MasterValidation) ProofOfWork() string {
	return mv.pow
}

//Validation returns the mining performed by the master peer
func (mv MasterValidation) Validation() MinerValidation {
	return mv.Validation()
}
