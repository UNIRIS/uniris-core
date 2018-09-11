package repositories

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//InMemoryPeerRepository implements the PeerRepository interface in memory
type InMemoryPeerRepository struct {
	peers     []*entities.Peer
	localPeer *entities.Peer
}

//ListSeedPeers retrieves the seed peers
func (ps *InMemoryPeerRepository) ListSeedPeers() ([]*entities.Peer, error) {
	seeds := make([]*entities.Peer, 0)
	for _, peer := range ps.peers {
		if peer.Category == entities.SeedCategory {
			seeds = append(seeds, peer)
		}
	}
	return seeds, nil
}

//ListDiscoveredPeers retrieves the discovered peers
func (ps *InMemoryPeerRepository) ListDiscoveredPeers() ([]*entities.Peer, error) {
	discovered := make([]*entities.Peer, 0)
	for _, peer := range ps.peers {
		if peer.Category == entities.DiscoveredCategory {
			discovered = append(discovered, peer)
		}
	}
	return discovered, nil
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

//SetLocalPeer stores the peer starting
func (ps *InMemoryPeerRepository) SetLocalPeer(p *entities.Peer) error {
	ps.localPeer = p
	return nil
}

//GetLocalPeer returns the local stored peer
func (ps *InMemoryPeerRepository) GetLocalPeer() (*entities.Peer, error) {
	return ps.localPeer, nil
}
