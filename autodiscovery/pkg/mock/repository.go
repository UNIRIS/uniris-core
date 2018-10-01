package mock

import (
	"encoding/hex"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Repository mock
type Repository struct {
	Peers            []discovery.Peer
	Seeds            []discovery.Seed
	UnreachablePeers []discovery.PublicKey
}

//CountKnownPeers retrun the number of Known peers
func (r *Repository) CountKnownPeers() (int, error) {
	return len(r.Peers), nil
}

//GetOwnedPeer return the local peer
func (r *Repository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.Peers {
		if p.Owned() {
			return p, nil
		}
	}
	return
}

//ListSeedPeers return all the seed on the repository
func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.Seeds, nil
}

//ListKnownPeers returns all the peers on the repository
func (r *Repository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.Peers, nil
}

//ListReachablePeers returns all the reacheable peers on the repository
func (r *Repository) ListReachablePeers() ([]discovery.Peer, error) {
	rp := make([]discovery.Peer, 0)
	for i := 0; i < len(r.Peers); i++ {
		if !r.containsUnreacheablePeer(r.Peers[i].Identity().PublicKey()) {
			rp = append(rp, r.Peers[i])
		}
	}
	return rp, nil
}

//ListUnreachablePeers returns all unreacheable peers on the repository
func (r *Repository) ListUnreachablePeers() ([]discovery.Peer, error) {
	unrp := make([]discovery.Peer, 0)
	for i := 0; i < len(r.Peers); i++ {
		if r.containsUnreacheablePeer(r.Peers[i].Identity().PublicKey()) {
			unrp = append(unrp, r.Peers[i])
		}
	}
	return unrp, nil
}

//AddPeer add a peer to the repository
func (r *Repository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.Peers = append(r.Peers, p)
	return nil
}

//AddSeed add a seed to the repository
func (r *Repository) AddSeed(s discovery.Seed) error {
	r.Seeds = append(r.Seeds, s)
	return nil
}

//AddUnreachablePeer add an unreachable peer to the repository
func (r *Repository) AddUnreachablePeer(pk discovery.PublicKey) error {
	if !r.containsUnreacheablePeer(pk) {
		r.UnreachablePeers = append(r.UnreachablePeers, pk)
	}
	return nil
}

//DelUnreachablePeer remove an unreachable peer to the repository
func (r *Repository) DelUnreachablePeer(pk discovery.PublicKey) error {
	if r.containsUnreacheablePeer(pk) {
		for i := 0; i < len(r.UnreachablePeers); i++ {
			if r.UnreachablePeers[i].Equals(pk) {
				r.UnreachablePeers = r.UnreachablePeers[:i+copy(r.UnreachablePeers[i:], r.UnreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

//UpdatePeer update an existing peer on the repository
func (r *Repository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.Peers {
		if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

//GetPeerByIP get a peer from the repository using its ip
func (r *Repository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.Peers); i++ {
		if r.Peers[i].Identity().IP().Equal(ip) {
			return r.Peers[i], nil
		}
	}
	return
}

func (r *Repository) containsPeer(p discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.Peers {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}

func (r *Repository) containsUnreacheablePeer(pk discovery.PublicKey) bool {
	for _, up := range r.UnreachablePeers {
		if up.Equals(pk) {
			return true
		}
	}
	return false
}
