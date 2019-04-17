package memtransport

import (
	"github.com/uniris/uniris-core/pkg/crypto"
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
func (n DiscoveryNotifier) NotifyReachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.Logger.Info("Peer reached " + string(p))
	return nil
}

//NotifyUnreachable notifies the peer's public key which became unreachable
func (n DiscoveryNotifier) NotifyUnreachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.Logger.Info("Peer unreached " + string(p))
	return nil
}
