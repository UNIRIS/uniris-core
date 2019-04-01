package memtransport

import (
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/logging"
)

//DiscoveryNotifier is a discovery notifier in memory
type DiscoveryNotifier struct {
	discovery.Notifier
	Logger logging.Logger
}

//NotifyDiscovery notifies the peer's which has been discovered
func (n DiscoveryNotifier) NotifyDiscovery(p discovery.Peer) error {
	n.Logger.Info("New peer discovered " + p.String())
	return nil
}

//NotifyReachable notifies the peer's public key which became reachable
func (n DiscoveryNotifier) NotifyReachable(pk string) error {
	n.Logger.Info("Peer reached " + pk)
	return nil
}

//NotifyUnreachable notifies the peer's public key which became unreachable
func (n DiscoveryNotifier) NotifyUnreachable(pk string) error {
	n.Logger.Info("Peer unreached " + pk)
	return nil
}
