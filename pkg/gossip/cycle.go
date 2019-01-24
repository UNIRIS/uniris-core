package gossip

import (
	"errors"
	"math/rand"
	"sync"

	uniris "github.com/uniris/uniris-core/pkg"
)

type cycle struct {
	initator      uniris.Peer
	msg           RoundMessenger
	prevReaches   []uniris.Peer
	prevUnreaches []uniris.Peer
	discoveryChan chan uniris.Peer
	unreachChan   chan uniris.Peer
	reachChan     chan uniris.Peer
	errChan       chan error
}

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

func newCycle(initiator uniris.Peer, msg RoundMessenger, rP []uniris.Peer, unrP []uniris.Peer) cycle {
	return cycle{
		msg:           msg,
		initator:      initiator,
		prevReaches:   rP,
		prevUnreaches: unrP,
		discoveryChan: make(chan uniris.Peer),
		unreachChan:   make(chan uniris.Peer),
		errChan:       make(chan error, 1),
		reachChan:     make(chan uniris.Peer),
	}
}

//run starts gossip cycle by creating rounds from a peer selection to spread the known peers and discover new peers
func (c cycle) run(p uniris.Peer, ss []uniris.Seed, kp []uniris.Peer) {

	selected, err := c.selectRandomPeers(ss, c.prevReaches, c.prevUnreaches)
	if err != nil {
		c.errChan <- err
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(selected))

	for _, p := range selected {
		go func(target uniris.Peer) {
			defer wg.Done()

			r := round{c.initator, target, c.msg}
			if err := r.run(kp, c.discoveryChan, c.reachChan, c.unreachChan); err != nil {
				c.errChan <- err
			}
		}(p)
	}

	wg.Wait()
}

func (c cycle) selectRandomPeers(seeds []uniris.Seed, reachP []uniris.Peer, unreachP []uniris.Peer) ([]uniris.Peer, error) {
	if seeds == nil || len(seeds) == 0 {
		return nil, ErrEmptySeed
	}

	peers := make([]uniris.Peer, 0)

	//We pick a random seed peer
	//and exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	ppSeeds := make([]uniris.Peer, 0)
	for _, seed := range seeds {
		ppSeeds = append(ppSeeds, seed.AsPeer())
	}
	s := c.random(ppSeeds)
	if s.Endpoint() != c.initator.Endpoint() {
		peers = append(peers, s)
	}

	//We pick a random reachables(discovered) peer and we filter ourself (we don't want gossip with ourself)
	filteredReachP := make([]uniris.Peer, 0)
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

func (c cycle) random(items []uniris.Peer) uniris.Peer {
	if len(items) > 1 {
		rnd := rand.Intn(len(items))
		return items[rnd]
	}
	return items[0]
}
