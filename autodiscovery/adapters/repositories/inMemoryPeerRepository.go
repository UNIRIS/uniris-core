package repositories

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//InMemoryPeerRepository implements the PeerRepository interface on memory
type InMemoryPeerRepository struct {
	peers []*entities.Peer
}

//GetOwnPeer returns the owned peer
func (ps *InMemoryPeerRepository) GetOwnPeer() (*entities.Peer, error) {
	for _, peer := range ps.peers {
		if peer.IsSelf {
			return peer, nil
		}
	}
	return nil, nil
}

//ListPeers get all peers
func (ps *InMemoryPeerRepository) ListPeers() ([]*entities.Peer, error) {
	return ps.peers, nil
}

//AddPeer stores a new peer
func (ps *InMemoryPeerRepository) AddPeer(peer *entities.Peer) error {
	ps.peers = append(ps.peers, peer)
	return nil
}

//UpdatePeer changes a peer
func (ps *InMemoryPeerRepository) UpdatePeer(newPeer *entities.Peer) error {
	peers := make([]*entities.Peer, 0)
	for _, peer := range ps.peers {
		if string(peer.PublicKey) == string(newPeer.PublicKey) {
			peers = append(peers, newPeer)
		} else {
			peers = append(peers, peer)
		}
	}
	ps.peers = peers
	return nil
}
