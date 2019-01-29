package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create transaction lock
	Given a transaction hash, an address and a master robot key
	When I want to lock the transaction
	Then the lock is stored
*/
func TestStoreLock(t *testing.T) {

	repo := &mockLockRepository{}
	s := LockService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	lock, err := NewLock(crypto.HashString("hash"), crypto.HashString("addr"), hex.EncodeToString(pub))
	assert.Nil(t, err)
	assert.Nil(t, s.StoreLock(lock))

	assert.Len(t, repo.locks, 1)
	assert.Equal(t, crypto.HashString("hash"), repo.locks[0].TransactionHash())
	assert.Equal(t, crypto.HashString("addr"), repo.locks[0].Address())
	assert.Equal(t, hex.EncodeToString(pub), repo.locks[0].MasterRobotKey())
}

/*
Scenario: Create two lock identicals
	Given a lock created
	When a want to created again the same lock
	Then I get an error
*/
func TestCreatedExistingLock(t *testing.T) {
	repo := &mockLockRepository{}
	s := LockService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	lock, err := NewLock(crypto.HashString("hash"), crypto.HashString("addr"), hex.EncodeToString(pub))
	assert.Nil(t, err)
	assert.Nil(t, s.StoreLock(lock))

	assert.EqualError(t, s.StoreLock(lock), "a lock already exist for this transaction")
}

/*
Scenario: Check if a lock exists
	Given a transaction lock created
	When I want to check if a lock exists for this transaction
	Then I get true
*/
func TestContainsLock(t *testing.T) {
	repo := &mockLockRepository{}
	s := LockService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	lock, _ := NewLock(crypto.HashString("hash"), crypto.HashString("addr"), hex.EncodeToString(pub))
	repo.StoreLock(lock)

	ok, err := s.ContainsLock(lock)
	assert.Nil(t, err)
	assert.True(t, ok)
}

/*
Scenario: Remove a lock
	Given a lock created
	When I want to remove it
	Then I get no lock after
*/
func TestRemoveLock(t *testing.T) {
	repo := &mockLockRepository{}
	s := LockService{
		repo: repo,
	}

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	lock, _ := NewLock(crypto.HashString("hash"), crypto.HashString("addr"), hex.EncodeToString(pub))
	repo.StoreLock(lock)

	assert.Nil(t, s.RemoveLock(lock))
	assert.Len(t, repo.locks, 0)
}

type mockLockRepository struct {
	locks []Lock
}

func (r *mockLockRepository) StoreLock(l Lock) error {
	r.locks = append(r.locks, l)
	return nil
}

func (r *mockLockRepository) RemoveLock(l Lock) error {
	pos := r.findLockPosition(l)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}

func (r mockLockRepository) ContainsLock(l Lock) (bool, error) {
	return r.findLockPosition(l) > -1, nil
}

func (r mockLockRepository) findLockPosition(l Lock) int {
	for i, lock := range r.locks {
		if lock.TransactionHash() == l.TransactionHash() && l.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}
