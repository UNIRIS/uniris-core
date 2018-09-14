package usecases

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Insert a new peer
	Given a new peer
	When we store the peer
	Then the peer is inserted in the repository
*/
func TestInsertPeer(t *testing.T) {
	repo := new(StorePeerTestPeerRepository)
	err := StorePeer(repo, domain.Peer{PublicKey: []byte("my key")})
	assert.Nil(t, err)
	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
	assert.Equal(t, "my key", string(peers[0].PublicKey))
}

/*
Scenario: Update an existing
	Given a peer
	When we store the peer
	Then the peer is updated in the repository
*/
func TestUpdatePeer(t *testing.T) {
	repo := new(StorePeerTestPeerRepository)
	repo.InsertPeer(domain.Peer{PublicKey: []byte("my key")})

	time.Sleep(100 * time.Millisecond)
	err := StorePeer(repo, domain.Peer{
		PublicKey:      []byte("my key"),
		IP:             net.ParseIP("127.0.0.1"),
		GenerationTime: time.Now(),
	})

	assert.Nil(t, err)
	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
	assert.Equal(t, "127.0.0.1", peers[0].IP.String())
}

type StorePeerTestPeerRepository struct {
	peers []domain.Peer
}

func (r StorePeerTestPeerRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *StorePeerTestPeerRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *StorePeerTestPeerRepository) UpdatePeer(p domain.Peer) error {
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

func (r StorePeerTestPeerRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}
