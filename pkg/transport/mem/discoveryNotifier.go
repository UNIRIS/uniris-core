package memtransport

import (
	"log"

	"github.com/uniris/uniris-core/pkg/discovery"
)

//DiscoveryNotifier is a discovery notifier in memory
type DiscoveryNotifier struct {
	discovery.Notifier
}

func (n DiscoveryNotifier) NotifyDiscovery(p discovery.Peer) error {
	log.Printf("New peer discovered %s", p.String())
	return nil
}

func (n DiscoveryNotifier) NotifyReachable(p discovery.PeerIdentity) error {
	log.Printf("Peer reached %s", p.Endpoint())
	return nil
}

func (n DiscoveryNotifier) NotifyUnreachable(p discovery.PeerIdentity) error {
	log.Printf("Peer unreached %s", p.Endpoint())
	return nil
}
