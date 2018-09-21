package rabbitmq

import (
	"log"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type notifier struct{}

//DisptachNewPeer notifies a new peers has been discovered
func (n notifier) DisptachNewPeer(p discovery.Peer) {
	log.Printf("New peer discovered %s", p.GetEndpoint())

	//TODO connect with rabbitmq
}

//NewNotifier creates an rabbitmq implementation of the gossip Notifier interface
func NewNotifier() discovery.GossipRoundNotifier {
	return notifier{}
}
