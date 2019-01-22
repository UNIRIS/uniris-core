package uniris

import (
	"errors"
	"time"
)

//TransactionVerifier defines methods to verify Transaction signatures
type TransactionVerifier interface {
	VerifyTransactionSignature(tx Transaction, pubKey string, sig string) (bool, error)
	VerifyValidationSignature(v MinerValidation) (bool, error)
}

//TransactionHasher defines methods to hash Transaction
type TransactionHasher interface {
	HashTransaction(tx Transaction) (string, error)
}

//TransactionStatus represents the status for the transaction
type TransactionStatus int

const (
	//UnknownTransaction define a transaction as unknown (the transaction hash is invalid)
	UnknownTransaction TransactionStatus = 0

	//PendingTransaction define a transaction in pending. (mining has not been finished)
	PendingTransaction TransactionStatus = 2

	//SuccessTransaction define a transaction in success (mining and storage succeed)
	SuccessTransaction TransactionStatus = 1

	//FailureTransaction define a transaction in failure (mining failed due to an invalid transaction/signatures)
	FailureTransaction TransactionStatus = 3
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
	address       string
	txType        TransactionType
	data          string
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

//NewTransactionBase creates a basic transaction
func NewTransactionBase(addr string, txType TransactionType, data string, timestamp time.Time, pubK string, sig string, emSig string, prop Proposal, txHash string) Transaction {
	return Transaction{
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
}

//NewChainedTransaction creates a transaction chained to another
func NewChainedTransaction(tx Transaction, prevTx Transaction) Transaction {
	return Transaction{
		address:   tx.address,
		txType:    tx.txType,
		data:      tx.data,
		timestamp: tx.timestamp,
		pubKey:    tx.pubKey,
		sig:       tx.sig,
		emSig:     tx.emSig,
		prop:      tx.prop,
		txHash:    tx.txHash,
		prevTx:    &prevTx,
	}
}

//NewMinedTransaction creates a mined transaction
func NewMinedTransaction(tx Transaction, masterV MasterValidation, confirms []MinerValidation) Transaction {
	return Transaction{
		address:       tx.address,
		txType:        tx.txType,
		data:          tx.data,
		timestamp:     tx.timestamp,
		pubKey:        tx.pubKey,
		sig:           tx.sig,
		emSig:         tx.emSig,
		prop:          tx.prop,
		txHash:        tx.txHash,
		prevTx:        tx.prevTx,
		masterV:       masterV,
		confirmValids: confirms,
	}
}

//Address returns the Transaction's address (use for the sharding and identify the owner of the Transaction)
func (t Transaction) Address() string {
	return t.address
}

//Type returns the type of the Transaction
func (t Transaction) Type() TransactionType {
	return t.txType
}

//Data returns Transaction's data
func (t Transaction) Data() string {
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
func (t *Transaction) CheckChainTransactionIntegrity(h TransactionHasher, tv TransactionVerifier) error {
	if t.prevTx != nil {
		if t.prevTx.TransactionHash() == "" {
			return errors.New("Transaction integrity violated")
		}
		return t.prevTx.CheckChainTransactionIntegrity(h, tv)
	}
	return t.CheckTransactionIntegrity(h, tv)
}

//CheckTransactionIntegrity insure the Transaction integrity
func (t Transaction) CheckTransactionIntegrity(h TransactionHasher, tv TransactionVerifier) error {
	hash, err := h.HashTransaction(t)
	if err != nil {
		return err
	}
	if hash != t.TransactionHash() {
		return errors.New("Transaction integrity violated")
	}

	ok, err := tv.VerifyTransactionSignature(t, t.PublicKey(), t.Signature())
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Transaction signature invalid")
	}

	return nil
}

//CheckProofOfWork ensures the proof of work is valid
func (t Transaction) CheckProofOfWork(tv TransactionVerifier) error {
	if t.MasterValidation().ProofOfWork() == "" || t.MasterValidation().Validation() == (MinerValidation{}) {
		return errors.New("Missing master validation")
	}

	ok, err := tv.VerifyTransactionSignature(t, t.MasterValidation().ProofOfWork(), t.emSig)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Invalid Proof of work")
	}
	return nil
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
func (v MinerValidation) CheckValidation(tv TransactionVerifier) error {
	ok, err := tv.VerifyValidationSignature(v)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Invalid validation signature")
	}
	return nil
}

//ValidationStatus defines a validation status
type ValidationStatus int

const (

	//ValidationOK defines when a validation successed
	ValidationOK ValidationStatus = iota

	//ValidationKO defines when a validation failed
	ValidationKO ValidationStatus = 1
)

//Proposal describe a proposal for a Transaction
type Proposal interface {

	//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
	SharedEmitterKeyPair() SharedKeys
}

type prop struct {
	sharedEmitterKP SharedKeys
}

//NewProposal create a new proposal for a Transaction
func NewProposal(shdEmitterKP SharedKeys) Proposal {
	return prop{
		sharedEmitterKP: shdEmitterKP,
	}
}

func (p prop) SharedEmitterKeyPair() SharedKeys {
	return p.sharedEmitterKP
}
