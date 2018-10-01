package discovery

import (
	"errors"
	"math/rand"
)

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//ErrNoOwnedPeer is returnes when no owned peers has been stored
var ErrNoOwnedPeer = errors.New("Cannot start a gossip round without a startuping peer")

//GossipRound describes a round in a gossip protocol
type GossipRound struct {
	initiator         Peer
	reacheablePeers   []Peer
	seedPeers         []Seed
	unreacheablePeers []Peer
}

//SelectPeers returns a seed , a known peer and an unreacheable peer randomly
func (r GossipRound) SelectPeers() ([]Peer, error) {
	peers := make([]Peer, 0)

	//We pick a random seed peer
	s := r.randomSeed().AsPeer()

	//Exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	owned := r.getOwnedPeer()
	if owned == nil {
		return nil, ErrNoOwnedPeer
	}
	if s.Endpoint() != owned.Endpoint() {
		peers = append(peers, s)
	}

	//Filter ourself (we don't want gossip with ourself)
	filtered := make([]Peer, 0)
	for _, p := range r.reacheablePeers {
		if !p.Owned() {
			filtered = append(filtered, p)
		}
	}
	r.reacheablePeers = filtered

	//We pick a random known peer
	if len(r.reacheablePeers) > 0 {
		peers = append(peers, r.randomPeer())
	}

	//We pick a random unreachable peer
	if len(r.unreacheablePeers) > 0 {
		peers = append(peers, r.randomUnreacheablePeer())

	}

	return peers, nil
}

func (r GossipRound) getOwnedPeer() Peer {
	for _, p := range r.reacheablePeers {
		if p.Owned() {
			return p
		}
	}
	return nil
}

func (r GossipRound) randomPeer() Peer {

	if len(r.reacheablePeers) > 1 {
		rnd := rand.Intn(len(r.reacheablePeers) - 1)
		return r.reacheablePeers[rnd]
	}
	return r.reacheablePeers[0]
}

func (r GossipRound) randomSeed() Seed {
	if len(r.seedPeers) > 1 {
		rnd := rand.Intn(len(r.seedPeers) - 1)
		return r.seedPeers[rnd]
	}
	return r.seedPeers[0]
}

func (r GossipRound) randomUnreacheablePeer() Peer {
	if len(r.unreacheablePeers) == 1 {
		return r.unreacheablePeers[0]
	}
	if len(r.unreacheablePeers) > 1 {
		rnd := rand.Intn(len(r.unreacheablePeers) - 1)
		return r.unreacheablePeers[rnd]
	}

	return nil
}

//NewGossipRound creates a gossip round
//
//If an empty list of seeds is provided an error is returned
func NewGossipRound(init Peer, rp []Peer, sp []Seed, up []Peer) (*GossipRound, error) {

	if sp == nil || len(sp) == 0 {
		return nil, ErrEmptySeed
	}

	return &GossipRound{
		initiator:         init,
		reacheablePeers:   rp,
		seedPeers:         sp,
		unreacheablePeers: up,
	}, nil
}
