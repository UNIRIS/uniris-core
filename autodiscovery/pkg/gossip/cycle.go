package gossip

import (
	"errors"
	"log"
	"math/rand"
	"sync"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Cycle describes a cycle in a gossip protocol
type Cycle struct {
	initator   discovery.Peer
	knownPeers []discovery.Peer
	seedPeers  []discovery.Seed
	msg        Messenger
	result     gossipChannel
}

type gossipChannel struct {
	discoveries  chan discovery.Peer
	unreachables chan discovery.Peer
	errors       chan error
}

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//NewGossipCycle creates a gossip cucle
//
//If an empty list of seeds is provided an error is returned
func NewGossipCycle(initator discovery.Peer, kp []discovery.Peer, sp []discovery.Seed, msg Messenger) (*Cycle, error) {
	if sp == nil || len(sp) == 0 {
		return nil, ErrEmptySeed
	}

	return &Cycle{
		initator:   initator,
		knownPeers: kp,
		seedPeers:  sp,
		msg:        msg,
		result: gossipChannel{
			discoveries:  make(chan discovery.Peer),
			unreachables: make(chan discovery.Peer),
			errors:       make(chan error, 1),
		},
	}, nil
}

//Run starts gossip by creating rounds from a peer selection to spread the known peers and discover new peers
func (c Cycle) Run() {
	defer close(c.result.discoveries)
	defer close(c.result.unreachables)
	defer close(c.result.errors)

	pp := c.selectPeers()

	var wg sync.WaitGroup
	wg.Add(len(pp))

	for _, p := range pp {
		log.Printf("Gossip with %s", p.Endpoint())
		go func(target discovery.Peer) {
			defer wg.Done()

			r := NewGossipRound(c.initator, target, c.msg)
			if err := r.Spread(c.knownPeers, c.result.discoveries, c.result.unreachables); err != nil {
				c.result.errors <- err
			}
		}(p)
	}

	wg.Wait()
}

//selectPeers chooses random seed and peer from the repository excluding ourself
func (c Cycle) selectPeers() []discovery.Peer {
	peers := make([]discovery.Peer, 0)

	//We pick a random seed peer
	s := c.randomSeed().AsPeer()

	//Exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	if s.Endpoint() != c.initator.Endpoint() {
		peers = append(peers, s)
	}

	filterP := make([]discovery.Peer, 0)
	for _, p := range c.knownPeers {
		if p.Endpoint() != c.initator.Endpoint() {
			filterP = append(filterP, p)
		}
	}

	//We pick a random known peer
	if len(filterP) > 0 {
		peers = append(peers, c.randomPeer(filterP))
	}

	return peers
}

func (c Cycle) randomPeer(peers []discovery.Peer) discovery.Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return c.knownPeers[0]
}

func (c Cycle) randomSeed() discovery.Seed {
	if len(c.seedPeers) > 1 {
		rnd := rand.Intn(len(c.seedPeers) - 1)
		return c.seedPeers[rnd]
	}
	return c.seedPeers[0]
}
