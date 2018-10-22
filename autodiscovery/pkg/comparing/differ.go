package comparing

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerDiffer describes methods to perform a diff on a list of peers comparing them by a source of peer
type PeerDiffer interface {
	UnknownPeers([]discovery.Peer) []discovery.Peer
	ProvidePeers([]discovery.Peer) []discovery.Peer
}

type peerDiffer struct {
	source []discovery.Peer
}

//NewPeerDiffer creates a PeerDiffer
func NewPeerDiffer(source []discovery.Peer) PeerDiffer {
	return peerDiffer{source}
}

//UnknownPeers returns peers the source does not contains or outdated
func (d peerDiffer) UnknownPeers(comparePP []discovery.Peer) []discovery.Peer {
	mapRepo := d.mapSlice(d.source)

	diff := make([]discovery.Peer, 0)

	for _, p := range comparePP {

		//Checks if the compared peer is include inside the repository
		kp, exist := mapRepo[p.Identity().PublicKey()]

		if !exist {
			//Adds to the list if the peer is unknown
			diff = append(diff, p)
		} else if p.HeartbeatState().MoreRecentThan(kp.HeartbeatState()) {
			//Adds to the list if the peer is more recent
			diff = append(diff, p)
		}
	}

	return diff
}

//ProvidePeers returns peers the comparee list does not contains or outdated
func (d peerDiffer) ProvidePeers(comparePP []discovery.Peer) []discovery.Peer {
	mapComparee := d.mapSlice(comparePP)

	diff := make([]discovery.Peer, 0)

	for _, p := range d.source {

		//Checks if the known peer is include inside the list of compared peer
		c, exist := mapComparee[p.Identity().PublicKey()]

		if !exist {
			//Adds to the list if the peer is unknown
			diff = append(diff, p)
		} else if p.HeartbeatState().MoreRecentThan(c.HeartbeatState()) {
			//Adds to the list if the peer is more recent
			diff = append(diff, p)
		}
	}

	return diff
}

func (d peerDiffer) mapSlice(pp []discovery.Peer) map[string]discovery.Peer {
	mPeers := make(map[string]discovery.Peer)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
