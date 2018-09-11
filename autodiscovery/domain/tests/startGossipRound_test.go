package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Initiate the start gossip round
	Given a peer startup
	When we start to gossip
	Then we discover unknown peers
*/
func TestStartGossipRound(t *testing.T) {
	repo := GetRepo()
	usecases.LoadSeedPeers(&SeedLoader{}, repo)
	err := usecases.StartGossipRound(repo, &FullGossipService{})
	assert.Nil(t, err)
	discoveredPeers, _ := repo.ListDiscoveredPeers()
	assert.NotEmpty(t, discoveredPeers)
	assert.Equal(t, GetValidPublicKey(), discoveredPeers[0].PublicKey)
}
