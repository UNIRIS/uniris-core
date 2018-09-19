package rabbitmq

import (
	"log"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type GossipNotifier struct{}

func (n GossipNotifier) Notify(p discovery.Peer) {
	log.Printf("New peer discovered %s", p.GetEndpoint())

	//TODO connect with rabbitmq
}
