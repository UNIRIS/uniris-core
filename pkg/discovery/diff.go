package discovery

//ExcludedOrRecent compares a source of peers with an other list of peers
//and returns the peers that are not included inside the source or
func ExcludedOrRecent(source []Peer, comparees []Peer) []Peer {
	mapSource := mapPeers(source)

	diff := make([]Peer, 0)

	for _, p := range comparees {

		//Add it to the list if the compared peer is include inside the source
		sp, exist := mapSource[p.Identity().PublicKey()]
		if !exist {
			diff = append(diff, p)
		} else if p.HeartbeatState().MoreRecentThan(sp.HeartbeatState()) {
			//Add it to the list if the comparee peer is more recent than the founded peer in the source
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
