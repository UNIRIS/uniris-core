package mem

import (
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type notifier struct{}

func (n notifier) Notify(p discovery.Peer) error {
	log.Printf("New peer discovered %s", p.Endpoint())
	return nil
}

//NewNotifier  creates a notifier in memory
func NewNotifier() gossip.Notifier {
	return notifier{}
}
