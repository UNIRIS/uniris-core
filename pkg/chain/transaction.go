package chain

import (
	"bytes"
	"encoding/hex"
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
	TransactionStatusInProgress TransactionStatus = 1

	//TransactionStatusSuccess define a transaction in success (mining and storage succeed)
	TransactionStatusSuccess TransactionStatus = 2

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
	addr          crypto.VersionnedHash
	txType        TransactionType
	data          map[string][]byte
	timestamp     time.Time
	pubKey        crypto.PublicKey
	sig           []byte
	emSig         []byte
	prop          shared.EmitterCrossKeyPair
	hash          crypto.VersionnedHash
	prevTx        *Transaction
	masterV       MasterValidation
	confirmValids []Validation
}

//NewTransaction creates a new transaction
func NewTransaction(addr crypto.VersionnedHash, txType TransactionType, data map[string][]byte, timestamp time.Time, pubK crypto.PublicKey, prop shared.EmitterCrossKeyPair, sig []byte, emSig []byte, hash crypto.VersionnedHash) (Transaction, error) {
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
func (t Transaction) Address() crypto.VersionnedHash {
	return t.addr
}

//TransactionType returns the type of the Transaction
func (t Transaction) TransactionType() TransactionType {
	return t.txType
}

//Data returns Transaction's data
func (t Transaction) Data() map[string][]byte {
	return t.data
}

//Timestamp returns the Transaction sending timestamp
func (t Transaction) Timestamp() time.Time {
	return t.timestamp
}

//PublicKey returns Transaction's public key
func (t Transaction) PublicKey() crypto.PublicKey {
	return t.pubKey
}

//Signature returns Transaction's signature
func (t Transaction) Signature() []byte {
	return t.sig
}

//EmitterSignature returns Transaction's client signature (use to perform POW)
func (t Transaction) EmitterSignature() []byte {
	return t.emSig
}

//EmitterSharedKeyProposal returns Transaction's proposal
func (t Transaction) EmitterSharedKeyProposal() shared.EmitterCrossKeyPair {
	return t.prop
}

//TransactionHash returns the Transaction's hash
func (t Transaction) TransactionHash() crypto.VersionnedHash {
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
func (t Transaction) ConfirmationsValidations() []Validation {
	return t.confirmValids
}

//CheckChainTransactionIntegrity insure the Transaction chain integrity
func (t *Transaction) CheckChainTransactionIntegrity() error {
	if t.prevTx != nil {
		if !t.prevTx.TransactionHash().IsValid() {
			return errors.New("invalid previous transaction hash")
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

	hash := crypto.Hash(txBytesForHash)
	if !bytes.Equal(hash, t.hash) {
		return errors.New("transaction integrity violated")
	}

	txBytesBeforeSig, err := t.MarshalBeforeSignature()
	if err != nil {
		return err
	}

	if ok := t.pubKey.Verify(txBytesBeforeSig, t.sig); !ok {
		return errors.New("transaction signature invalid")
	}
	return nil
}

//CheckMasterValidation ensures the proof of work is valid
func (t Transaction) CheckMasterValidation() error {
	if _, err := t.masterV.IsValid(); err != nil {
		return err
	}

	txBytesBeforeEmSig, err := t.MarshalBeforeEmitterSignature()
	if err != nil {
		return err
	}

	if ok := t.masterV.pow.Verify(txBytesBeforeEmSig, t.emSig); !ok {
		return errors.New("invalid proof of work")
	}
	return nil
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
func (t *Transaction) Mined(mv MasterValidation, confs []Validation) error {
	t.masterV = mv
	if len(confs) == 0 {
		return errors.New("confirmation validations of the transaction are missing")
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

	pubK, err := t.pubKey.Marshal()
	if err != nil {
		return nil, err
	}

	propPub, err := t.prop.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	for k, v := range t.data {
		data[k] = hex.EncodeToString(v)
	}

	return json.Marshal(map[string]interface{}{
		"addr":       hex.EncodeToString(t.addr),
		"data":       data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": hex.EncodeToString(pubK),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(t.prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(propPub),
		},
	})
}

//MarshalBeforeEmitterSignature serializes as JSON the transaction before the emitter signature
func (t Transaction) MarshalBeforeEmitterSignature() ([]byte, error) {

	pubK, err := t.pubKey.Marshal()
	if err != nil {
		return nil, err
	}

	propPub, err := t.prop.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	for k, v := range t.data {
		data[k] = hex.EncodeToString(v)
	}

	return json.Marshal(map[string]interface{}{
		"addr":       hex.EncodeToString(t.addr),
		"data":       data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": hex.EncodeToString(pubK),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(t.prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(propPub),
		},
		"signature": hex.EncodeToString(t.sig),
	})
}

//MarshalHash serializes as JSON the transaction to produce its hash
func (t Transaction) MarshalHash() ([]byte, error) {
	pubK, err := t.pubKey.Marshal()
	if err != nil {
		return nil, err
	}

	propPub, err := t.prop.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	for k, v := range t.data {
		data[k] = hex.EncodeToString(v)
	}

	return json.Marshal(map[string]interface{}{
		"addr":       hex.EncodeToString(t.addr),
		"data":       data,
		"timestamp":  t.timestamp.Unix(),
		"type":       t.txType,
		"public_key": hex.EncodeToString(pubK),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(t.prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(propPub),
		},
		"signature":    hex.EncodeToString(t.sig),
		"em_signature": hex.EncodeToString(t.emSig),
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

	if len(t.addr) == 0 {
		return errors.New("transaction address is missing")
	}
	if !t.addr.IsValid() {
		return errors.New("transaction address is not a valid hash")
	}

	if t.pubKey == nil {
		return errors.New("transaction public key is missing")
	}

	if len(t.sig) == 0 {
		return errors.New("transaction signature is missing")
	}

	if len(t.emSig) == 0 {
		return errors.New("transaction emitter signature is missing")
	}

	if len(t.hash) == 0 {
		return errors.New("transaction hash is missing")
	}
	if !t.hash.IsValid() {
		return errors.New("transaction hash is not a valid hash")
	}

	if len(t.data) == 0 {
		return errors.New("transaction data is missing")
	}

	if t.timestamp.Unix() > time.Now().Unix() {
		return errors.New("transaction timestamp must be greater lower than now")
	}

	switch t.txType {
	case KeychainTransactionType:
	case IDTransactionType:
	case ContractTransactionType:
	case ContractMessageTransactionType:
	default:
		return errors.New("transaction type is not allowed")
	}

	if t.prop.EncryptedPrivateKey() == nil {
		return errors.New("transaction proposal private key is missing")
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

//Validation represents a Transaction validation made by a node
type Validation struct {
	status    ValidationStatus
	timestamp time.Time
	nodePubk  crypto.PublicKey
	nodeSig   []byte
}

//NewValidation creates a new node validation
func NewValidation(status ValidationStatus, t time.Time, nodePubk crypto.PublicKey, nodeSig []byte) (Validation, error) {
	v := Validation{
		status:    status,
		timestamp: t,
		nodePubk:  nodePubk,
		nodeSig:   nodeSig,
	}

	_, err := v.IsValid()
	if err != nil {
		return Validation{}, err
	}
	return v, nil
}

//Status return the validation status
func (v Validation) Status() ValidationStatus {
	return v.status
}

//Timestamp return the validation timestamp
func (v Validation) Timestamp() time.Time {
	return v.timestamp
}

//PublicKey return the node's public key performed this validation
func (v Validation) PublicKey() crypto.PublicKey {
	return v.nodePubk
}

//Signature returne the node's signature which performed this validation
func (v Validation) Signature() []byte {
	return v.nodeSig
}

//IsValid checks if the node validation is valid
func (v Validation) IsValid() (bool, error) {

	if v.nodePubk == nil {
		return false, errors.New("validation public key is missing")
	}

	if len(v.nodeSig) == 0 {
		return false, errors.New("validation signature is missing")
	}

	if v.timestamp.Unix() > time.Now().Unix() {
		return false, errors.New("validation timestamp must be anterior or equal to now")
	}

	switch v.status {
	case ValidationKO:
	case ValidationOK:
	default:
		return false, errors.New("validation status is not allowed")
	}

	vBytes, err := json.Marshal(v)
	if err != nil {
		return false, err
	}

	if !v.nodePubk.Verify(vBytes, v.nodeSig) {
		return false, errors.New("validation signature is not valid")
	}

	return true, nil
}

//MarshalJSON serializes as JSON a node validation
func (v Validation) MarshalJSON() ([]byte, error) {
	nodeKey, err := v.nodePubk.Marshal()
	if err != nil {
		return nil, err
	}
	return json.Marshal(map[string]interface{}{
		"status":     v.status,
		"public_key": nodeKey,
		"timestamp":  v.timestamp.Unix(),
	})
}

//MasterValidation describe the master Transaction validation
type MasterValidation struct {
	prevValidNodes []crypto.PublicKey
	pow            crypto.PublicKey
	validation     Validation
	wHeaders       []NodeHeader
	vHeaders       []NodeHeader
	sHeaders       []NodeHeader
}

//NodeHeader identifies a node entry in the transaction headers
type NodeHeader struct {
	pubKey        crypto.PublicKey
	isUnreachable bool
	isMaster      bool
	patchNumber   int
	isOK          bool
}

//NewNodeHeader creates a new node header
func NewNodeHeader(pubK crypto.PublicKey, isUnreachable bool, isMaster bool, patchNumber int, isOk bool) NodeHeader {
	return NodeHeader{
		pubKey:        pubK,
		isUnreachable: isUnreachable,
		isMaster:      isMaster,
		patchNumber:   patchNumber,
		isOK:          isOk,
	}
}

//PublicKey returns the node public key
func (h NodeHeader) PublicKey() crypto.PublicKey {
	return h.pubKey
}

//IsUnreachable determinates if the node is unreachable
func (h NodeHeader) IsUnreachable() bool {
	return h.isUnreachable
}

//IsMaster returns determinates if the node is a master
func (h NodeHeader) IsMaster() bool {
	return h.isMaster
}

//PatchNumber returns the node geo patch number
func (h NodeHeader) PatchNumber() int {
	return h.patchNumber
}

//IsOk determinates if the node is in status OK not faulty and not in bootstraping
func (h NodeHeader) IsOk() bool {
	return h.isOK
}

//NewMasterValidation creates a new master Transaction validation
func NewMasterValidation(prevValidNodes []crypto.PublicKey, pow crypto.PublicKey, valid Validation, wHeaders []NodeHeader, vHeaders []NodeHeader, sHeaders []NodeHeader) (MasterValidation, error) {
	mv := MasterValidation{
		prevValidNodes: prevValidNodes,
		pow:            pow,
		validation:     valid,
		wHeaders:       wHeaders,
		vHeaders:       vHeaders,
		sHeaders:       sHeaders,
	}
	if _, err := mv.IsValid(); err != nil {
		return MasterValidation{}, err
	}
	return mv, nil
}

//PreviousValidationNodes returns the validation nodes for the previous transaction
func (mv MasterValidation) PreviousValidationNodes() []crypto.PublicKey {
	return mv.prevValidNodes
}

//ProofOfWork returns the Transaction proof of work (emitter public key) validated the emitter signature
func (mv MasterValidation) ProofOfWork() crypto.PublicKey {
	return mv.pow
}

//Validation returns the mining performed by the master peer
func (mv MasterValidation) Validation() Validation {
	return mv.validation
}

//WelcomeHeaders returns the headers determining the master nodes election
func (mv MasterValidation) WelcomeHeaders() []NodeHeader {
	return mv.wHeaders
}

//ValidationHeaders returns the headers determining the validation nodes election
func (mv MasterValidation) ValidationHeaders() []NodeHeader {
	return mv.vHeaders
}

//StorageHeaders returns the node headers determining the storage nodes election
func (mv MasterValidation) StorageHeaders() []NodeHeader {
	return mv.sHeaders
}

//IsValid check is the master validation is correct
func (mv MasterValidation) IsValid() (bool, error) {

	if len(mv.wHeaders) == 0 {
		return false, fmt.Errorf("master validation: missing welcome node headers")
	}

	if len(mv.vHeaders) == 0 {
		return false, fmt.Errorf("master validation: missing validation pool headers")
	}

	if len(mv.sHeaders) == 0 {
		return false, fmt.Errorf("master validation: missing storage pool headers")
	}

	//TODO: ensure the validaty of the headers

	if mv.ProofOfWork() == nil {
		return false, errors.New("proof of work is missing")
	}

	if _, err := mv.Validation().IsValid(); err != nil {
		return false, fmt.Errorf("master validation is not valid: %s", err.Error())
	}

	return true, nil
}
