package rabbitmq

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type notifier struct{}

//Notify a new peer has been discovered
func (n notifier) Notify(p discovery.Peer) {
	log.Printf("New peer discovered %s", p.GetEndpoint())

	//TODO connect with rabbitmq
}

//NewNotifier creates a gossip notifier instance using RabbitMQ
func NewNotifier() gossip.Notifier {
	return notifier{}
}
