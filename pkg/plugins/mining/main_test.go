package main

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	os.Setenv("PLUGINS_DIR", filepath.Join(dir, "../"))
	m.Run()
}

/*
Scenario: Perform Proof of work
	Given a transaction and em chain keypair stored
	When I want to perform the proof of work of this transaction
	Then I get the valid public key
*/
func TestPerformPOW(t *testing.T) {

	pub1, _, _ := ed25519.GenerateKey(rand.Reader)
	pub2, _, _ := ed25519.GenerateKey(rand.Reader)
	pub3, pv3, _ := ed25519.GenerateKey(rand.Reader)

	tx := mockTransaction{}
	tRaw, _ := tx.MarshalBeforeOriginSignature()
	sig, _ := mockPrivateKey{bytes: pv3}.Sign(tRaw)
	tx.originSig = sig

	pow, err := performPOW(tx, []interface{}{mockPublicKey{bytes: pub1}, mockPublicKey{bytes: pub2}, mockPublicKey{bytes: pub3}})
	assert.Nil(t, err)
	assert.EqualValues(t, pub3, pow.(mockPublicKey).Marshal()[1:])
}

func TestCrossValidateTransaction(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)
	tx := mockTransaction{}
	tRaw, _ := tx.MarshalBeforeOriginSignature()
	sig, _ := mockPrivateKey{bytes: pv}.Sign(tRaw)
	tx.originSig = sig

	stamp, err := CrossValidateTransaction(t, mockPublicKey{bytes: pub}, mockPrivateKey{bytes: pv})
	assert.Nil(t, err)
	assert.NotNil(t, stamp)
	assert.EqualValues(t, pub, stamp.(mockValidation).NodePublicKey().(publicKey).Marshal()[1:])
	assert.NotNil(t, stamp.(mockValidation).NodeSignature())
	assert.Equal(t, 1, stamp.(mockValidation).Status())
}

//TODO: Make test for the coordinator transaction processing

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

type mockTransaction struct {
	addr      []byte
	txType    int
	data      map[string]interface{}
	timestamp time.Time
	pub       publicKey
	sig       []byte
	originSig []byte
	coord     interface{}
	crossV    []interface{}
}

func (t mockTransaction) Address() []byte {
	return t.addr
}
func (t mockTransaction) Type() int {
	return t.txType
}
func (t mockTransaction) Data() map[string]interface{} {
	return t.data
}
func (t mockTransaction) Timestamp() time.Time {
	return t.timestamp
}
func (t mockTransaction) PreviousPublicKey() interface{} {
	return t.pub
}
func (t mockTransaction) Signature() []byte {
	return t.sig
}
func (t mockTransaction) OriginSignature() []byte {
	return t.originSig
}
func (t mockTransaction) CoordinatorStamp() interface{} {
	return t.coord
}
func (t mockTransaction) CrossValidations() []interface{} {
	return t.crossV
}

func (t mockTransaction) MarshalRoot() ([]byte, error) {
	return []byte("dummy"), nil
}
func (t mockTransaction) MarshalBeforeOriginSignature() ([]byte, error) {
	return []byte("dummy"), nil
}

type mockValidation interface {
	Status() int
	Timestamp() time.Time
	NodePublicKey() interface{}
	NodeSignature() []byte
}
