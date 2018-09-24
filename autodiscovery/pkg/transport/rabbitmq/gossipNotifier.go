package rabbitmq

import (
	"log"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

type notifier struct{}

//DisptachNewPeer notifies a new peers has been discovered
func (n notifier) Notify(p discovery.Peer) {
	log.Printf("New peer discovered %s", p.Endpoint())

	//TODO connect with rabbitmq
}

//NewNotifier creates an rabbitmq implementation of the gossip Notifier interface
func NewNotifier() gossip.Notifier {
	return notifier{}
}
