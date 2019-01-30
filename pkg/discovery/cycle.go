package discovery

import (
	"errors"
	"math/rand"
	"sync"
)

type cycle struct {
	initator      Peer
	cli           Client
	prevReaches   []PeerIdentity
	prevUnreaches []PeerIdentity
	discoveryChan chan Peer
	unreachChan   chan PeerIdentity
	reachChan     chan PeerIdentity
	errChan       chan error
}

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

func newCycle(initiator Peer, cli Client, rP []PeerIdentity, unrP []PeerIdentity) cycle {
	return cycle{
		cli:           cli,
		initator:      initiator,
		prevReaches:   rP,
		prevUnreaches: unrP,
		discoveryChan: make(chan Peer),
		unreachChan:   make(chan PeerIdentity),
		errChan:       make(chan error, 1),
		reachChan:     make(chan PeerIdentity),
	}
}

//run starts gossip cycle by creating rounds from a peer selection to spread the known peers and discover new peers
func (c cycle) run(p Peer, ss []PeerIdentity, kp []Peer) {

	selected, err := c.selectRandomPeers(ss, c.prevReaches, c.prevUnreaches)
	if err != nil {
		c.errChan <- err
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(selected))

	for _, p := range selected {
		go func(target PeerIdentity) {
			defer wg.Done()

			r := round{target, c.cli}
			if err := r.run(kp, c.discoveryChan, c.reachChan, c.unreachChan); err != nil {
				c.errChan <- err
			}
		}(p)
	}

	wg.Wait()
}

func (c cycle) selectRandomPeers(seeds []PeerIdentity, reachP []PeerIdentity, unreachP []PeerIdentity) ([]PeerIdentity, error) {
	if seeds == nil || len(seeds) == 0 {
		return nil, ErrEmptySeed
	}

	peers := make([]PeerIdentity, 0)

	//We pick a random seed peer
	//and exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	s := c.random(seeds)
	if s.Endpoint() != c.initator.Identity().Endpoint() {
		peers = append(peers, s)
	}

	//We pick a random reachables(discovered) peer and we filter ourself (we don't want gossip with ourself)
	filteredReachP := make([]PeerIdentity, 0)
	for _, p := range reachP {
		if p.Endpoint() != c.initator.Identity().Endpoint() {
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

func (c cycle) random(items []PeerIdentity) PeerIdentity {
	if len(items) > 1 {
		rnd := rand.Intn(len(items))
		return items[rnd]
	}
	return items[0]
}
