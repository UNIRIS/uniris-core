package rabbitmq

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type notifier struct{}

//Notify a new peer has been discovered
func (n notifier) Notify(p discovery.Peer) {
	log.Printf("Peer discovered, %s - beats: %d", p.Endpoint(), p.HeartbeatState().ElapsedHeartbeats())

	//TODO connect with rabbitmq
}

//NewNotifier creates an rabbitmq implementation of the gossip Notifier interface
func NewNotifier() gossip.Notifier {
	return notifier{}
}
