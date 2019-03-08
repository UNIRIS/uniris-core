package memtransport

import (
	"log"

	"github.com/uniris/uniris-core/pkg/discovery"
)

//DiscoveryNotifier is a discovery notifier in memory
type DiscoveryNotifier struct {
	discovery.Notifier
}

//NotifyDiscovery notifies the peer's which has been discovered
func (n DiscoveryNotifier) NotifyDiscovery(p discovery.Peer) error {
	log.Printf("New peer discovered %s", p.String())
	return nil
}

//NotifyReachable notifies the peer's public key which became reachable
func (n DiscoveryNotifier) NotifyReachable(pk string) error {
	log.Printf("Peer reached %s", pk)
	return nil
}

//NotifyUnreachable notifies the peer's public key which became unreachable
func (n DiscoveryNotifier) NotifyUnreachable(pk string) error {
	log.Printf("Peer unreached %s", pk)
	return nil
}
