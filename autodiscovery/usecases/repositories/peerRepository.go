package repositories

import "github.com/uniris/uniris-core/autodiscovery/domain"

//PeerRepository represents the interface for peer storage operations
type PeerRepository interface {
	ListPeers() ([]domain.Peer, error)
	InsertPeer(p domain.Peer) error
	UpdatePeer(p domain.Peer) error
	GetOwnedPeer() (domain.Peer, error)
}
