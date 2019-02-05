package rpc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Receive lock transaction request
	Given a transaction to lock
	When I want to request to lock it
	Then I get not error and the lock is stored
*/
func TestHandleLockTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	lockDB := &mockLockDb{}
	lockSrv := NewLockServer(lockDB, techDB)

	req := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := lockSrv.LockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, lockDB.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockDB.locks[0]["transaction_address"])
}

/*
Scenario: Receive unlock transaction request
	Given a transaction already
	When I want to request to unlock
	Then I get not error and the lock is removed
*/
func TestHandleUnlockTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	lockDB := &mockLockDb{}
	lockSrv := NewLockServer(lockDB, techDB)

	req := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := lockSrv.LockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, lockDB.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockDB.locks[0]["transaction_address"])

	req2 := &api.LockRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx"),
		MasterPeerPublicKey: pub,
		Address:             crypto.HashString("addr"),
	}
	reqBytes2, _ := json.Marshal(req2)
	sig2, _ := crypto.Sign(string(reqBytes2), pv)
	req2.SignatureRequest = sig2

	res2, err := lockSrv.UnlockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes2, _ := json.Marshal(&api.LockResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes2), pub, res2.SignatureResponse))

	assert.Len(t, lockDB.locks, 0)
}

type mockLockDb struct {
	locks []map[string]string
}

func (l *mockLockDb) WriteLock(txHash string, txAddr string) error {
	l.locks = append(l.locks, map[string]string{
		"transaction_address": txAddr,
		"transaction_hash":    txHash,
	})
	return nil
}
func (l *mockLockDb) RemoveLock(txHash string, txAddr string) error {
	pos := l.findLockPosition(txHash, txAddr)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l mockLockDb) ContainsLock(txHash string, txAddr string) (bool, error) {
	return l.findLockPosition(txHash, txAddr) > -1, nil
}

func (l mockLockDb) findLockPosition(txHash string, txAddr string) int {
	for i, lock := range l.locks {
		if lock["transaction_hash"] == txHash && lock["transaction_address"] == txAddr {
			return i
		}
	}
	return -1
}
