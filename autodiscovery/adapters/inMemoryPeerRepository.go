package adapters

import "github.com/uniris/uniris-core/autodiscovery/domain"

//InMemoryPeerRepository implements the interface of PeerRepository in memory
type InMemoryPeerRepository struct {
	peers []domain.Peer
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
