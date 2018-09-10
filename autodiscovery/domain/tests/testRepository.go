package tests

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//PeerRepository implements the IPeerRepository on memory
type PeerRepository struct {
	peerStorage []*entities.Peer
}

//ListPeer retrieves the peers locally as array
func (ps PeerRepository) ListPeers() ([]*entities.Peer, error) {
	return ps.peerStorage, nil
}

//GetPeers retrieves the peers stored locally as map
func (ps PeerRepository) GetPeers() (map[string]*entities.Peer, error) {
	mapPeers := make(map[string]*entities.Peer)

	for _, peer := range ps.peerStorage {
		mapPeers[string(peer.PublicKey)] = peer
	}
	return mapPeers, nil
}

//AddPeer stores a peer locally
func (ps *PeerRepository) AddPeer(p *entities.Peer) error {
	ps.peerStorage = append(ps.peerStorage, p)
	return nil
}

//UpdatePeer changes a peer locally
func (ps *PeerRepository) UpdatePeer(newPeer *entities.Peer) error {
	newPeers := make([]*entities.Peer, 0)
	for _, peer := range ps.peerStorage {
		if string(peer.PublicKey) == string(newPeer.PublicKey) {
			newPeers = append(newPeers, newPeer)
		} else {
			newPeers = append(newPeers, peer)
		}
	}
	ps.peerStorage = newPeers
	return nil
}
