package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Discover the peers from another peer by gossip which known different peers
	Given a node which have different peer's registered
	When we want discover its peers
	Then we load the unknown peers in our local repository
*/
func TestDiscoverNewPeers(t *testing.T) {
	repo := GetRepo()
	err := usecases.LoadSeedPeers(&SeedLoader{}, repo)
	assert.Nil(t, err)
	err = usecases.DiscoverPeers(repo, &FullGossipService{})
	assert.Nil(t, err)
	knownPeers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(knownPeers))
	assert.Equal(t, GetValidPublicKey(), knownPeers[0].PublicKey)
	assert.Equal(t, GetSecondValidPublicKey(), knownPeers[1].PublicKey, "First new peer discovered must have the public key %s", string(GetSecondValidPublicKey()))
	assert.Equal(t, GetThirdValidPublicKey(), knownPeers[2].PublicKey, "Second new peer discovered must have the public key %s", string(GetThirdValidPublicKey()))

}

/*
Scenario: Discover the peers from another peer by gossip which known the same peers
	Given a node which have the same peer's registered
	When we want discover its peers
	Then we don't change our peer's list
*/
func TestDiscoverSamePeers(t *testing.T) {
	repo := GetRepo()
	err := usecases.LoadSeedPeers(&SeedLoader{}, repo)
	assert.Nil(t, err)
	err = usecases.DiscoverPeers(repo, &SameGossipService{})
	assert.Nil(t, err)
	knownPeers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(knownPeers))
	assert.Equal(t, GetValidPublicKey(), knownPeers[0].PublicKey, "Public key of the known peer must be the same")
}
