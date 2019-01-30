package memstorage

import "github.com/uniris/uniris-core/pkg/discovery"

type discoveryDb struct {
	knownPeers       []discovery.Peer
	seedPeers        []discovery.PeerIdentity
	unreachablePeers []string
}

//NewDiscoveryDatabase creates a new memory database
func NewDiscoveryDatabase() discovery.Repository {
	return &discoveryDb{}
}

func (db discoveryDb) CountKnownPeers() (int, error) {
	return len(db.knownPeers), nil
}

func (db discoveryDb) ListSeedPeers() ([]discovery.PeerIdentity, error) {
	return db.seedPeers, nil
}

func (db discoveryDb) ListKnownPeers() ([]discovery.Peer, error) {
	return db.knownPeers, nil
}

func (db *discoveryDb) StoreKnownPeer(peer discovery.Peer) error {
	for i, p := range db.knownPeers {
		if p.Identity().PublicKey() == peer.Identity().PublicKey() {
			db.knownPeers[i] = peer
			return nil
		}
	}
	db.knownPeers = append(db.knownPeers, peer)
	return nil
}

func (db discoveryDb) ListReachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)
	for i := 0; i < len(db.knownPeers); i++ {
		if !db.ContainsUnreachablePeer(db.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, db.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

//ListunreachablePeers returns all unreachable peers
func (db discoveryDb) ListUnreachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)
	for i := 0; i < len(db.knownPeers); i++ {
		if db.ContainsUnreachablePeer(db.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, db.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

func (db *discoveryDb) StoreSeedPeer(s discovery.PeerIdentity) error {
	db.seedPeers = append(db.seedPeers, s)
	return nil
}

func (db *discoveryDb) StoreUnreachablePeer(pk string) error {
	if !db.ContainsUnreachablePeer(pk) {
		db.unreachablePeers = append(db.unreachablePeers, pk)
	}
	return nil
}

func (db *discoveryDb) RemoveUnreachablePeer(pk string) error {
	if db.ContainsUnreachablePeer(pk) {
		for i := 0; i < len(db.unreachablePeers); i++ {
			if db.unreachablePeers[i] == pk {
				db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (db discoveryDb) ContainsUnreachablePeer(peerPubK string) bool {
	for _, up := range db.unreachablePeers {
		if up == peerPubK {
			return true
		}
	}
	return false
}
