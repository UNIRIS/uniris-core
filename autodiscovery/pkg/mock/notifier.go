package mock

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

type Notifier struct {
	notifiedPeers []discovery.Peer
}

func (n Notifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

func (n *Notifier) Notify(p discovery.Peer) {
	n.notifiedPeers = append(n.notifiedPeers, p)
}
