package uniris

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

//NewTransaction creates a transaction
func NewTransaction(addr string, txType TransactionType, data string, timestamp time.Time, pubK string, sig string, emSig string, prop TransactionProposal, txHash string) (Transaction, error) {
	if _, err := crypto.IsHash(addr); err != nil {
		return Transaction{}, fmt.Errorf("Transaction: %s", err.Error())
	}

	if _, err := crypto.IsHash(txHash); err != nil {
		return Transaction{}, fmt.Errorf("Transaction: %s", err.Error())
	}

	if data == "" {
		return Transaction{}, errors.New("Transaction: data is empty")
	}
	if _, err := hex.DecodeString(data); err != nil {
		return Transaction{}, errors.New("Transaction: data is not in hexadecimal format")
	}

	if timestamp.Unix() > time.Now().Unix() {
		return Transaction{}, errors.New("Transaction: timestamp must be greater lower than now")
	}

	if _, err := crypto.IsPublicKey(pubK); err != nil {
		return Transaction{}, fmt.Errorf("Transaction: %s", err.Error())
	}

	if _, err := crypto.IsSignature(sig); err != nil {
		return Transaction{}, fmt.Errorf("Transaction: %s", err.Error())
	}

	if _, err := crypto.IsSignature(emSig); err != nil {
		return Transaction{}, fmt.Errorf("Transaction: %s", err.Error())
	}

	switch txType {
	case KeychainTransactionType:
	case IDTransactionType:
	case ContractTransactionType:
	case ContractMessageTransactionType:
	default:
		return Transaction{}, errors.New("Transaction: type not allowed")
	}

	if prop == (TransactionProposal{}) {
		return Transaction{}, errors.New("Transaction: proposal is missing")
	}

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
	}, nil
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

func (t *Transaction) AddMining(mv MasterValidation, confs []MinerValidation) error {
	t.masterV = mv
	if len(confs) == 0 {
		return errors.New("Transaction: Missing confirmation validations")
	}

	t.confirmValids = confs
	return nil
}

func (t *Transaction) Chain(prevTx *Transaction) {
	if prevTx != nil && prevTx.TransactionHash() != "" {
		t.prevTx = prevTx
	}
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Address   string              `json:"address"`
		Data      string              `json:"data"`
		Type      TransactionType     `json:"type"`
		PublicKey string              `json:"public_key"`
		Proposal  TransactionProposal `json:"proposal"`
	}{
		Address:   t.Address(),
		Data:      t.Data(),
		Type:      t.Type(),
		PublicKey: t.PublicKey(),
		Proposal:  t.Proposal(),
	})
}

//TransactionProposal describe a proposal for a Transaction
type TransactionProposal struct {
	sharedEmitterKP SharedKeys
}

//NewTransactionProposal create a new proposal for a Transaction
func NewTransactionProposal(shdEmitterKP SharedKeys) (TransactionProposal, error) {
	if (shdEmitterKP == SharedKeys{}) {
		return TransactionProposal{}, errors.New("Transaction proposal: missing shared keys")
	}
	return TransactionProposal{
		sharedEmitterKP: shdEmitterKP,
	}, nil
}

//SharedEmitterKeyPair returns the keypair proposed for the shared emitter keys
func (p TransactionProposal) SharedEmitterKeyPair() SharedKeys {
	return p.sharedEmitterKP
}

func (p TransactionProposal) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SharedEmitterKP SharedKeys `json:"shared_emitter_keys"`
	}{
		SharedEmitterKP: p.SharedEmitterKeyPair(),
	})
}
