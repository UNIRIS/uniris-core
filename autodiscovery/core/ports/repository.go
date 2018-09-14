package ports

import "github.com/uniris/uniris-core/autodiscovery/core/domain"

//PeerRepository represents the interface for peer storage operations
type PeerRepository interface {
	ListPeers() ([]domain.Peer, error)
	InsertPeer(p domain.Peer) error
	UpdatePeer(p domain.Peer) error
	GetOwnedPeer() (domain.Peer, error)
	ContainsPeer(p domain.Peer) (bool, error)

	StoreSeed(p domain.Peer) error
	ListSeeds() ([]domain.Peer, error)
}
