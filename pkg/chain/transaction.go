package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//TransactionStatus represents the status for the transaction
type TransactionStatus int

const (
	//TransactionStatusUnknown define a transaction as unknown (the transaction hash is invalid)
	TransactionStatusUnknown TransactionStatus = 0

	//TransactionStatusInProgress define a transaction in in progress. (mining has not been finished)
	TransactionStatusInProgress TransactionStatus = 2

	//TransactionStatusSuccess define a transaction in success (mining and storage succeed)
	TransactionStatusSuccess TransactionStatus = 1

	//TransactionStatusFailure define a transaction in failure (mining failed due to an invalid transaction/signatures)
	TransactionStatusFailure TransactionStatus = 3
)

//TransactionType represents the Transaction type
type TransactionType int

const (
	//KeychainTransactionType represents a Transaction related to keychain
	KeychainTransactionType TransactionType = 0

	//IDTransactionType represents a Transaction related to ID data
	IDTransactionType TransactionType = 1

	//ContractTransactionType represents a Transaction related to a smart contract
	ContractTransactionType TransactionType = 2

	//ContractMessageTransactionType represents a Transaction related to a smart contract message
	ContractMessageTransactionType TransactionType = 3
)

//Transaction describe a root Transaction
type Transaction struct {
	addr          string
	txType        TransactionType
	data          map[string]string
	timestamp     time.Time
	pubKey        string
	sig           string
	emSig         string
	prop          shared.EmitterKeyPair
	hash          string
	prevTx        *Transaction
	masterV       MasterValidation
	confirmValids []MinerValidation
}

//NewTransaction creates a new transaction
func NewTransaction(addr string, txType TransactionType, data map[string]string, timestamp time.Time, pubK string, prop shared.EmitterKeyPair, sig string, emSig string, hash string) (Transaction, error) {
	tx := Transaction{
		addr:      addr,
		txType:    txType,
		data:      data,
		timestamp: timestamp,
		pubKey:    pubK,
		sig:       sig,
		emSig:     emSig,
		prop:      prop,
		hash:      hash,
	}
	if err := tx.checkFields(); err != nil {
		return Transaction{}, err
	}
	return tx, nil
}

//Address returns the Transaction's addr (use for the sharding and identify the owner of the Transaction)
func (t Transaction) Address() string {
	return t.addr
}

//TransactionType returns the type of the Transaction
func (t Transaction) TransactionType() TransactionType {
	return t.txType
}

//Data returns Transaction's data
func (t Transaction) Data() map[string]string {
	return t.data
}

//Timestamp returns the Transaction sending timestamp
func (t Transaction) Timestamp() time.Time {
	return t.timestamp
}

//PublicKey returns Transaction's public key
func (t Transaction) PublicKey() string {
	return t.pubKey
}

//Signature returns Transaction's signature
func (t Transaction) Signature() string {
	return t.sig
}

//EmitterSignature returns Transaction's client signature (use to perform POW)
func (t Transaction) EmitterSignature() string {
	return t.emSig
}

//EmitterSharedKeyProposal returns Transaction's proposal
func (t Transaction) EmitterSharedKeyProposal() shared.EmitterKeyPair {
	return t.prop
}

//TransactionHash returns the Transaction's hash
func (t Transaction) TransactionHash() string {
	return t.hash
}

//PreviousTransaction returns the previous (chained) Transaction
func (t Transaction) PreviousTransaction() *Transaction {
	if t.prevTx != nil {
		return t.prevTx
	}
	return nil
}

//MasterValidation returns the Transaction validation performed by the master peer (including the Proof of Work)
func (t Transaction) MasterValidation() MasterValidation {
	return t.masterV
}

//ConfirmationsValidations returns the Transaction confirmation validations performed by the validation pool
func (t Transaction) ConfirmationsValidations() []MinerValidation {
	return t.confirmValids
}

//CheckChainTransactionIntegrity insure the Transaction chain integrity
func (t *Transaction) CheckChainTransactionIntegrity() error {
	if t.prevTx != nil {
		if _, err := crypto.IsHash(t.prevTx.TransactionHash()); err != nil {
			return err
		}
		if t.prevTx.timestamp.Unix() >= t.timestamp.Unix() {
			return errors.New("previous chained transaction must be anterior to the current transaction")
		}
		return t.prevTx.CheckChainTransactionIntegrity()
	}
	return t.checkTransactionIntegrity()
}

func (t Transaction) checkTransactionIntegrity() error {
	txBytesForHash, err := t.MarshalHash()
	if err != nil {
		return err
	}
	hash := crypto.HashBytes(txBytesForHash)
	if hash != t.hash {
		return errors.New("transaction integrity violated")
	}

	txBytesBeforeSig, err := t.MarshalBeforeSignature()
	if err != nil {
		return err
	}

	err = crypto.VerifySignature(string(txBytesBeforeSig), t.pubKey, t.sig)
	if err == crypto.ErrInvalidSignature {
		return errors.New("transaction signature invalid")
	}
	return err
}

//CheckMasterValidation ensures the proof of work is valid
func (t Transaction) CheckMasterValidation() error {
	if _, err := t.masterV.IsValid(); err != nil {
		return err
	}

	txBytesBeforeSig, err := t.MarshalBeforeSignature()
	if err != nil {
		return err
	}

	err = crypto.VerifySignature(string(txBytesBeforeSig), t.masterV.pow, t.emSig)
	if err == crypto.ErrInvalidSignature {
		return errors.New("invalid proof of work")
	}
	return err
}

//IsKO determinates is the Transaction is KO (plan to be in the KO storage)
func (t Transaction) IsKO() bool {
	if t.masterV.validation.status == ValidationKO {
		return true
	}
	for _, v := range t.confirmValids {
		if v.status == ValidationKO {
			return true
		}
	}
	return false
}

//Mined define the transaction as mined by providing the master validation and confirmation validations
func (t *Transaction) Mined(mv MasterValidation, confs []MinerValidation) error {
	t.masterV = mv
	if len(confs) == 0 {
		return errors.New("transaction: missing confirmation validations")
	}

	t.confirmValids = confs
	return nil
}

//Chain links a transaction to another one
func (t *Transaction) Chain(prevTx *Transaction) error {
	if prevTx != nil {
		if prevTx.timestamp.Unix() >= t.timestamp.Unix() {
			return errors.New("previous chained transaction must be anterior to the current transaction")
		}
		if err := prevTx.CheckChainTransactionIntegrity(); err != nil {
			return err
		}
		t.prevTx = prevTx
	}

	return nil
}

//MarshalBeforeSignature serializes as JSON the transaction before its signature
func (t Transaction) MarshalBeforeSignature() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":                    t.addr,
		"data":                    t.data,
		"timestamp":               t.timestamp.Unix(),
		"type":                    t.txType,
		"public_key":              t.pubKey,
		"em_shared_keys_proposal": t.prop,
	})
}

//MarshalHash serializes as JSON the transaction to produce its hash
func (t Transaction) MarshalHash() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"addr":                    t.addr,
		"data":                    t.data,
		"timestamp":               t.timestamp.Unix(),
		"type":                    t.txType,
		"public_key":              t.pubKey,
		"em_shared_keys_proposal": t.prop,
		"signature":               t.sig,
		"em_signature":            t.emSig,
	})
}

//IsValid checks if the transaction fields are valid and its the integrity is respected
func (t Transaction) IsValid() (bool, error) {
	if err := t.checkFields(); err != nil {
		return false, err
	}
	if err := t.checkTransactionIntegrity(); err != nil {
		return false, err
	}
	return true, nil
}

func (t Transaction) checkFields() error {
	if _, err := crypto.IsHash(t.addr); err != nil {
		return fmt.Errorf("transaction: addr %s", err.Error())
	}

	if _, err := crypto.IsHash(t.hash); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	if len(t.data) == 0 {
		return errors.New("transaction: data is empty")
	}

	if t.timestamp.Unix() > time.Now().Unix() {
		return errors.New("transaction: timestamp must be greater lower than now")
	}

	if _, err := crypto.IsPublicKey(t.pubKey); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	if _, err := crypto.IsSignature(t.sig); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	if _, err := crypto.IsSignature(t.emSig); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	switch t.txType {
	case KeychainTransactionType:
	case IDTransactionType:
	case ContractTransactionType:
	case ContractMessageTransactionType:
	default:
		return errors.New("transaction: type not allowed")
	}

	if t.prop == (shared.EmitterKeyPair{}) {
		return errors.New("transaction: proposal is missing")
	}

	return nil
}

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = iota

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = 1
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
	return json.Marshal(map[string]interface{}{
		"status":     v.status,
		"public_key": v.minerPubk,
		"timestamp":  v.timestamp.Unix(),
	})
}

//MasterValidation describe the master Transaction validation
type MasterValidation struct {
	prevMiners []string
	pow        string
	validation MinerValidation
}

//NewMasterValidation creates a new master Transaction validation
func NewMasterValidation(prevMiners []string, pow string, valid MinerValidation) (MasterValidation, error) {
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
func (mv MasterValidation) PreviousTransactionMiners() []string {
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

	//Ensure the previous miners are public keys
	if len(mv.prevMiners) > 0 {
		for _, m := range mv.prevMiners {
			if _, err := crypto.IsPublicKey(m); err != nil {
				return false, err
			}
		}
	}

	if _, err := crypto.IsPublicKey(mv.ProofOfWork()); err != nil {
		return false, fmt.Errorf("master validation POW: %s", err.Error())
	}

	if _, err := mv.Validation().IsValid(); err != nil {
		return false, fmt.Errorf("master validation: %s", err.Error())
	}

	return true, nil
}
