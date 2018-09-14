package adapters

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
)

//InMemoryPeerRepository handles the repository queries and commands
type InMemoryPeerRepository struct {
	peers []domain.Peer
	seeds []domain.Peer
}

//ListPeers retrieves the peers
func (r InMemoryPeerRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

//InsertPeer executes command insert into the db
func (r *InMemoryPeerRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

//UpdatePeer executes command to update the db
func (r *InMemoryPeerRepository) UpdatePeer(p domain.Peer) error {
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

//GetOwnedPeer retrieves the peer flagged as IsOwned
func (r InMemoryPeerRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}

//ContainsPeer checks if a peer exist in the repo
func (r InMemoryPeerRepository) ContainsPeer(p domain.Peer) (bool, error) {
	cMapped := make(map[string]domain.Peer, 0)
	for _, peer := range r.peers {
		cMapped[hex.EncodeToString(peer.PublicKey)] = peer
	}
	_, exist := cMapped[hex.EncodeToString(p.PublicKey)]
	return exist, nil
}

//ListSeeds retrieves in memory the list of loaded seeds
func (r InMemoryPeerRepository) ListSeeds() ([]domain.Peer, error) {
	return r.seeds, nil
}

//StoreSeed stores in memory a seed
func (r *InMemoryPeerRepository) StoreSeed(p domain.Peer) error {
	r.seeds = append(r.seeds, p)
	return nil
}
