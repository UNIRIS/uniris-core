package repositories

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//PeerRepository represents the interface for peer storage operations
type PeerRepository interface {
	ListPeers() ([]*entities.Peer, error)
	ListSeedPeers() ([]*entities.Peer, error)
	ListDiscoveredPeers() ([]*entities.Peer, error)

	AddPeer(p *entities.Peer) error
	UpdatePeer(p *entities.Peer) error

	SetLocalPeer(p *entities.Peer) error
	GetLocalPeer() (*entities.Peer, error)
}
