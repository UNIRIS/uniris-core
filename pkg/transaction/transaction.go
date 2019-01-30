package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
)

//Status represents the status for the transaction
type Status int

const (
	//StatusUnknown define a transaction as unknown (the transaction hash is invalid)
	StatusUnknown Status = 0

	//StatusPending define a transaction in pending. (mining has not been finished)
	StatusPending Status = 2

	//StatusSuccess define a transaction in success (mining and storage succeed)
	StatusSuccess Status = 1

	//StatusFailure define a transaction in failure (mining failed due to an invalid transaction/signatures)
	StatusFailure Status = 3
)

//Type represents the Transaction type
type Type int

const (
	//KeychainType represents a Transaction related to keychain
	KeychainType Type = 0

	//IDType represents a Transaction related to ID data
	IDType Type = 1

	//ContractType represents a Transaction related to a smart contract
	ContractType Type = 2

	//ContractMessageType represents a Transaction related to a smart contract message
	ContractMessageType Type = 3
)

//Transaction describe a root Transaction
type Transaction struct {
	address       string
	txType        Type
	data          map[string]string
	timestamp     time.Time
	pubKey        string
	sig           string
	emSig         string
	prop          Proposal
	txHash        string
	prevTx        *Transaction
	masterV       MasterValidation
	confirmValids []MinerValidation
}

//New creates a transaction
func New(addr string, txType Type, data map[string]string, timestamp time.Time, pubK string, sig string, emSig string, prop Proposal, txHash string) (Transaction, error) {

	tx := Transaction{
		address:   addr,
		txType:    txType,
		data:      data,
		timestamp: timestamp,
		pubKey:    pubK,
		sig:       sig,
		emSig:     emSig,
		prop:      prop,
		txHash:    txHash,
	}
	if err := tx.checkFields(); err != nil {
		return Transaction{}, err
	}
	return tx, nil
}

//Address returns the Transaction's address (use for the sharding and identify the owner of the Transaction)
func (t Transaction) Address() string {
	return t.address
}

//Type returns the type of the Transaction
func (t Transaction) Type() Type {
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

//Proposal returns Transaction's proposal
func (t Transaction) Proposal() Proposal {
	return t.prop
}

//TransactionHash returns the Transaction's hash
func (t Transaction) TransactionHash() string {
	return t.txHash
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
	txHash := crypto.HashBytes(txBytesForHash)
	if txHash != t.TransactionHash() {
		return errors.New("transaction integrity violated")
	}

	txBytesBeforeSig, err := t.MarshalBeforeSignature()
	if err != nil {
		return err
	}

	err = crypto.VerifySignature(string(txBytesBeforeSig), t.PublicKey(), t.Signature())
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

	err = crypto.VerifySignature(string(txBytesBeforeSig), t.MasterValidation().ProofOfWork(), t.EmitterSignature())
	if err == crypto.ErrInvalidSignature {
		return errors.New("invalid proof of work")
	}
	return err
}

//IsKO determinates is the Transaction is KO (plan to be in the KO storage)
func (t Transaction) IsKO() bool {
	if t.masterV.Validation().Status() == ValidationKO {
		return true
	}
	for _, v := range t.confirmValids {
		if v.Status() == ValidationKO {
			return true
		}
	}
	return false
}

//AddMining add master validation and confirmation validations to the transaction
func (t *Transaction) AddMining(mv MasterValidation, confs []MinerValidation) error {
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
	return json.Marshal(struct {
		Address   string            `json:"address"`
		Data      map[string]string `json:"data"`
		Timestamp int64             `json:"timestamp"`
		Type      Type              `json:"type"`
		PublicKey string            `json:"public_key"`
		Proposal  Proposal          `json:"proposal"`
	}{
		Address:   t.Address(),
		Data:      t.Data(),
		Timestamp: t.Timestamp().Unix(),
		Type:      t.Type(),
		PublicKey: t.PublicKey(),
		Proposal:  t.Proposal(),
	})
}

//MarshalHash serializes as JSON the transaction to produce its hash
func (t Transaction) MarshalHash() ([]byte, error) {
	return json.Marshal(struct {
		Address          string            `json:"address"`
		Data             map[string]string `json:"data"`
		Timestamp        int64             `json:"timestamp"`
		Type             Type              `json:"type"`
		PublicKey        string            `json:"public_key"`
		Proposal         Proposal          `json:"proposal"`
		Signature        string            `json:"signature"`
		EmitterSignature string            `json:"em_signature"`
	}{
		Address:          t.Address(),
		Data:             t.Data(),
		Timestamp:        t.Timestamp().Unix(),
		Type:             t.Type(),
		PublicKey:        t.PublicKey(),
		Proposal:         t.Proposal(),
		Signature:        t.Signature(),
		EmitterSignature: t.EmitterSignature(),
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
	if _, err := crypto.IsHash(t.address); err != nil {
		return fmt.Errorf("transaction: %s", err.Error())
	}

	if _, err := crypto.IsHash(t.txHash); err != nil {
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
	case KeychainType:
	case IDType:
	case ContractType:
	case ContractMessageType:
	default:
		return errors.New("transaction: type not allowed")
	}

	if t.prop == (Proposal{}) {
		return errors.New("transaction: proposal is missing")
	}

	return nil
}
