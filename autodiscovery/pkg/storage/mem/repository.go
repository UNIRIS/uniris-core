package mem

import (
	"encoding/hex"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Repository provides access to the peer repository
type Repository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

//GetOwnedPeer return the local peer
func (r *Repository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.IsOwned() {
			return p, nil
		}
	}
	return
}

//ListSeedPeers return all the seed on the repository
func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

//ListKnownPeers returns all the peers on the repository
func (r *Repository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

//AddPeer add a peer to the repository
func (r *Repository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

//AddSeed add a seed to the repository
func (r *Repository) AddSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

//UpdatePeer update an existing peer on the repository
func (r *Repository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.peers {
		if string(p.PublicKey()) == string(peer.PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *Repository) containsPeer(p discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}

	_, exist := mPeers[hex.EncodeToString(p.PublicKey())]
	return exist
}
