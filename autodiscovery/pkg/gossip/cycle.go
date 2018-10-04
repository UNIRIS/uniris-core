package gossip

import (
	"errors"
	"math/rand"
	"sync"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Cycle describes a cycle in a gossip protocol
type Cycle struct {
	initator discovery.Peer
	msg      Messenger
	result   gossipChannel
}

type gossipChannel struct {
	discoveries  chan discovery.Peer
	unreachables chan discovery.Peer
	reaches      chan discovery.Peer
	errors       chan error
}

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//NewGossipCycle creates a gossip cucle
//
//If an empty list of seeds is provided an error is returned
func NewGossipCycle(initiator discovery.Peer, msg Messenger) *Cycle {
	return &Cycle{
		msg:      msg,
		initator: initiator,
		result: gossipChannel{
			discoveries:  make(chan discovery.Peer),
			unreachables: make(chan discovery.Peer),
			errors:       make(chan error, 1),
			reaches:      make(chan discovery.Peer),
		},
	}
}

//Run starts gossip by creating rounds from a peer selection to spread the known peers and discover new peers
func (c Cycle) Run(init discovery.Peer, selectedP []discovery.Peer, knownPeers []discovery.Peer) {

	var wg sync.WaitGroup
	wg.Add(len(selectedP))

	for _, p := range selectedP {
		go func(target discovery.Peer) {
			defer wg.Done()

			r := NewGossipRound(init, target, c.msg)

			if err := r.Spread(knownPeers, c.result.discoveries, c.result.reaches, c.result.unreachables); err != nil {
				c.result.errors <- err
			}
		}(p)
	}

	wg.Wait()
}

//SelectPeers chooses random seed and peer from the repository excluding ourself
func (c Cycle) SelectPeers(seeds []discovery.Seed, reachP []discovery.Peer, unreachP []discovery.Peer) ([]discovery.Peer, error) {
	if seeds == nil || len(seeds) == 0 {
		return nil, ErrEmptySeed
	}

	peers := make([]discovery.Peer, 0)

	//We pick a random seed peer
	//and exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	ppSeeds := make([]discovery.Peer, 0)
	for _, s := range seeds {
		ppSeeds = append(ppSeeds, s.AsPeer())
	}
	s := c.random(ppSeeds)
	if s.Endpoint() != c.initator.Endpoint() {
		peers = append(peers, s)
	}

	//We pick a random reachables(discovered) peer and we filter ourself (we don't want gossip with ourself)
	filteredReachP := make([]discovery.Peer, 0)
	for _, p := range reachP {
		if p.Endpoint() != c.initator.Endpoint() {
			filteredReachP = append(filteredReachP, p)
		}
	}
	if len(filteredReachP) > 0 {
		peers = append(peers, c.random(filteredReachP))
	}

	//We pick a random unreachable peer
	if len(unreachP) > 0 {
		peers = append(peers, c.random(unreachP))
	}

	return peers, nil
}

func (c Cycle) random(items []discovery.Peer) discovery.Peer {
	if len(items) > 1 {
		rnd := rand.Intn(len(items))
		return items[rnd]
	}
	return items[0]
}
