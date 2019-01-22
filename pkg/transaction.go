package uniris

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
)

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
	prop          TransactionProposal
	txHash        string
	prevTx        *Transaction
	masterV       MasterValidation
	confirmValids []MinerValidation
}

//NewTransactionBase creates a basic transaction
func NewTransactionBase(addr string, txType TransactionType, data string, timestamp time.Time, pubK string, sig string, emSig string, prop TransactionProposal, txHash string) Transaction {
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
func (t Transaction) Proposal() TransactionProposal {
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
		if t.prevTx.TransactionHash() == "" {
			return errors.New("Transaction integrity violated")
		}
		return t.prevTx.CheckChainTransactionIntegrity()
	}
	return t.CheckTransactionIntegrity()
}

//CheckTransactionIntegrity insure the Transaction integrity
func (t Transaction) CheckTransactionIntegrity() error {
	txBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	txHash := crypto.HashBytes(txBytes)
	if txHash != t.TransactionHash() {
		return errors.New("Transaction integrity violated")
	}

	err = crypto.VerifySignature(string(txBytes), t.PublicKey(), t.Signature())
	if err == crypto.ErrInvalidSignature {
		return errors.New("Transaction signature invalid")
	}
	return err
}

//CheckProofOfWork ensures the proof of work is valid
func (t Transaction) CheckProofOfWork() error {
	if t.MasterValidation().ProofOfWork() == "" || t.MasterValidation().Validation() == (MinerValidation{}) {
		return errors.New("Missing master validation")
	}

	txBytes, err := json.Marshal(t)
	if err != nil {
		return err
	}

	err = crypto.VerifySignature(string(txBytes), t.MasterValidation().ProofOfWork(), t.EmitterSignature())
	if err == crypto.ErrInvalidSignature {
		return errors.New("Invalid Proof of work")
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

func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address   string
		Data      string
		Type      TransactionType
		PublicKey string
		Proposal  TransactionProposal
	}{
		Address:   t.Address(),
		Data:      t.Data(),
		Type:      t.Type(),
		PublicKey: t.PublicKey(),
		Proposal:  t.Proposal(),
	})
}
