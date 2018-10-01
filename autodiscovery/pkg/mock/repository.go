package mock

import (
	"encoding/hex"
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type Repository struct {
	ownedPeer       discovery.Peer
	discoveredPeers []discovery.Peer
	seedPeers       []discovery.Seed
}

func (r *Repository) CountDiscoveredPeers() (int, error) {
	return len(r.discoveredPeers), nil
}

//GetOwnedPeer return the local peer
func (r *Repository) GetOwnedPeer() (discovery.Peer, error) {
	return r.ownedPeer, nil
}

//ListSeedPeers return all the seed on the mockRepository
func (r *Repository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seedPeers, nil
}

//ListDiscoveredPeers returns all the discoveredPeers on the mockRepository
func (r *Repository) ListDiscoveredPeers() ([]discovery.Peer, error) {
	return r.discoveredPeers, nil
}

func (r *Repository) SetPeer(peer discovery.Peer) error {
	if peer.Owned() {
		r.ownedPeer = peer
		return nil
	}
	if r.containsPeer(peer) {
		for _, p := range r.discoveredPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.discoveredPeers = append(r.discoveredPeers, peer)
	}
	return nil
}

func (r *Repository) SetSeed(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//GetPeerByIP get a peer from the mockRepository using its ip
func (r *Repository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	if r.ownedPeer.Identity().IP().Equal(ip) {
		return r.ownedPeer, nil
	}
	for i := 0; i < len(r.discoveredPeers); i++ {
		if r.discoveredPeers[i].Identity().IP().Equal(ip) {
			return r.discoveredPeers[i], nil
		}
	}
	return
}

func (r *Repository) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.discoveredPeers {
		mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}
