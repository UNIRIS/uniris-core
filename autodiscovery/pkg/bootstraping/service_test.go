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
	Then we can retreive them into the repository
*/
func TestLoadSeeds(t *testing.T) {
	seeds := []discovery.Seed{discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}}
	repo := new(mockPeerRepository)

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
	Given a peer repository and a peer localizer
	When a peer startups
	Then the peer is stored on the peer repository
*/
func TestStartup(t *testing.T) {

	repo := new(mockPeerRepository)
	pos := new(mockPeerPositionner)
	net := new(mockPeerNetworker)

	srv := NewService(repo, pos, net)
	p, err := srv.Startup([]byte("key"), 3000, "1.0")
	assert.NotNil(t, p)
	assert.Nil(t, err)

	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, "key", p.Identity().PublicKey().String())
	assert.Equal(t, "1.0", p.AppState().Version())
	assert.Equal(t, discovery.BootstrapingStatus, p.AppState().Status())

	pp, _ := repo.ListKnownPeers()
	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].Identity().IP().String())
	assert.True(t, pp[0].Owned())
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *mockPeerRepository) CountKnownPeers() (int, error) {
	return len(r.peers), nil
}

func (r *mockPeerRepository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.Owned() {
			return p, nil
		}
	}
	return
}

func (r *mockPeerRepository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
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

func (r *mockPeerRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.peers); i++ {
		if string(ip) == string(r.peers[i].IP()) {
			return r.peers[i], nil
		}
	}
	return
}

func (r *mockPeerRepository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.peers {
		if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *mockPeerRepository) containsPeer(peer discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = peer
	}

	_, exist := mPeers[hex.EncodeToString(peer.Identity().PublicKey())]
	return exist
}

type mockPeerPositionner struct{}

func (l mockPeerPositionner) Position() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}

type mockPeerNetworker struct{}

func (n mockPeerNetworker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockPeerNetworker) CheckInternetState() error {
	return nil
}

func (n mockPeerNetworker) CheckNtpState() error {
	return nil
}
