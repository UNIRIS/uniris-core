package bootstraping

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Loads initial seed peers
	Given list of seeds
	When we want to load them
	Then we can retreive them into the repository
*/
func TestLoadSeeds(t *testing.T) {
	seeds := []discovery.Seed{discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}}
	repo := new(mockPeerRepository)

	srv := NewService(repo, nil)
	err := srv.LoadSeeds(seeds)
	assert.Nil(t, err)

	rSeeds, _ := repo.ListSeedPeers()
	assert.NotEmpty(t, rSeeds)
	assert.Equal(t, 1, len(rSeeds))
	assert.Equal(t, "127.0.0.1", rSeeds[0].IP.String())
}

/*
Scenario: Starts a peer
	Given a peer repository and a peer localizer
	When a peer startups
	Then the peer is stored on the peer repository
*/
func TestStartup(t *testing.T) {

	repo := new(mockPeerRepository)
	loc := new(mockPeerLocalizer)

	srv := NewService(repo, loc)
	p, err := srv.Startup([]byte("key"), 3000, 1, "1.0")
	assert.NotNil(t, p)
	assert.Nil(t, err)

	assert.Equal(t, "127.0.0.1", p.IP().String())
	assert.Equal(t, "key", string(p.PublicKey()))
	assert.Equal(t, 1, p.P2PFactor())
	assert.Equal(t, "1.0", p.Version())
	assert.Equal(t, discovery.BootstrapingStatus, p.Status())

	pp, _ := repo.ListKnownPeers()
	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].IP().String())
	assert.True(t, pp[0].IsOwned())
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *mockPeerRepository) AddPeer(p discovery.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *mockPeerRepository) AddSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

func (r *mockPeerRepository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

func (r *mockPeerRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

type mockPeerLocalizer struct{}

func (l mockPeerLocalizer) GetIP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (l mockPeerLocalizer) GetGeoPosition() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}
