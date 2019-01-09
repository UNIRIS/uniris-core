package uniris

import (
	"errors"
	"time"
)

//TransactionVerifier defines methods to verify transaction signatures
type TransactionVerifier interface {
	VerifyTransactionSignature(tx Transaction, pubKey string, sig string) (bool, error)
	VerifyValidationSignature(v MinerValidation) (bool, error)
}

//TransactionHasher defines methods to hash transaction
type TransactionHasher interface {
	HashTransaction(tx Transaction) (string, error)
}

//TransactionType represents the Transaction type
type TransactionType int

const (
	//KeychainTransactionType represents Transaction related to keychain
	KeychainTransactionType TransactionType = 0

	//IDTransactionType represents Transaction related to ID data
	IDTransactionType TransactionType = 1
)

//Transaction describe a transaction
type Transaction struct {
	idPubKey            string
	idSig               string
	emSig               string
	prop                Proposal
	txHash              string
	previousTransaction *Transaction
	mining              TransactionMining
}

//IDPublicKey returns transaction's ID public key
func (t Transaction) IDPublicKey() string {
	return t.idPubKey
}

//IDSignature returns transaction's ID signature
func (t Transaction) IDSignature() string {
	return t.idSig
}

//Proposal returns transaction's proposal
func (t Transaction) Proposal() Proposal {
	return t.prop
}

//EmitterSignature returns transaction's emitter signature (use to perform POW)
func (t Transaction) EmitterSignature() string {
	return t.emSig
}

//TransactionHash returns the transaction's hash
func (t Transaction) TransactionHash() string {
	return t.txHash
}

//PreviousTransaction returns the previous (chained) transaction
func (t Transaction) PreviousTransaction() *Transaction {
	if t.previousTransaction != nil {
		return t.previousTransaction
	}
	return nil
}

//Mining returns the mining of the Transaction
func (t Transaction) Mining() TransactionMining {
	return t.mining
}

//CheckChainTransactionIntegrity insure the transaction chain integrity
func (t Transaction) CheckChainTransactionIntegrity(h TransactionHasher, tv TransactionVerifier) error {
	if t.PreviousTransaction() != nil {
		if t.PreviousTransaction().TransactionHash() == "" {
			return errors.New("Transaction integrity violated")
		}
		return t.previousTransaction.CheckChainTransactionIntegrity(h, tv)
	}
	return t.CheckTransactionIntegrity(h, tv)
}

//CheckTransactionIntegrity insure the transaction integrity
func (t Transaction) CheckTransactionIntegrity(h TransactionHasher, tv TransactionVerifier) error {
	hash, err := h.HashTransaction(t)
	if err != nil {
		return err
	}
	if hash != t.TransactionHash() {
		return errors.New("Transaction integrity violated")
	}

	ok, err := tv.VerifyTransactionSignature(t, t.IDPublicKey(), t.IDSignature())
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
	ok, err := tv.VerifyTransactionSignature(t, t.mining.pow, t.emSig)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Invalid Proof of work")
	}
	return nil
}

//IsKO determinates is the transaction is KO (plan to be in the KO storage)
func (t Transaction) IsKO() bool {
	if t.mining.masterValidation.status == ValidationKO {
		return true
	}
	for _, v := range t.mining.validations {
		if v.Status() == ValidationKO {
			return true
		}
	}
	return false
}

//AddMining include the transaction mining
func (t *Transaction) AddMining(mining TransactionMining) {
	t.mining = mining
}

//TransactionMining describe the mining of a Transaction
type TransactionMining struct {
	prevMiners       []string
	pow              string
	masterValidation MinerValidation
	validations      []MinerValidation
}

//NewTransactionMining creates a new Transaction mining
func NewTransactionMining(prevMiners []string, pow string, masterV MinerValidation, valids []MinerValidation) TransactionMining {
	return TransactionMining{
		prevMiners:       prevMiners,
		pow:              pow,
		masterValidation: masterV,
		validations:      valids,
	}
}

//PreviousTransactionMiners returns the miners for the previous transaction
func (tx TransactionMining) PreviousTransactionMiners() []string {
	return tx.prevMiners
}

//ProofOfWork returns the transaction proof of work (emitter public key) validated the emitter signature
func (tx TransactionMining) ProofOfWork() string {
	return tx.pow
}

//MasterValidation returns the validation performed by the master peer
func (tx TransactionMining) MasterValidation() MinerValidation {
	return tx.masterValidation
}

//Validations returns the validations performed by the validation pool
func (tx TransactionMining) Validations() []MinerValidation {
	return tx.validations
}

//MinerValidation represents a transaction validation made by a miner
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
