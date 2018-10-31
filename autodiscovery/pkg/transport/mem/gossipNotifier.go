package mem

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type notifier struct{}

func (n notifier) NotifyDiscoveries(p discovery.Peer) error {
	log.Printf("New peer discovered %s", p.Endpoint())
	return nil
}

func (n notifier) NotifyReachable(pubk string) error {
	log.Printf("New reachable Peer %s", pubk)
	return nil
}

func (n notifier) NotifyUnreachable(pubk string) error {
	log.Printf("New Unreachable Peer %s", pubk)
	return nil
}

//NewNotifier  creates a notifier in memory
func NewNotifier() gossip.Notifier {
	return notifier{}
}
