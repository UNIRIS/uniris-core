package main

import (
	"encoding/json"
	"errors"
	"time"
)

var (

	//ValidationKO defines when a validation failed
	ValidationKO = 0

	//ValidationOK defines when a validation successed
	ValidationOK = 1
)

//ValidationStamp represents a Transaction validation stamp
type ValidationStamp interface {
	Status() int
	Timestamp() time.Time
	NodePublicKey() interface{}
	NodeSignature() []byte
	IsValid() (bool, string)
}

type validationStamp struct {
	status    int
	timestamp time.Time
	nodePubk  interface{}
	nodeSig   []byte
}

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Marshal() []byte
}

//NewValidationStamp creates a new validation stamp
func NewValidationStamp(status int, t time.Time, nodePubk interface{}, nodeSig []byte) (interface{}, error) {

	v := validationStamp{
		status:    status,
		timestamp: t,
		nodePubk:  nodePubk,
		nodeSig:   nodeSig,
	}

	if ok, reason := v.IsValid(); !ok {
		return nil, errors.New(reason)
	}
	return v, nil
}

//Status return the validation status
func (v validationStamp) Status() int {
	return v.status
}

//Timestamp return the validation timestamp
func (v validationStamp) Timestamp() time.Time {
	return v.timestamp
}

//PublicKey return the node's public key performed this validation
func (v validationStamp) NodePublicKey() interface{} {
	return v.nodePubk
}

//Signature return the node's signature which performed this validation
func (v validationStamp) NodeSignature() []byte {
	return v.nodeSig
}

//IsValid checks if the node validation is valid
func (v validationStamp) IsValid() (bool, string) {

	if v.nodePubk == nil {
		return false, "validation stamp: public key is missing"
	}

	nodePub, ok := v.nodePubk.(publicKey)
	if !ok {
		return false, "validation stamp: public key is not valid"
	}

	if len(v.nodeSig) == 0 {
		return false, "validation stamp: signature is missing"
	}

	if v.timestamp.Unix() > time.Now().Unix() {
		return false, "validation stamp: timestamp must be anterior or equal to now"
	}

	switch v.status {
	case ValidationKO:
	case ValidationOK:
	default:
		return false, "validation stamp: invalid status"
	}

	vBytes, err := json.Marshal(v)
	if err != nil {
		return false, err.Error()
	}

	if ok, err := nodePub.Verify(vBytes, v.nodeSig); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "validation stamp: signature is not valid"
	}

	return true, ""
}

//MarshalJSON serializes as JSON a validation stamp
func (v validationStamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"status":     v.status,
		"public_key": v.nodePubk.(publicKey).Marshal(),
		"timestamp":  v.timestamp.Unix(),
	})
}
