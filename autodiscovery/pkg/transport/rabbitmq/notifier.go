package rabbitmq

import "github.com/uniris/uniris-core/autodiscovery/pkg/discovery"

type GossipNotifier struct{}

func (n GossipNotifier) DispatchDiscovery(p discovery.Peer) error {
	return nil
}
