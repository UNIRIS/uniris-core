package shared

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Store a new shared emitter keypair
	Given a shared key pair
	When I want to store as an emitter keypair
	Then I the keypair is stored
*/
func TestStoreSharedEmitterKeys(t *testing.T) {
	repo := &mockRepo{}
	s := NewService(repo)

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	kp, err := NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	assert.Nil(t, err)
	assert.Nil(t, s.StoreSharedEmitterKeyPair(kp))
	assert.Len(t, repo.emKeys, 1)
	log.Print(repo.emKeys[0])
	assert.Equal(t, hex.EncodeToString([]byte("encpvkey")), repo.emKeys[0].EncryptedPrivateKey())
}

/*
Scenario: List shared emitter keypair
	Given a shared key pair stored
	When I want to get the list of the emitter keypairs
	Then I get all the keypairs
*/
func TestListSharedEmitterKeys(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	kp, err := NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))

	repo := &mockRepo{
		emKeys: EmitterKeys{kp},
	}
	s := NewService(repo)

	assert.Nil(t, err)
	emKeys, err := s.ListSharedEmitterKeyPairs()
	assert.Nil(t, err)
	assert.Len(t, emKeys, 1)
	log.Print(emKeys[0])
	assert.Equal(t, hex.EncodeToString([]byte("encpvkey")), emKeys[0].EncryptedPrivateKey())
}

type mockRepo struct {
	emKeys    EmitterKeys
	minerKeys MinerKeyPair
}

func (r mockRepo) ListSharedEmitterKeyPairs() (EmitterKeys, error) {
	return r.emKeys, nil
}
func (r *mockRepo) StoreSharedEmitterKeyPair(kp EmitterKeyPair) error {
	r.emKeys = append(r.emKeys, kp)
	return nil
}

func (r *mockRepo) GetLastSharedMinersKeyPair() (MinerKeyPair, error) {
	return r.minerKeys, nil
}
