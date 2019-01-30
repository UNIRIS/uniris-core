package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
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
func NewMinerValidation(status ValidationStatus, t time.Time, minerPubk string, minerSig string) (MinerValidation, error) {
	v := MinerValidation{
		status:    status,
		timestamp: t,
		minerPubk: minerPubk,
		minerSig:  minerSig,
	}

	_, err := v.IsValid()
	if err != nil {
		return MinerValidation{}, err
	}
	return v, nil
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

//IsValid checks if the miner validation is valid
func (v MinerValidation) IsValid() (bool, error) {

	if v.timestamp.Unix() > time.Now().Unix() {
		return false, errors.New("miner validation: timestamp must be anterior or equal to now")
	}

	if _, err := crypto.IsPublicKey(v.minerPubk); err != nil {
		return false, fmt.Errorf("miner validation: %s", err.Error())
	}
	switch v.status {
	case ValidationKO:
	case ValidationOK:
	default:
		return false, errors.New("miner validation: status not allowed")
	}

	if _, err := crypto.IsSignature(v.minerSig); err != nil {
		return false, fmt.Errorf("miner validation: %s", err.Error())
	}
	vBytes, err := json.Marshal(v)
	if err != nil {
		return false, err
	}
	if err := crypto.VerifySignature(string(vBytes), v.minerPubk, v.minerSig); err != nil {
		if err == crypto.ErrInvalidSignature {
			return false, errors.New("miner validation: signature is invalid")
		}
		return false, err
	}
	return true, nil
}

//MarshalJSON serializes as JSON a miner validation
func (v MinerValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status    ValidationStatus `json:"status"`
		PublicKey string           `json:"public_key"`
		Timestamp int64            `json:"timestamp"`
	}{
		Status:    v.status,
		PublicKey: v.minerPubk,
		Timestamp: v.timestamp.Unix(),
	})
}

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = iota

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = 1
)

//MasterValidation describe the master Transaction validation
type MasterValidation struct {
	prevMiners Pool
	pow        string
	validation MinerValidation
}

//NewMasterValidation creates a new master Transaction validation
func NewMasterValidation(prevMiners Pool, pow string, valid MinerValidation) (MasterValidation, error) {
	mv := MasterValidation{
		prevMiners: prevMiners,
		pow:        pow,
		validation: valid,
	}
	if _, err := mv.IsValid(); err != nil {
		return MasterValidation{}, err
	}
	return mv, nil
}

//PreviousTransactionMiners returns the miners for the previous Transaction
func (mv MasterValidation) PreviousTransactionMiners() Pool {
	return mv.prevMiners
}

//ProofOfWork returns the Transaction proof of work (emitter public key) validated the emitter signature
func (mv MasterValidation) ProofOfWork() string {
	return mv.pow
}

//Validation returns the mining performed by the master peer
func (mv MasterValidation) Validation() MinerValidation {
	return mv.validation
}

//IsValid check is the master validation is correct
func (mv MasterValidation) IsValid() (bool, error) {
	if _, err := crypto.IsPublicKey(mv.ProofOfWork()); err != nil {
		return false, fmt.Errorf("master validation POW: %s", err.Error())
	}

	if _, err := mv.Validation().IsValid(); err != nil {
		return false, fmt.Errorf("master validation: %s", err.Error())
	}

	return true, nil
}
