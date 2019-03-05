package memstorage

import "github.com/uniris/uniris-core/pkg/discovery"

type discoveryDB struct {
	discoveredPeers  []discovery.Peer
	unreachablePeers []discovery.PeerIdentity
}

//NewDiscoveryDatabase creates a new discovery database in memory
func NewDiscoveryDatabase() discovery.Database {
	return &discoveryDB{}
}

func (db discoveryDB) DiscoveredPeers() ([]discovery.Peer, error) {
	return db.discoveredPeers, nil
}

func (db *discoveryDB) WriteDiscoveredPeer(peer discovery.Peer) error {
	for i, p := range db.discoveredPeers {
		if p.Identity().PublicKey() == peer.Identity().PublicKey() {
			db.discoveredPeers[i] = peer
			return nil
		}
	}
	db.discoveredPeers = append(db.discoveredPeers, peer)
	return nil
}

func (db discoveryDB) UnreachablePeers() ([]discovery.PeerIdentity, error) {
	return db.unreachablePeers, nil
}

func (db *discoveryDB) WriteUnreachablePeer(p discovery.PeerIdentity) error {
	db.unreachablePeers = append(db.unreachablePeers, p)
	return nil
}

func (db *discoveryDB) RemoveUnreachablePeer(p discovery.PeerIdentity) error {
	for i := 0; i < len(db.unreachablePeers); i++ {
		if db.unreachablePeers[i].PublicKey() == p.PublicKey() {
			db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
		}
	}
	return nil
}
