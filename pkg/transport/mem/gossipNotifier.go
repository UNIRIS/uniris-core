package memtransport

import (
	"log"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/gossip"
)

type gossipNotif struct{}

//NewGossipNotifier creates a new gossip notifier in memory
func NewGossipNotifier() gossip.Notifier {
	return gossipNotif{}
}

func (n gossipNotif) NotifyDiscovery(p uniris.Peer) error {
	log.Printf("New peer discovered %s", p.String())
	return nil
}

func (n gossipNotif) NotifyReachable(p uniris.Peer) error {
	log.Printf("Peer reached %s", p.String())
	return nil
}

func (n gossipNotif) NotifyUnreachable(p uniris.Peer) error {
	log.Printf("Peer unreached %s", p.String())
	return nil
}
