package ports

import (
	"github.com/uniris/uniris-core/autodiscovery/core/domain"
)

//DiscoveryNotifier notifies when a new peer is discovered
type DiscoveryNotifier interface {
	NotifyNewPeer(p domain.Peer) error
}
