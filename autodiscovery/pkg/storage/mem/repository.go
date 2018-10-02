package mem

import (
	"encoding/hex"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type repo struct {
	seedPeers         []discovery.Seed
	KnownPeers        []discovery.Peer
	UnreacheablePeers []discovery.PublicKey
}

//NewRepository implements the repository in memory
func NewRepository() discovery.Repository {
	return &repo{}
}

//CountKnownPeers
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

func (r *repo) SetKnownPeer(peer discovery.Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.KnownPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
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

func (r *repo) SetSeedPeer(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//SetUnreachablePeer add an unreachable peer to the repository
func (r *repo) SetUnreachablePeer(pk discovery.PublicKey) error {
	if !r.containsUnreachablePeer(pk) {
		r.UnreacheablePeers = append(r.UnreacheablePeers, pk)
	}
	return nil
}

//DelUnreacheablePeer remove an unreachable peer to the repository
func (r *repo) RemoveUnreachablePeer(pk discovery.PublicKey) error {
	if r.containsUnreachablePeer(pk) {
		for i := 0; i < len(r.UnreacheablePeers); i++ {
			if r.UnreacheablePeers[i].Equals(pk) {
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

func (r *repo) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.KnownPeers {
		mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}

func (r *repo) containsUnreachablePeer(pk discovery.PublicKey) bool {
	for _, up := range r.UnreacheablePeers {
		if up.Equals(pk) {
			return true
		}
	}
	return false
}
