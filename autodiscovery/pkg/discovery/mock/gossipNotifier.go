package mock

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/gossip"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
)

type MockGossipNotifier struct{}

func NewGossipNotifier() gossip.GossipNotifier {
	return MockGossipNotifier{}
}

func (n MockGossipNotifier) Notify(p discovery.Peer) error {
	log.Print("New peer discovered %s", p.Endpoint())
	return nil
}
