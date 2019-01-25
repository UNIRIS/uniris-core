package uniris

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
	if t.Unix() > time.Now().Unix() {
		return MinerValidation{}, errors.New("Miner validation: timestamp must be lower than now")
	}
	if _, err := crypto.IsPublicKey(minerPubk); err != nil {
		return MinerValidation{}, fmt.Errorf("Miner validation: %s", err.Error())
	}

	switch status {
	case ValidationKO:
	case ValidationOK:
	default:
		return MinerValidation{}, errors.New("Miner validation: status not allowed")
	}

	if minerSig != "" {
		if _, err := crypto.IsSignature(minerSig); err != nil {
			return MinerValidation{}, fmt.Errorf("Miner validation: %s", err.Error())
		}
	}

	return MinerValidation{
		status:    status,
		timestamp: t,
		minerPubk: minerPubk,
		minerSig:  minerSig,
	}, nil
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

func (v MinerValidation) CheckValidation() error {
	if _, err := crypto.IsSignature(v.minerSig); err != nil {
		return fmt.Errorf("Miner validation: %s", err.Error())
	}
	vBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err = crypto.VerifySignature(string(vBytes), v.minerPubk, v.minerSig); err != nil {
		if err == crypto.ErrInvalidSignature {
			return errors.New("Miner validation: signature is invalid")
		}
		return err
	}
	return nil
}

func (v MinerValidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status    ValidationStatus `json:"status"`
		PublicKey string           `json:"public_key"`
		Timestamp time.Time        `json:"timestamp"`
		Signature string           `json:"signature,omitempty"`
	}{
		Status:    v.status,
		PublicKey: v.minerPubk,
		Timestamp: v.timestamp,
	})
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
	prevMiners []PeerIdentity
	pow        string
	validation MinerValidation
}

//NewMasterValidation creates a new master Transaction validation
func NewMasterValidation(prevMiners []PeerIdentity, pow string, valid MinerValidation) (MasterValidation, error) {

	ok, err := crypto.IsPublicKey(pow)
	if ok == false && err != nil {
		return MasterValidation{}, fmt.Errorf("Master validation Proof of work: %s", err.Error())
	}

	if valid == (MinerValidation{}) {
		return MasterValidation{}, errors.New("Master validation: Missing pre-validation")
	}

	return MasterValidation{
		prevMiners: prevMiners,
		pow:        pow,
		validation: valid,
	}, nil
}

//PreviousTransactionMiners returns the miners for the previous Transaction
func (mv MasterValidation) PreviousTransactionMiners() []PeerIdentity {
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
