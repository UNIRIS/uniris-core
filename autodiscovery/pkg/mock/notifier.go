package mock

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

//Notifier mock
type Notifier struct {
	notifiedPeers             []discovery.Peer
	notifiedUnreacheablePeers []string
}

//NotifiedPeers mock
func (n Notifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

//NotifyReachable notifies for an unreachable peer which is now reachable
func (n Notifier) NotifyReachable(pubk string) error {
	for i := 0; i < len(n.notifiedUnreacheablePeers); i++ {
		if n.notifiedUnreacheablePeers[i] == pubk {
			n.notifiedUnreacheablePeers = n.notifiedUnreacheablePeers[:i+copy(n.notifiedUnreacheablePeers[i:], n.notifiedUnreacheablePeers[i+1:])]
		}
	}
	return nil
}

//NotifyUnreachable notifies for an unreachable peer
func (n Notifier) NotifyUnreachable(pubk string) error {
	n.notifiedUnreacheablePeers = append(n.notifiedUnreacheablePeers, pubk)
	return nil
}

//NotifyDiscoveries notifies a new peers has been discovered
func (n *Notifier) NotifyDiscoveries(p discovery.Peer) error {
	n.notifiedPeers = append(n.notifiedPeers, p)
	return nil
}
