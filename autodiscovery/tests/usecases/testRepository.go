package tests

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//PeerRepository implements the IPeerRepository on memory
type PeerRepository struct {
	peerStorage []entities.Peer
}

//GetPeers retrieves the peers stored locally as map
func (ps PeerRepository) GetPeers() (map[string]entities.Peer, error) {
	mapPeers := make(map[string]entities.Peer)

	for _, peer := range ps.peerStorage {
		mapPeers[peer.IP.String()] = peer
	}
	return mapPeers, nil
}

//AddPeer stores a peer locally
func (ps *PeerRepository) AddPeer(p entities.Peer) error {
	ps.peerStorage = append(ps.peerStorage, p)
	return nil
}

//UpdatePeer changes a peer locally
func (ps *PeerRepository) UpdatePeer(newPeer entities.Peer) error {
	newPeers := make([]entities.Peer, 0)
	for _, peer := range ps.peerStorage {
		if peer.IP.Equal(newPeer.IP) {
			newPeers = append(newPeers, newPeer)
		} else {
			newPeers = append(newPeers, peer)
		}
	}
	ps.peerStorage = newPeers
	return nil
}
