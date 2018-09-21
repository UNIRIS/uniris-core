package mem

import (
	"encoding/hex"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type Repository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *Repository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.IsOwned() {
			return p, nil
		}
	}
	return
}

func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

func (r *Repository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

func (r *Repository) SetPeer(peer discovery.Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.peers {
			if string(p.PublicKey()) == string(peer.PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.peers = append(r.peers, peer)
	}
	return nil
}

func (r *Repository) SetSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
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
