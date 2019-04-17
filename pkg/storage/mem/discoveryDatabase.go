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

func (db *discoveryDB) DiscoveredPeers() ([]discovery.Peer, error) {
	return db.discoveredPeers, nil
}

func (db *discoveryDB) WriteDiscoveredPeer(peer discovery.Peer) error {
	if db.containsPeer(peer) {
		for _, p := range db.discoveredPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		db.discoveredPeers = append(db.discoveredPeers, peer)
	}
	return nil
}

func (db *discoveryDB) UnreachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)
	for i := 0; i < len(db.discoveredPeers); i++ {
		if exist, _ := db.ContainsUnreachablePeer(db.discoveredPeers[i].Identity()); exist {
			pp = append(pp, db.discoveredPeers[i].Identity())
		}
	}
	return pp, nil
}

func (db *discoveryDB) WriteUnreachablePeer(pi discovery.PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pi); !exist {
		db.unreachablePeers = append(db.unreachablePeers, pi)
	}
	return nil
}

func (db *discoveryDB) RemoveUnreachablePeer(pi discovery.PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pi); exist {
		for i := 0; i < len(db.unreachablePeers); i++ {
			if db.unreachablePeers[i].PublicKey().Equals(pi.PublicKey()) {
				db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (db *discoveryDB) ContainsUnreachablePeer(pi discovery.PeerIdentity) (bool, error) {
	for _, up := range db.unreachablePeers {
		if up.PublicKey().Equals(pi.PublicKey()) {
			return true, nil
		}
	}
	return false, nil
}

func (db *discoveryDB) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range db.discoveredPeers {
		mdiscoveredPeers[string(p.Identity().PublicKey().Bytes())] = p
	}

	_, exist := mdiscoveredPeers[string(p.Identity().PublicKey().Bytes())]
	return exist
}
