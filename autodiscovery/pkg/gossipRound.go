package discovery

import (
	"errors"
	"math/rand"
)

//ErrEmptySeed is returnes when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//GossipRound describes a round in a gossip protocol
type GossipRound struct {
	initiator  Peer
	knownPeers []Peer
	seedPeers  []Seed
}

//SelectPeers returns a seed and a known peer randomly
func (r GossipRound) SelectPeers() []Peer {
	peers := make([]Peer, 0)
	peers = append(peers, r.randomSeed().ToPeer())

	//We exclude the owned peer, to not dial with ourself
	filtered := make([]Peer, 0)
	for _, peer := range r.knownPeers {
		if !peer.IsOwned() {
			filtered = append(filtered, peer)
		}
	}
	r.knownPeers = filtered

	if len(r.knownPeers) > 0 {
		peers = append(peers, r.randomPeer())
	}
	return peers
}

func (r GossipRound) randomPeer() Peer {
	if len(r.knownPeers) > 1 {
		rnd := rand.Intn(len(r.knownPeers) - 1)
		return r.knownPeers[rnd]
	}
	return r.knownPeers[0]
}

func (r GossipRound) randomSeed() Seed {
	if len(r.seedPeers) > 1 {
		rnd := rand.Intn(len(r.seedPeers) - 1)
		return r.seedPeers[rnd]
	}
	return r.seedPeers[0]
}

//NewGossipRound creates a gossip round
//
//If an empty list of seeds is provided an error is returned
func NewGossipRound(init Peer, kp []Peer, sp []Seed) (*GossipRound, error) {

	if sp == nil || len(sp) == 0 {
		return nil, ErrEmptySeed
	}

	return &GossipRound{
		initiator:  init,
		knownPeers: kp,
		seedPeers:  sp,
	}, nil
}
