package main

import (
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"golang.org/x/crypto/ed25519"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a new coordinator stamp
	Given a proof of work, node validation and the elected nodes
	When I want to create the coordinator stamp
	Then I get a valid stamp
*/
func TestNewCoordinatorStamp(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	pubK := mockPublicKey{bytes: pub}
	pvK := mockPrivateKey{bytes: pv}
	v := mockValidation{
		pubKey: pubK,
		status: 1,
		ts:     time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := pvK.Sign(b)
	v.sig = sig

	elecNode := mockElectedNode{
		isCoord:       true,
		isOK:          true,
		isUnreachable: false,
		patchNumber:   1,
		publicKey:     pubK,
	}

	coordN := mockElectedNodeList{nodes: []electedNode{elecNode}, creatorPublicKey: pubK}
	coordNJSON, _ := json.Marshal(coordN.nodes)
	coordSig, _ := pvK.Sign(coordNJSON)
	coordN.creatorSignature = coordSig

	crossVN := mockElectedNodeList{nodes: []electedNode{elecNode}, creatorPublicKey: pubK}
	crossVNJSON, _ := json.Marshal(crossVN.nodes)
	crossVNSig, _ := pvK.Sign(crossVNJSON)
	crossVN.creatorSignature = crossVNSig

	storN := mockElectedNodeList{nodes: []electedNode{elecNode}, creatorPublicKey: pubK}
	storNJSON, _ := json.Marshal(storN.nodes)
	storNSig, _ := pvK.Sign(storNJSON)
	storN.creatorSignature = storNSig

	coorStmp, err := NewCoordinatorStamp(nil, pubK, v, []byte("hash"), coordN, crossVN, storN)
	assert.Nil(t, err)
	assert.EqualValues(t, pubK, coorStmp.(CoordinatorStamp).ProofOfWork().(publicKey))
	assert.Equal(t, v, coorStmp.(CoordinatorStamp).ValidationStamp())
	assert.Empty(t, coorStmp.(CoordinatorStamp).PreviousCrossValidators())
}

/*
Scenario: Create a coordinator stamp with POW invalid
	Given a no POW or invalid public key
	When I want to create coordinate stamp
	Then I get an error indicating the POW is missing or invalid
*/
func TestCreateCoordinatorStampWithInvalidPOW(t *testing.T) {
	_, err := NewCoordinatorStamp(nil, nil, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, "coordinator stamp: proof of work is missing")

	pub := struct{}{}

	_, err = NewCoordinatorStamp(nil, pub, nil, nil, nil, nil, nil)
	assert.EqualError(t, err, "coordinator stamp: proof of work is not a valid public key")
}

//TODO: add tests to handle the check of the coordinator nodes
//TODO: add tests to handle the check of the cross validation nodes
//TODO: add tests to handle the check of the storage nodes

type mockValidation struct {
	pubKey publicKey
	sig    []byte
	ts     time.Time
	status int
}

func (v mockValidation) IsValid() (bool, error) {
	return true, nil
}

type mockPublicKey struct {
	bytes []byte
	curve int
}

func (pb mockPublicKey) Marshal() []byte {
	out := make([]byte, 1+len(pb.bytes))
	out[0] = byte(int(pb.curve))
	copy(out[1:], pb.bytes)
	return out
}

func (pb mockPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	return ed25519.Verify(pb.bytes, data, sig), nil
}

type mockPrivateKey struct {
	bytes []byte
}

func (pv mockPrivateKey) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(pv.bytes, data), nil
}

type mockElectedNode struct {
	publicKey     mockPublicKey
	isUnreachable bool
	isCoord       bool
	patchNumber   int
	isOK          bool
}

func (e mockElectedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"publicKey":     e.publicKey.Marshal(),
		"isUnreachable": e.isUnreachable,
		"isCoordinator": e.isCoord,
		"patchNumber":   e.patchNumber,
		"isOk":          e.isOK,
	})
}

type mockElectedNodeList struct {
	nodes            []electedNode
	creatorPublicKey publicKey
	creatorSignature []byte
}

func (e mockElectedNodeList) Nodes() []electedNode {
	return e.nodes
}

func (e mockElectedNodeList) CreatorPublicKey() publicKey {
	return e.creatorPublicKey
}

func (e mockElectedNodeList) CreatorSignature() []byte {
	return e.creatorSignature
}
