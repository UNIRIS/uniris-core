package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
}

type electedNode interface {
	MarshalJSON() ([]byte, error)
}

type electedNodeList interface {
	Nodes() []interface{}
	CreatorPublicKey() interface{}
	CreatorSignature() []byte
}

type coordinatorStamp interface {
	PreviousCrossValidators() [][]byte
	ProofOfWork() interface{}
	ValidationStamp() interface{}
	TransactionHash() []byte
	ElectedCoordinatorNodes() interface{}
	ElectedCrossValidationNodes() interface{}
	ElectedStorageNodes() interface{}
}

//coordStmp describe the Transaction validation made by the coordinator
type coordStmp struct {
	prevCrossV      [][]byte
	pow             interface{}
	validStamp      interface{}
	txHash          []byte
	elecCoordNodes  interface{}
	elecCrossVNodes interface{}
	elecStorNodes   interface{}
}

//NewCoordinatorStamp creates a new coordinator stamp
func NewCoordinatorStamp(prevCrossV [][]byte, pow interface{}, validStamp interface{}, txHash []byte, elecCoordNodes interface{}, elecCrossVNodes interface{}, elecStorNodes interface{}) (interface{}, error) {

	cs := coordStmp{
		prevCrossV:      prevCrossV,
		pow:             pow,
		validStamp:      validStamp,
		txHash:          txHash,
		elecCoordNodes:  elecCoordNodes,
		elecCrossVNodes: elecCrossVNodes,
		elecStorNodes:   elecStorNodes,
	}

	if ok, reason := IsCoordinatorStampValid(cs); !ok {
		return nil, errors.New(reason)
	}

	return cs, nil
}

func (c coordStmp) PreviousCrossValidators() [][]byte {
	return c.prevCrossV
}

func (c coordStmp) ProofOfWork() interface{} {
	return c.pow
}

func (c coordStmp) ValidationStamp() interface{} {
	return c.validStamp
}

func (c coordStmp) TransactionHash() []byte {
	return c.txHash
}

func (c coordStmp) ElectedCoordinatorNodes() interface{} {
	return c.elecCoordNodes
}

func (c coordStmp) ElectedCrossValidationNodes() interface{} {
	return c.elecCrossVNodes
}

func (c coordStmp) ElectedStorageNodes() interface{} {
	return c.elecStorNodes
}

//IsCoordinatorStampValid checks if the coordinator stamp is valid
func IsCoordinatorStampValid(c interface{}) (bool, string) {

	cs, ok := c.(coordinatorStamp)
	if !ok {
		return false, "coordinator stamp: not valid"
	}

	if cs.ProofOfWork() == nil {
		return false, "coordinator stamp: proof of work is missing"
	}

	if _, ok := cs.ProofOfWork().(publicKey); !ok {
		return false, "coordinator stamp: proof of work is not a valid public key"
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "validationStamp/plugin.so"))
	if err != nil {
		return false, fmt.Sprintf("coordinator stamp: %s", err.Error())
	}
	sym, err := p.Lookup("IsValidStamp")
	if err != nil {
		return false, fmt.Sprintf("coordinator stamp: %s", err.Error())
	}
	isValidF := sym.(func(interface{}) (bool, string))

	if ok, reason := isValidF(cs.ValidationStamp()); !ok {
		return false, fmt.Sprintf("coordinator stamp: invalid validation stamp: %s", reason)
	}

	if cs.TransactionHash() == nil {
		return false, "coordinator stamp: missing transaction hash"
	}

	if cs.ElectedCoordinatorNodes() == nil {
		return false, "coordinator stamp: missing elected coordinator nodes"
	}

	elecCoordN, ok := cs.ElectedCoordinatorNodes().(electedNodeList)
	if !ok {
		return false, "coordinator stamp: invalid elected coordinator nodes"
	}

	if len(elecCoordN.Nodes()) == 0 {
		return false, "coordinator stamp: missing elected coordinates nodes"
	}
	if elecCoordN.CreatorPublicKey() == nil {
		return false, "coordinator stamp: missing elected coordinates nodes creator public key"
	}
	if elecCoordN.CreatorSignature() == nil {
		return false, "coordinator stamp: missing elected coordinates nodes creator signature"
	}

	elecNJSON, err := json.Marshal(elecCoordN.Nodes())
	if err != nil {
		return false, err.Error()
	}
	if ok, err := elecCoordN.CreatorPublicKey().(publicKey).Verify(elecNJSON, elecCoordN.CreatorSignature()); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "coordinator stamp: invalid elected coordinate node signature"
	}

	if cs.ElectedCrossValidationNodes() == nil {
		return false, "coordinator stamp: missing elected cross validation nodes"
	}

	elecValidN, ok := cs.ElectedCrossValidationNodes().(electedNodeList)
	if !ok {
		return false, "coordinator stamp: invalid elected cross validation nodes"
	}

	if len(elecValidN.Nodes()) == 0 {
		return false, "coordinator stamp: missing elected  cross validation nodes"
	}
	if elecValidN.CreatorPublicKey() == nil {
		return false, "coordinator stamp: missing elected cross validation nodes creator public key"
	}
	if elecValidN.CreatorSignature() == nil {
		return false, "coordinator stamp: missing elected cross validation nodes creator signature"
	}

	elecNJSON, err = json.Marshal(elecValidN.Nodes())
	if err != nil {
		return false, err.Error()
	}
	if ok, err := elecValidN.CreatorPublicKey().(publicKey).Verify(elecNJSON, elecValidN.CreatorSignature()); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "coordinator stamp: invalid elected cross validation node signature"
	}

	if cs.ElectedStorageNodes() == nil {
		return false, "coordinator stamp: missing elected storage nodes"
	}

	elecStrN, ok := cs.ElectedStorageNodes().(electedNodeList)
	if !ok {
		return false, "coordinator stamp: invalid elected storage nodes"
	}

	if len(elecStrN.Nodes()) == 0 {
		return false, "coordinator stamp: missing elected storage nodes"
	}

	if elecStrN.CreatorPublicKey() == nil {
		return false, "coordinator stamp: missing elected storage nodes creator public key"
	}
	if elecStrN.CreatorSignature() == nil {
		return false, "coordinator stamp: missing elected storage nodes creator signature"
	}

	elecNJSON, err = json.Marshal(elecStrN.Nodes())
	if err != nil {
		return false, err.Error()
	}
	if ok, err := elecStrN.CreatorPublicKey().(publicKey).Verify(elecNJSON, elecStrN.CreatorSignature()); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "coordinator stamp: invalid elected storage node signature"
	}

	return true, ""
}
