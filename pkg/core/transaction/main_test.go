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
	assert.Equal(t, addr, tx.(Transaction).Address())
	assert.Equal(t, data, tx.(Transaction).Data())
	assert.Equal(t, KeychainTransactionType, tx.(Transaction).Type())
	assert.Equal(t, mockPublicKey{bytes: pub}, tx.(Transaction).PreviousPublicKey())
	assert.Equal(t, sig, tx.(Transaction).Signature())
	assert.Equal(t, sig, tx.(Transaction).OriginSignature())
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

	_, err := NewTransaction(addr, KeychainTransactionType, nil, time.Now(), nil, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: public key is missing")

	fake := struct{}{}
	_, err = NewTransaction(addr, KeychainTransactionType, nil, time.Now(), fake, nil, nil, nil, nil)
	assert.EqualError(t, err, "transaction: public key type is invalid")
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
	assert.EqualError(t, err, "transaction: coordinator stamp type is invalid")
}

/*
Scenario: Create a new transaction with invalid coordinator stamp
	Given a transaction with an invalid coordinator stamp values
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidCoordinatorStampValues(t *testing.T) {

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
	assert.Equal(t, err, "transaction: cross validation type is invalid")
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

type mockValidation struct {
	pubKey publicKey
	sig    []byte
	ts     time.Time
	status int
}

func (v mockValidation) IsValid() (bool, string) {
	return true, ""
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
func (ct mockCoordinatorStamp) ValidationStamp() validationStamp {
	return mockValidation{}
}
func (ct mockCoordinatorStamp) IsValid() (bool, string) {
	return true, ""
}
