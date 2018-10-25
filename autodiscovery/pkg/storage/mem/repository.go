package mem

import (
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	gossip "github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

type repo struct {
	seedPeers         []discovery.Seed
	KnownPeers        []discovery.Peer
	UnreacheablePeers []string
}

//NewRepository implements the repository in memory
func NewRepository() discovery.Repository {
	return &repo{}
}

//CountKnownPeers return the number of Known peers
func (r *repo) CountKnownPeers() (int, error) {
	return len(r.KnownPeers), nil
}

//GetOwnedPeer return the local peer
func (r *repo) GetOwnedPeer() (discovery.Peer, error) {
	for _, p := range r.KnownPeers {
		if p.Owned() {
			return p, nil
		}
	}
	return nil, nil
}

//ListSeedPeers return all the seed on the repository
func (r *repo) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seedPeers, nil
}

//ListKnownPeers returns all the discoveredPeers on the repository
func (r *repo) ListKnownPeers() ([]discovery.Peer, error) {
	return r.KnownPeers, nil
}

//SetKnownPeer add a peer to the repository
func (r *repo) SetKnownPeer(peer discovery.Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.KnownPeers {
			if p.Identity().PublicKey() == peer.Identity().PublicKey() {
				p = peer
				break
			}
		}
	} else {
		r.KnownPeers = append(r.KnownPeers, peer)
	}
	return nil
}

//ListReachablePeers returns all the reachable peers on the repository
func (r *repo) ListReachablePeers() ([]discovery.Peer, error) {
	pp := make([]discovery.Peer, 0)
	for i := 0; i < len(r.KnownPeers); i++ {
		if !r.containsUnreachablePeer(r.KnownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.KnownPeers[i])
		}
	}
	return pp, nil
}

//ListUnreacheablePeers returns all unreachable peers on the repository
func (r *repo) ListUnreachablePeers() ([]discovery.Peer, error) {
	pp := make([]discovery.Peer, 0)
	for i := 0; i < len(r.KnownPeers); i++ {
		if r.containsUnreachablePeer(r.KnownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.KnownPeers[i])
		}
	}
	return pp, nil
}

//SetSeedPeer add a seed to the repository
func (r *repo) SetSeedPeer(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//SetUnreachablePeer add an unreachable peer to the repository
func (r *repo) SetUnreachablePeer(pk string) error {
	if !r.containsUnreachablePeer(pk) {
		r.UnreacheablePeers = append(r.UnreacheablePeers, pk)
	}
	return nil
}

//RemoveUnreachablePeer remove an unreachable peer to the repository
func (r *repo) RemoveUnreachablePeer(pk string) error {
	if r.containsUnreachablePeer(pk) {
		for i := 0; i < len(r.UnreacheablePeers); i++ {
			if r.UnreacheablePeers[i] == pk {
				r.UnreacheablePeers = r.UnreacheablePeers[:i+copy(r.UnreacheablePeers[i:], r.UnreacheablePeers[i+1:])]
			}
		}
	}
	return nil
}

//GetPeerByIP get a peer from the repository using its ip
func (r *repo) GetKnownPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.KnownPeers); i++ {
		if r.KnownPeers[i].Identity().IP().Equal(ip) {
			return r.KnownPeers[i], nil
		}
	}
	return
}

//ContainsUnreachableKey check if the pubk is in the list of unreacheable keys
func (r *repo) ContainsUnreachableKey(pubk string) error {
	if r.containsUnreachablePeer(pubk) {
		return nil
	}

	return gossip.ErrNotFoundOnUnreachableList
}

func (r *repo) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.KnownPeers {
		mdiscoveredPeers[p.Identity().PublicKey()] = p
	}

	_, exist := mdiscoveredPeers[p.Identity().PublicKey()]
	return exist
}

func (r *repo) containsUnreachablePeer(pk string) bool {
	for _, up := range r.UnreacheablePeers {
		if up == pk {
			return true
		}
	}
	return false
}
