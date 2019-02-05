package discovery

//GetUnknownPeers compare local peers and incoming peers.
//It returns the peers not included inside the incoming list.
func GetUnknownPeers(localPeers []Peer, comparePP []Peer) []Peer {
	mPeers := mapPeers(localPeers)

	diff := make([]Peer, 0)

	for _, p := range comparePP {

		//Checks if the compared peer is include inside the repository
		kp, exist := mPeers[p.Identity().PublicKey()]

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

//GetNewPeers compare local peers and incoming peers.
//It returns the peers not included inside the local list.
func GetNewPeers(localPeers []Peer, comparePP []Peer) []Peer {
	mPeers := mapPeers(comparePP)

	diff := make([]Peer, 0)

	for _, p := range localPeers {

		//Checks if the known peer is include inside the list of compared peer
		c, exist := mPeers[p.Identity().PublicKey()]

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

func mapPeers(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
