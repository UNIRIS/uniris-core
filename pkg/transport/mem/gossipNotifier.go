package memtransport

import (
	"log"

	"github.com/uniris/uniris-core/pkg/discovery"
)

type gossipNotif struct{}

//NewGossipNotifier creates a new gossip notifier in memory
func NewGossipNotifier() discovery.Notifier {
	return gossipNotif{}
}

func (n gossipNotif) NotifyDiscovery(p discovery.Peer) error {
	log.Printf("New peer discovered %s", p.String())
	return nil
}

func (n gossipNotif) NotifyReachable(p discovery.Peer) error {
	log.Printf("Peer reached %s", p.String())
	return nil
}

func (n gossipNotif) NotifyUnreachable(p discovery.Peer) error {
	log.Printf("Peer unreached %s", p.String())
	return nil
}
