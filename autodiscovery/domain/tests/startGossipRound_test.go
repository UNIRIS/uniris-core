package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Initiate the start gossip round with the seed
	Given empty peer repo
	When we start to gossip
	Then we discover seed peer
*/
func TestStartGossipRoundWithEmptyRepo(t *testing.T) {
	repo := GetRepo()

	err := usecases.StartGossipRound(&SeedLoader{}, repo, &FullGossipService{})
	assert.Nil(t, err)
	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
	assert.Equal(t, 1, len(peers))
	assert.Equal(t, GetValidPublicKey(), peers[0].PublicKey)
}

/*
Scenario: Initiate the start gossip round with the seed
	Given a filled peer repo
	When we start to gossip
	Then we discover the seed peer
*/
func TestStartGossipRoundWithFilledRepo(t *testing.T) {
	repo := GetRepo()

	usecases.StartupPeer(repo, &GeolocService{}, GetSecondValidPublicKey(), 3545)

	err := usecases.StartGossipRound(&SeedLoader{}, repo, &FullGossipService{})
	assert.Nil(t, err)
	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
	assert.Equal(t, 2, len(peers))
	assert.Equal(t, GetSecondValidPublicKey(), peers[0].PublicKey)
	assert.Equal(t, GetValidPublicKey(), peers[1].PublicKey)
}
