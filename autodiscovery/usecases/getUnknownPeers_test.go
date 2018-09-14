package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Get the receiver unknown peers
	Given a list of known peers
	When we received a list of new peers
	Then we returns the list of node the receiver does not known
*/
func TestGetUnknownPeers(t *testing.T) {
	repo := new(GetUnknownPeersTestRepository)
	repo.InsertPeer(domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, false))

	receivedPeers := make([]domain.Peer, 0)
	receivedPeers = append(receivedPeers, domain.NewPeer([]byte("my new key"), net.ParseIP("10.0.0.0"), 3545, false))

	unknownPeers, err := GetUnknownPeers(repo, receivedPeers)
	assert.Nil(t, err)
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, "my new key", string(unknownPeers[0].PublicKey))

}

//=========================
//INTERFACE IMPLEMENTATIONS
//=========================

type GetUnknownPeersTestRepository struct {
	peers []domain.Peer
}

func (r GetUnknownPeersTestRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *GetUnknownPeersTestRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *GetUnknownPeersTestRepository) UpdatePeer(p domain.Peer) error {
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

func (r GetUnknownPeersTestRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}
