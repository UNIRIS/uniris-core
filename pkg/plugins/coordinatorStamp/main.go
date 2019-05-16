package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
}

type validationStamp interface {
	IsValid() (bool, error)
}

type electedNode interface {
	MarshalJSON() ([]byte, error)
}

type electedNodeList interface {
	Nodes() []electedNode
	CreatorPublicKey() publicKey
	CreatorSignature() []byte
}

//CoordinatorStamp represents a coordinator's stamp
type CoordinatorStamp interface {
	PreviousCrossValidators() [][]byte
	ProofOfWork() interface{}
	ValidationStamp() interface{}
	TransactionHash() []byte
	ElectedCoordinatorNodes() interface{}
	ElectedCrossValidationNodes() interface{}
	ElectedStorageNodes() interface{}
	IsValid() (bool, string)
}

//CoordinatorStamp describe the Transaction validation made by the coordinator
type coordinatorStamp struct {
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

	cs := coordinatorStamp{
		prevCrossV:      prevCrossV,
		pow:             pow,
		validStamp:      validStamp,
		txHash:          txHash,
		elecCoordNodes:  elecCoordNodes,
		elecCrossVNodes: elecCrossVNodes,
		elecStorNodes:   elecStorNodes,
	}

	if ok, reason := cs.IsValid(); !ok {
		return nil, errors.New(reason)
	}

	return cs, nil
}

func (c coordinatorStamp) PreviousCrossValidators() [][]byte {
	return c.prevCrossV
}

func (c coordinatorStamp) ProofOfWork() interface{} {
	return c.pow
}

func (c coordinatorStamp) ValidationStamp() interface{} {
	return c.validStamp
}

func (c coordinatorStamp) TransactionHash() []byte {
	return c.txHash
}

func (c coordinatorStamp) ElectedCoordinatorNodes() interface{} {
	return c.elecCoordNodes
}

func (c coordinatorStamp) ElectedCrossValidationNodes() interface{} {
	return c.elecCrossVNodes
}

func (c coordinatorStamp) ElectedStorageNodes() interface{} {
	return c.elecStorNodes
}

func (c coordinatorStamp) IsValid() (bool, string) {

	if c.pow == nil {
		return false, "coordinator stamp: proof of work is missing"
	}

	if _, ok := c.pow.(publicKey); !ok {
		return false, "coordinator stamp: proof of work is not a valid public key"
	}

	vStamp, ok := c.validStamp.(validationStamp)
	if !ok {
		return false, "coordinator stamp: invalid validation stamp"
	}

	if _, err := vStamp.IsValid(); err != nil {
		return false, fmt.Sprintf("coordinator stamp: invalid validations stamp: %s", err.Error())
	}

	if c.txHash == nil {
		return false, "coordinator stamp: missing transaction hash"
	}

	if c.elecCoordNodes == nil {
		return false, "coordinator stamp: missing elected coordinator nodes"
	}

	elecCoordN, ok := c.elecCoordNodes.(electedNodeList)
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

	if c.elecCrossVNodes == nil {
		return false, "coordinator stamp: missing elected cross validation nodes"
	}

	elecValidN, ok := c.elecCrossVNodes.(electedNodeList)
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

	if c.elecStorNodes == nil {
		return false, "coordinator stamp: missing elected storage nodes"
	}

	elecStrN, ok := c.elecStorNodes.(electedNodeList)
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
