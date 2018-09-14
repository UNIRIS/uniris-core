package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Starts a peer
	Given a need peer
	When the peer starting up
	Then the peer is stored in the repository
*/
func TestStartPeer(t *testing.T) {
	repo := new(StartPeerTestPeerRepository)
	geo := new(StartPeerMockGeo)

	err := StartPeer(repo, geo, domain.NewPeerConfiguration("1.0", []byte("key"), 3545, 1))
	assert.Nil(t, err)

	peers, err := repo.ListPeers()
	assert.Nil(t, err)
	assert.NotEmpty(t, peers)
	assert.Equal(t, "key", string(peers[0].PublicKey))
	assert.Equal(t, "127.0.0.1", peers[0].IP.String())
	assert.Equal(t, "1.0", peers[0].State.Version)
	assert.Equal(t, 3545, peers[0].Port)
	assert.Equal(t, 1, peers[0].State.P2PFactor)
}

//=========================
//INTERFACE IMPLEMENTATIONS
//=========================

type StartPeerMockGeo struct{}

func (geo StartPeerMockGeo) Lookup() (domain.GeoPosition, error) {
	return domain.GeoPosition{Lat: 10, Lon: 50, IP: net.ParseIP("127.0.0.1")}, nil
}

type StartPeerTestPeerRepository struct {
	peers []domain.Peer
}

func (r StartPeerTestPeerRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *StartPeerTestPeerRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *StartPeerTestPeerRepository) UpdatePeer(p domain.Peer) error {
	newPeers := make([]domain.Peer, 0)
	for _, peer := range r.peers {
		if peer.Equals(p) {
			newPeers = append(newPeers, p)
		} else {
			newPeers = append(newPeers, peer)
		}
	}
	r.peers = newPeers
	return nil
}

func (r StartPeerTestPeerRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}
