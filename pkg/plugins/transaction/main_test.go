package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
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
Scenario: Create a new transaction
	Given transaction data (addr, type, timestamp, data, public key, signature, originSig)
	When I want to create the transaction
	Then I get it
*/
func TestNewTransaction(t *testing.T) {
	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	data := map[string]interface{}{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	tJSON, _ := json.Marshal(map[string]interface{}{
		"addr":       addr,
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       KeychainTransactionType,
		"public_key": mockPublicKey{bytes: pub}.Marshal(),
	})
	sig, _ := mockPrivateKey{bytes: pv}.Sign(tJSON)

	tx, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), mockPublicKey{bytes: pub}, sig, sig, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, addr, tx.(transaction).Address())
	assert.Equal(t, data, tx.(transaction).Data())
	assert.Equal(t, KeychainTransactionType, tx.(transaction).Type())
	assert.Equal(t, mockPublicKey{bytes: pub}, tx.(transaction).PreviousPublicKey())
	assert.Equal(t, sig, tx.(transaction).Signature())
	assert.Equal(t, sig, tx.(transaction).OriginSignature())
}

/*
Scenario: Create a new transaction with an invalid addr
	Given a invalid addr hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidAddress(t *testing.T) {
	_, err := NewTransaction(nil, KeychainTransactionType, nil, time.Now(), mockPublicKey{}, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: address is missing")

	_, err = NewTransaction([]byte("abc"), KeychainTransactionType, nil, time.Now(), mockPublicKey{}, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: address is an invalid hash")
}

/*
Scenario: Create a new transaction without public key
	Given a transaction without public key or with invalid public key type
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidPublicKey(t *testing.T) {

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	fake := struct{}{}
	_, err := NewTransaction(addr, KeychainTransactionType, nil, time.Now(), fake, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: invalid public key")
}

/*
Scenario: Create a new transaction without signature
	Given a transaction without signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutSignature(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	_, err := NewTransaction(addr, KeychainTransactionType, nil, time.Now(), mockPublicKey{bytes: pub}, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: signature is missing")
}

/*
Scenario: Create a new transaction without origin signature
	Given a transaction without origin signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutOriginSignature(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	_, err := NewTransaction(addr, KeychainTransactionType, nil, time.Now(), mockPublicKey{bytes: pub}, []byte("sig"), nil, nil, nil)
	assert.EqualError(t, err, "transaction: origin signature is missing")
}

/*
Scenario: Create a new transaction without data
	Given a transaction without data
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutData(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], addr)

	_, err := NewTransaction(addr, KeychainTransactionType, nil, time.Now(), mockPublicKey{bytes: pub}, []byte("sig"), []byte("sig"), nil, nil)
	assert.EqualError(t, err, "transaction: data is missing")
}

/*
Scenario: Create a new transaction with invalid timestamp
	Given a transaction with older timestamp
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidTimestamp(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	_, err := NewTransaction(addr, KeychainTransactionType, map[string]interface{}{
		"test": "test",
	}, time.Now().Add(2*time.Second), mockPublicKey{bytes: pub}, []byte("sig"), []byte("sig"), nil, nil)
	assert.EqualError(t, err, "transaction: invalid timestamp")
}

/*
Scenario: Create a new transaction with invalid type
	Given a transaction with invalid type
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidType(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	data := map[string]interface{}{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	_, err := NewTransaction(addr, 200, data, time.Now(), mockPublicKey{bytes: pub}, []byte("sig"), []byte("sig"), nil, nil)
	assert.EqualError(t, err, "transaction: invalid type")
}

/*
Scenario: Create a new transaction with invalid signature
	Given a transaction with invalid signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidSignature(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	data := map[string]interface{}{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	_, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), mockPublicKey{bytes: pub}, []byte("sig"), []byte("sig"), nil, nil)
	assert.EqualError(t, err, "transaction: invalid signature")
}

/*
Scenario: Create a new transaction with invalid coordinator stamp
	Given a transaction with an invalid coordinator stamp type
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidCoordinatorStampType(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	data := map[string]interface{}{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], hOut)

	tJSON, _ := json.Marshal(map[string]interface{}{
		"addr":       addr,
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       KeychainTransactionType,
		"public_key": mockPublicKey{bytes: pub}.Marshal(),
	})
	sig, _ := mockPrivateKey{bytes: pv}.Sign(tJSON)

	fake := struct{}{}
	_, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), mockPublicKey{bytes: pub}, sig, sig, fake, nil)
	assert.EqualError(t, err, "transaction: invalid coordinator stamp")
}

/*
Scenario: Create a new transaction with invalid cross valdiation stamp
	Given a transaction with an invalid cross valdiation stamp
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidCrossValidatorStamp(t *testing.T) {

	pub, pv, _ := ed25519.GenerateKey(rand.Reader)

	data := map[string]interface{}{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	h := sha256.New()
	h.Write([]byte("addr"))
	hOut := h.Sum(nil)
	addr := make([]byte, 1+len(hOut))
	addr[0] = byte(int(crypto.SHA256))
	copy(addr[1:], addr)

	tJSON, _ := json.Marshal(map[string]interface{}{
		"addr":       addr,
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       KeychainTransactionType,
		"public_key": mockPublicKey{bytes: pub}.Marshal(),
	})
	sig, _ := mockPrivateKey{bytes: pv}.Sign(tJSON)

	tJSON, _ = json.Marshal(map[string]interface{}{
		"addr":       addr,
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       KeychainTransactionType,
		"public_key": mockPublicKey{bytes: pub}.Marshal(),
		"signature":  sig,
	})
	originSig, _ := mockPrivateKey{bytes: pv}.Sign(tJSON)
	tJSON, _ = json.Marshal(map[string]interface{}{
		"addr":             addr,
		"data":             data,
		"timestamp":        time.Now().Unix(),
		"type":             KeychainTransactionType,
		"public_key":       mockPublicKey{bytes: pub}.Marshal(),
		"signature":        sig,
		"origin_signature": originSig,
	})
	h = sha256.New()
	h.Write(tJSON)
	hOut = h.Sum(nil)
	txHash := make([]byte, 1+len(hOut))
	txHash[0] = byte(int(crypto.SHA256))
	copy(txHash[1:], hOut)

	cs := mockCoordinatorStamp{
		txHash: txHash,
		pow:    mockPublicKey{bytes: pub},
	}
	fake := []interface{}{struct{}{}}
	_, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), mockPublicKey{bytes: pub}, sig, sig, cs, fake)
	assert.EqualError(t, err, "transaction: cross validation type is invalid")
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

type mockValidationStamp struct {
	status    int
	timestamp time.Time
	pubKey    mockPublicKey
	sig       []byte
}

func (v mockValidationStamp) Status() int {
	return v.status
}
func (v mockValidationStamp) Timestamp() time.Time {
	return v.timestamp
}
func (v mockValidationStamp) NodePublicKey() interface{} {
	return v.pubKey
}
func (v mockValidationStamp) NodeSignature() []byte {
	return v.sig
}

type mockCoordinatorStamp struct {
	txHash []byte
	pow    mockPublicKey
}

func (ct mockCoordinatorStamp) ProofOfWork() interface{} {
	return ct.pow
}
func (ct mockCoordinatorStamp) TransactionHash() []byte {
	return ct.txHash
}
