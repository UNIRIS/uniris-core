package adding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Store a biometric device  public key
	Given a biometric device public key
	When I want to store it
	Then it's stored. If the key already exists, it's replaced
*/
func TestStoreBiodPubKey(t *testing.T) {
	repo := &mockRepository{}
	s := NewService(repo)
	err := s.RegisterKey("my key")
	assert.Nil(t, err)

	assert.NotEmpty(t, repo.BiodPubKeys)
	assert.Equal(t, "my key", repo.BiodPubKeys[0])
}

type mockRepository struct {
	BiodPubKeys []string
}

func (r *mockRepository) StoreBiodPublicKey(key string) error {
	//Prevent to add multiple times
	for _, k := range r.BiodPubKeys {
		if k == key {
			return nil
		}
	}
	r.BiodPubKeys = append(r.BiodPubKeys, key)
	return nil
}
