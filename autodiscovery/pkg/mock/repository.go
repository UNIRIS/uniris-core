package mock

import (
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Repository implements the repository interface as mock
type Repository struct {
	SeedPeers        []discovery.Seed
	KnownPeers       []discovery.Peer
	UnreachablePeers []string
}

//CountKnownPeers counts the known peers
func (r *Repository) CountKnownPeers() (int, error) {
	return len(r.KnownPeers), nil
}

//GetOwnedPeer return the local peer
func (r *Repository) GetOwnedPeer() (discovery.Peer, error) {
	for _, p := range r.KnownPeers {
		if p.Owned() {
			return p, nil
		}
	}
	return nil, nil
}

//ListSeedPeers return all the seed on the repository
func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.SeedPeers, nil
}

//ListKnownPeers returns all the KnownPeers on the repository
func (r *Repository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.KnownPeers, nil
}

//SetKnownPeer add or update a known peer
func (r *Repository) SetKnownPeer(peer discovery.Peer) error {
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
func (r *Repository) ListReachablePeers() ([]discovery.Peer, error) {
	pp := make([]discovery.Peer, 0)
	for i := 0; i < len(r.KnownPeers); i++ {
		if !r.containsUnreachablePeer(r.KnownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.KnownPeers[i])
		}
	}
	return pp, nil
}

//ListUnreachablePeers returns all unreachable peers on the repository
func (r *Repository) ListUnreachablePeers() ([]discovery.Peer, error) {
	pp := make([]discovery.Peer, 0)

	for i := 0; i < len(r.SeedPeers); i++ {
		if r.containsUnreachablePeer(r.SeedPeers[i].PublicKey) {
			pp = append(pp, r.SeedPeers[i].AsPeer())
		}
	}

	for i := 0; i < len(r.KnownPeers); i++ {
		if r.containsUnreachablePeer(r.KnownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.KnownPeers[i])
		}
	}
	return pp, nil
}

//SetSeedPeer adds a seed
func (r *Repository) SetSeedPeer(s discovery.Seed) error {
	r.SeedPeers = append(r.SeedPeers, s)
	return nil
}

//SetUnreachablePeer add an unreachable peer to the repository
func (r *Repository) SetUnreachablePeer(pk string) error {
	if !r.containsUnreachablePeer(pk) {
		r.UnreachablePeers = append(r.UnreachablePeers, pk)
	}
	return nil
}

//RemoveUnreachablePeer remove an unreachable peer to the repository
func (r *Repository) RemoveUnreachablePeer(pk string) error {
	if r.containsUnreachablePeer(pk) {
		for i := 0; i < len(r.UnreachablePeers); i++ {
			if r.UnreachablePeers[i] == pk {
				r.UnreachablePeers = r.UnreachablePeers[:i+copy(r.UnreachablePeers[i:], r.UnreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

//GetKnownPeerByIP get a peer from the repository using its ip
func (r *Repository) GetKnownPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.KnownPeers); i++ {
		if r.KnownPeers[i].Identity().IP().Equal(ip) {
			return r.KnownPeers[i], nil
		}
	}
	return
}

func (r *Repository) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.KnownPeers {
		mdiscoveredPeers[p.Identity().PublicKey()] = p
	}

	_, exist := mdiscoveredPeers[p.Identity().PublicKey()]
	return exist
}

func (r *Repository) containsUnreachablePeer(pk string) bool {
	for _, up := range r.UnreachablePeers {
		if up == pk {
			return true
		}
	}
	return false
}
