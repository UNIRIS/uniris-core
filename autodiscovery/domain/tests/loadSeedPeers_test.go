package tests

import (
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenerio: Load seed peers
	Given a seed peer file
	When we load the seed peers
	Then we can retrieve the peers on the peer's database
*/
func TestLoadSeedPeersWithExistingFile(t *testing.T) {
	repo := &PeerRepository{}
	err := usecases.LoadSeedPeers(&SeedLoader{}, repo)
	assert.Nil(t, err)
	peers, _ := repo.ListSeedPeers()
	assert.NotEmpty(t, peers, "Seed peers must not be empty")
	assert.Equal(t, "127.0.0.1", peers[0].IP.String(), "Seed peer IP must be 127.0.0.1")
	assert.Equal(t, entities.SeedCategory, peers[0].Category, "Peer category must be seed")
}
