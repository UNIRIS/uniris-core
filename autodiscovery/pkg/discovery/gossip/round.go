package gossip

import (
	"encoding/hex"
	"math/rand"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
)

//GossipRound wraps gossip round mechanism
type GossipRound struct {
	Initiator  discovery.Peer
	KnownPeers []discovery.Peer
	Seeds      []discovery.SeedPeer
}

//NewGossipRound creates a gossip round
func NewGossipRound(initiator discovery.Peer, seeds []discovery.SeedPeer, knownPeers []discovery.Peer) GossipRound {
	return GossipRound{
		Initiator:  initiator,
		Seeds:      seeds,
		KnownPeers: knownPeers,
	}
}

//SelectPeers perform a gossip peer selection
func (r GossipRound) SelectPeers() []discovery.Peer {
	peers := make([]discovery.Peer, 0)
	randomSeed := r.randomSeed(r.Seeds)
	peers = append(peers, randomSeed.AsPeer())
	if len(r.excludeOwnPeer()) > 0 {
		randomPeer := r.randomPeer(r.KnownPeers)
		peers = append(peers, randomPeer)
	}
	return peers
}

//RequestedPeers returns the known peers from a list of requested peers
func (r GossipRound) RequestedPeers(pp []discovery.Peer) []discovery.Peer {
	knownPeersMap := make(map[string]discovery.Peer, 0)
	for _, p := range r.KnownPeers {
		knownPeersMap[hex.EncodeToString(p.PublicKey())] = p
	}

	peers := make([]discovery.Peer, 0)
	for _, p := range pp {
		if k, exist := knownPeersMap[hex.EncodeToString(p.PublicKey())]; exist && k.IsOk() {
			peers = append(peers, k)
		}
	}
	return peers
}

func (r GossipRound) excludeOwnPeer() []discovery.Peer {
	noExclude := make([]discovery.Peer, 0)
	for _, peer := range r.KnownPeers {
		if !peer.IsOwned() {
			noExclude = append(noExclude, peer)
		}
	}
	return noExclude
}

func (r GossipRound) randomPeer(peers []discovery.Peer) discovery.Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return peers[0]
}

func (r GossipRound) randomSeed(seeds []discovery.SeedPeer) discovery.SeedPeer {
	if len(seeds) > 1 {
		rnd := rand.Intn(len(seeds) - 1)
		return seeds[rnd]
	}
	return seeds[0]
}
