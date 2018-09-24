package discovery

import (
	"errors"
	"math/rand"
)

//GossipCycle describes a cycle in a gossip protocol
type GossipCycle struct {
	initator   Peer
	knownPeers []Peer
	seedPeers  []Seed
	rounds     []*GossipRound
}

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//NewGossipCycle creates a gossip cucle
//
//If an empty list of seeds is provided an error is returned
func NewGossipCycle(initator Peer, kp []Peer, sp []Seed) (*GossipCycle, error) {
	if sp == nil || len(sp) == 0 {
		return nil, ErrEmptySeed
	}

	return &GossipCycle{
		initator:   initator,
		knownPeers: kp,
		seedPeers:  sp,
	}, nil
}

//SelectPeers chooses random seed and peer from the repository excluding ourself
func (c GossipCycle) SelectPeers() []Peer {
	peers := make([]Peer, 0)

	//We pick a random seed peer
	s := c.randomSeed().ToPeer()

	//Exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	if s.Endpoint() != c.initator.Endpoint() {
		peers = append(peers, s)
	}

	//Filter ourself (we don't want gossip with ourself)
	filtered := make([]Peer, 0)
	for _, p := range c.knownPeers {
		if !p.IsOwned() {
			filtered = append(filtered, p)
		}
	}
	c.knownPeers = filtered

	//We pick a random known peer
	if len(c.knownPeers) > 0 {
		peers = append(peers, c.randomPeer())
	}
	return peers
}

//CreateRound creates a gossip round
func (c *GossipCycle) CreateRound(target Peer) *GossipRound {
	r := NewGossipRound(c.initator, target)
	c.rounds = append(c.rounds, r)
	return r
}

//Discoveries returns the discovered peers from the related started rounds
func (c *GossipCycle) Discoveries() []Peer {
	pp := make([]Peer, 0)
	for _, r := range c.rounds {
		for _, p := range r.discoveredPeers {
			pp = append(pp, p)
		}
	}
	return pp
}

func (c GossipCycle) randomPeer() Peer {
	if len(c.knownPeers) > 1 {
		rnd := rand.Intn(len(c.knownPeers) - 1)
		return c.knownPeers[rnd]
	}
	return c.knownPeers[0]
}

func (c GossipCycle) randomSeed() Seed {
	if len(c.seedPeers) > 1 {
		rnd := rand.Intn(len(c.seedPeers) - 1)
		return c.seedPeers[rnd]
	}
	return c.seedPeers[0]
}
