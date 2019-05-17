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

type validationStamp interface {
	Status() int
	Timestamp() time.Time
	NodePublicKey() interface{}
	NodeSignature() []byte
}

type vStamp struct {
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

	v := vStamp{
		status:    status,
		timestamp: t,
		nodePubk:  nodePubk,
		nodeSig:   nodeSig,
	}

	if ok, reason := IsValidStamp(v); !ok {
		return nil, errors.New(reason)
	}
	return v, nil
}

//Status return the validation status
func (v vStamp) Status() int {
	return v.status
}

//Timestamp return the validation timestamp
func (v vStamp) Timestamp() time.Time {
	return v.timestamp
}

//PublicKey return the node's public key performed this validation
func (v vStamp) NodePublicKey() interface{} {
	return v.nodePubk
}

//Signature return the node's signature which performed this validation
func (v vStamp) NodeSignature() []byte {
	return v.nodeSig
}

//IsValidStamp checks if the validation stamp is valid
func IsValidStamp(v interface{}) (bool, string) {

	stamp, ok := v.(validationStamp)
	if !ok {
		return false, "validation stamp: is not valid"
	}

	if stamp.NodePublicKey() == nil {
		return false, "validation stamp: public key is missing"
	}

	nodePub, ok := stamp.NodePublicKey().(publicKey)
	if !ok {
		return false, "validation stamp: public key is not valid"
	}

	if len(stamp.NodeSignature()) == 0 {
		return false, "validation stamp: signature is missing"
	}

	if stamp.Timestamp().Unix() > time.Now().Unix() {
		return false, "validation stamp: timestamp must be anterior or equal to now"
	}

	switch stamp.Status() {
	case ValidationKO:
	case ValidationOK:
	default:
		return false, "validation stamp: invalid status"
	}

	vBytes, err := json.Marshal(v)
	if err != nil {
		return false, err.Error()
	}

	if ok, err := nodePub.Verify(vBytes, stamp.NodeSignature()); err != nil {
		return false, err.Error()
	} else if !ok {
		return false, "validation stamp: signature is not valid"
	}

	return true, ""
}

//MarshalJSON serializes as JSON a validation stamp
func (v vStamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"status":     v.status,
		"public_key": v.nodePubk.(publicKey).Marshal(),
		"timestamp":  v.timestamp.Unix(),
	})
}
