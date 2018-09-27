package mock

import (
	"encoding/hex"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type Repository struct {
	Peers []discovery.Peer
	Seeds []discovery.Seed
}

func (r *Repository) CountKnownPeers() (int, error) {
	return len(r.Peers), nil
}

func (r *Repository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.Peers {
		if p.Owned() {
			return p, nil
		}
	}
	return
}

func (r *Repository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.Peers = append(r.Peers, p)
	return nil
}

func (r *Repository) AddSeed(s discovery.Seed) error {
	r.Seeds = append(r.Seeds, s)
	return nil
}

func (r *Repository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.Peers, nil
}

func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.Seeds, nil
}

func (r *Repository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.Peers {
		if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *Repository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.Peers); i++ {
		if ip.Equal(r.Peers[i].Identity().IP()) {
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
