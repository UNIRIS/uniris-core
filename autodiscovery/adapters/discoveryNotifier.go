package adapters

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
)

//InMemoryDiscoveryNotifier notifies when a new peer is discovered
type InMemoryDiscoveryNotifier struct {
}

//NotifyNewPeer notify when a new peer is created
func (n InMemoryDiscoveryNotifier) NotifyNewPeer(p domain.Peer) error {
	log.Printf("New peer discovered %s", p.GetDiscoveryEndpoint())
	return nil
}
