package tests

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenerio: Load seed peers
	Given not seed peer file
	When we load the seed peers
	Then a error is returned saying there is no file
*/
func TestLoadSeedPeersWithMissingFile(t *testing.T) {
	repo := &PeerRepository{}
	err := usecases.LoadSeedPeers(repo, "")
	assert.NotNil(t, err)
	assert.Equal(t, "open : no such file or directory", err.Error())
}

/*
Scenerio: Load seed peers
	Given a seed peer file
	When we load the seed peers
	Then we can retrieve the peers on the peer's database
*/
func TestLoadSeedPeersWithExistingFile(t *testing.T) {
	repo := &PeerRepository{}
	path, _ := filepath.Abs("./seed.json")
	err := usecases.LoadSeedPeers(repo, path)
	assert.Nil(t, err)
	peers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)
	assert.NotEmpty(t, peers, "Seed peers must not be empty")
	assert.Equal(t, "127.0.0.1", peers[0].IP.String(), "Seed peer IP must be 127.0.0.1")
}
