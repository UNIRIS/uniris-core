package discovery

//Database wrap discovery database queries, persistence and removals
type Database interface {
	DatabaseReader
	DatabaseWriter
	DatabaseRemover
}

//DatabaseReader handles network peer database query
type DatabaseReader interface {

	//DiscoveredPeers retrieves the discovered peers from the discovery database
	DiscoveredPeers() ([]Peer, error)

	//UnreachablePeers retrieves the unreachables peers from the unreachables database
	UnreachablePeers() ([]PeerIdentity, error)

	//ContainsUnreachablePeer determinates if the peer is found inside the unreachable database
	ContainsUnreachablePeer(p PeerIdentity) (bool, error)
}

//DatabaseWriter handles network peer discovery persistence
type DatabaseWriter interface {

	//WriteDiscoveredPeer inserts or updates the peer in the discovery database
	WriteDiscoveredPeer(p Peer) error

	//WriteUnreachablePeer inserts the peer in the unreachable database
	WriteUnreachablePeer(p PeerIdentity) error
}

//DatabaseRemover handles the removing of the unfreshed data
type DatabaseRemover interface {
	//RemoveUnreachablePeer deletes the peer from the unreachable database
	RemoveUnreachablePeer(p PeerIdentity) error
}

func reachablePeers(db DatabaseReader) ([]Peer, []PeerIdentity, error) {
	peers, err := db.DiscoveredPeers()
	if err != nil {
		return nil, nil, err
	}

	reachables := make([]Peer, 0)
	reachablesIdentity := make([]PeerIdentity, 0)

	for _, p := range peers {
		exist, err := db.ContainsUnreachablePeer(p.Identity())
		if err != nil {
			return nil, nil, err
		}
		if !exist {
			reachables = append(reachables, p)
			reachablesIdentity = append(reachablesIdentity, p.identity)

		}
	}

	return reachables, reachablesIdentity, nil

}
