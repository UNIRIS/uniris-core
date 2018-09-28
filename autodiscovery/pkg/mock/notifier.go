package mock

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

//Notifier mock
type Notifier struct {
	notifiedPeers []discovery.Peer
}

//NotifiedPeers mock
func (n Notifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

//Notify mock
func (n *Notifier) Notify(p discovery.Peer) {
	n.notifiedPeers = append(n.notifiedPeers, p)
}
