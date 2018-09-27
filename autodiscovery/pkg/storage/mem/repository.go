package mem

import (
	"encoding/hex"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type repo struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

//NewRepository implements the repository in memory
func NewRepository() discovery.Repository {
	return &repo{}
}

func (r *repo) CountKnownPeers() (int, error) {
	return len(r.peers), nil
}

//GetOwnedPeer return the local peer
func (r *repo) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.Owned() {
			return p, nil
		}
	}
	return
}

//ListSeedPeers return all the seed on the repository
func (r *repo) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

//ListKnownPeers returns all the peers on the repository
func (r *repo) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

//AddPeer add a peer to the repository
func (r *repo) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

//AddSeed add a seed to the repository
func (r *repo) AddSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

//UpdatePeer update an existing peer on the repository
func (r *repo) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.peers {
		if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

//GetPeerByIP get a peer from the repository using its ip
func (r *repo) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.peers); i++ {
		if r.peers[i].Identity().IP().Equal(ip) {
			return r.peers[i], nil
		}
	}
	return
}

func (r *repo) containsPeer(p discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}
