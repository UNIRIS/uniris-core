package bootstraping

import (
	"encoding/hex"
	"net"
	"testing"

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
	repo := new(mockRepository)

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

	repo := new(mockRepository)
	pos := new(mockPositioner)
	net := new(networker)

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

//////////////////////////////////////////////////////////
// 						MOCKS
/////////////////////////////////////////////////////////

type mockRepository struct {
	ownedPeer       discovery.Peer
	discoveredPeers []discovery.Peer
	seedPeers       []discovery.Seed
}

func (r *mockRepository) CountDiscoveredPeers() (int, error) {
	return len(r.discoveredPeers), nil
}

//GetOwnedPeer return the local peer
func (r *mockRepository) GetOwnedPeer() (discovery.Peer, error) {
	return r.ownedPeer, nil
}

//ListSeedPeers return all the seed on the mockRepository
func (r *mockRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seedPeers, nil
}

//ListDiscoveredPeers returns all the discoveredPeers on the mockRepository
func (r *mockRepository) ListDiscoveredPeers() ([]discovery.Peer, error) {
	return r.discoveredPeers, nil
}

func (r *mockRepository) SetPeer(peer discovery.Peer) error {
	if peer.Owned() {
		r.ownedPeer = peer
		return nil
	}
	if r.containsPeer(peer) {
		for _, p := range r.discoveredPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.discoveredPeers = append(r.discoveredPeers, peer)
	}
	return nil
}

func (r *mockRepository) SetSeed(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//GetPeerByIP get a peer from the mockRepository using its ip
func (r *mockRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	if r.ownedPeer.Identity().IP().Equal(ip) {
		return r.ownedPeer, nil
	}
	for i := 0; i < len(r.discoveredPeers); i++ {
		if r.discoveredPeers[i].Identity().IP().Equal(ip) {
			return r.discoveredPeers[i], nil
		}
	}
	return
}

func (r *mockRepository) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.discoveredPeers {
		mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}

type networker struct{}

func (n networker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n networker) CheckInternetState() error {
	return nil
}

func (n networker) CheckNtpState() error {
	return nil
}

type mockPositioner struct{}

func (l mockPositioner) Position() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}
