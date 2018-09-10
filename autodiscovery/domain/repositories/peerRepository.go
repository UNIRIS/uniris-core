package repositories

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//PeerRepository represents the interface for peer storage operations
type PeerRepository interface {
	GetPeers() (map[string]*entities.Peer, error)
	ListPeers() ([]*entities.Peer, error)
	AddPeer(p *entities.Peer) error
	UpdatePeer(p *entities.Peer) error
}
