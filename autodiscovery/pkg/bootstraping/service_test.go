package bootstraping

import (
	"net"
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Loads initial seed peers
	Given list of seeds
	When we want to load them
	Then we can retreive them into the mockRepository
*/
func TestLoadSeeds(t *testing.T) {
	seeds := []discovery.Seed{discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}}
	repo := new(mock.Repository)

	srv := NewService(repo, nil, nil)
	err := srv.LoadSeeds(seeds)
	assert.Nil(t, err)

	rSeeds, _ := repo.ListSeedPeers()
	assert.NotEmpty(t, rSeeds)
	assert.Equal(t, 1, len(rSeeds))
	assert.Equal(t, "127.0.0.1", rSeeds[0].IP.String())
}

/*
Scenario: Starts a peer
	Given a peer mockRepository and a peer localizer
	When a peer startups
	Then the peer is stored on the peer mockRepository
*/
func TestStartup(t *testing.T) {

	repo := new(mock.Repository)
	pos := new(mock.Positioner)
	net := new(mock.Networker)

	srv := NewService(repo, pos, net)
	p, err := srv.Startup([]byte("key"), 3000, "1.0")
	assert.NotNil(t, p)
	assert.Nil(t, err)

	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, "key", p.Identity().PublicKey().String())
	assert.Equal(t, "1.0", p.AppState().Version())
	assert.Equal(t, discovery.BootstrapingStatus, p.AppState().Status())

	selfPeer, _ := repo.GetOwnedPeer()
	assert.Equal(t, "127.0.0.1", selfPeer.Identity().IP().String())
}
