package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Get the sender unknown peers
	Given a list of known peers
	When we received a list of new peers
	Then we returns the list of node the sender does not known
*/
func TestProvideNewPeers(t *testing.T) {
	repo := new(ProvideNewPeersTestRepository)

	repo.InsertPeer(domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, true))
	repo.InsertPeer(domain.NewPeer([]byte("other key"), net.ParseIP("30.0.0.1"), 3545, false))

	receivedPeers := make([]domain.Peer, 0)
	receivedPeers = append(receivedPeers, domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, true))

	newPeers, err := ProvideNewPeers(repo, receivedPeers)
	assert.Nil(t, err)
	assert.NotEmpty(t, newPeers)
}

//=========================
//INTERFACE IMPLEMENTATIONS
//=========================

type ProvideNewPeersTestRepository struct {
	peers []domain.Peer
}

func (r ProvideNewPeersTestRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *ProvideNewPeersTestRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *ProvideNewPeersTestRepository) UpdatePeer(p domain.Peer) error {
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

func (r ProvideNewPeersTestRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}
