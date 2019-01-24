package gossip

import (
	uniris "github.com/uniris/uniris-core/pkg"
)

func getUnknownPeers(knownPeers []uniris.Peer, comparePP []uniris.Peer) []uniris.Peer {
	mapPeers := mapPeerSlice(knownPeers)

	diff := make([]uniris.Peer, 0)

	for _, p := range comparePP {

		//Checks if the compared peer is include inside the repository
		kp, exist := mapPeers[p.Identity().PublicKey()]

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

func getNewPeers(knownPeers []uniris.Peer, comparePP []uniris.Peer) []uniris.Peer {
	mapComparee := mapPeerSlice(comparePP)

	diff := make([]uniris.Peer, 0)

	for _, p := range knownPeers {

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

func mapPeerSlice(pp []uniris.Peer) map[string]uniris.Peer {
	mPeers := make(map[string]uniris.Peer)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
