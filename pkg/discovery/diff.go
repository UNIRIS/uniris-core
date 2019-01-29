package discovery

func getUnknownPeers(knownPeers []Peer, comparePP []Peer) []Peer {
	mapPeers := mapPeerSlice(knownPeers)

	diff := make([]Peer, 0)

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

func getNewPeers(knownPeers []Peer, comparePP []Peer) []Peer {
	mapComparee := mapPeerSlice(comparePP)

	diff := make([]Peer, 0)

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

func mapPeerSlice(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
