package tests

import (
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Get the unknown peers by comparing a peer's list with the known peer's list
	Given a peer's list
	When we compare with our local known peer's list
	Then we get the only the unknown peers
*/
func TestGetUnknownPeers(t *testing.T) {
	repo := GetRepo()
	err := usecases.LoadSeedPeers(&SeedLoader{}, repo)
	assert.Nil(t, err)
	peer := usecases.CreateNewPeer(GetSecondValidPublicKey(), "20.10.200.10")
	unknownPeers, err := usecases.GetUnknownPeers(repo, []*entities.Peer{peer})
	assert.Nil(t, err)
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, GetValidPublicKey(), unknownPeers[0].PublicKey)
	assert.Equal(t, "127.0.0.1", unknownPeers[0].IP.String())
}
