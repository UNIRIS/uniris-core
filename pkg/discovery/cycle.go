package discovery

import (
	"log"
	"math/rand"
	"sync"
)

//Cycle represents a gossip cycle
type Cycle struct {
	Discoveries []Peer
	Reaches     []PeerIdentity
	Unreaches   []PeerIdentity
}

//run starts gossip cycle by creating rounds from a peer selection to spread the known peers and discover new peers
func (c *Cycle) run(initator Peer, msg RoundMessenger, seeds []PeerIdentity, peers []Peer, reaches []PeerIdentity, unreaches []PeerIdentity) error {

	selected, err := c.selectRandomPeers(initator, seeds, reaches, unreaches)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(selected))

	for _, p := range selected {
		go func(target PeerIdentity) {
			defer wg.Done()
			if err := c.startRound(target, peers, msg); err != nil {
				log.Printf("unexpected error during round execution: %s", err.Error())
				return
			}
		}(p)
	}

	wg.Wait()

	return nil
}

func (c *Cycle) startRound(target PeerIdentity, peers []Peer, msg RoundMessenger) error {
	r := round{
		target: target,
		peers:  peers,
	}
	peers, err := r.run(msg)
	if err != nil {
		if err == ErrUnreachablePeer {
			c.Unreaches = append(c.Unreaches, target)
			return nil
		}
		return err
	}
	c.Discoveries = append(c.Discoveries, peers...)
	c.Reaches = append(c.Reaches, target)
	return nil
}

func (c Cycle) selectRandomPeers(initator Peer, seeds []PeerIdentity, reachP []PeerIdentity, unreachP []PeerIdentity) ([]PeerIdentity, error) {

	peers := make([]PeerIdentity, 0)

	//We pick a random seed peer
	//and exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	s := c.randomPeer(seeds)
	if s.Endpoint() != initator.Identity().Endpoint() {
		peers = append(peers, s)
	}

	//We pick a random reachables(discovered) peer and we filter ourself (we don't want gossip with ourself)
	filteredReachP := make([]PeerIdentity, 0)
	for _, p := range reachP {
		if p.Endpoint() != initator.Identity().Endpoint() {
			filteredReachP = append(filteredReachP, p)
		}
	}
	if len(filteredReachP) > 0 {
		peers = append(peers, c.randomPeer(filteredReachP))
	}

	//We pick a random unreachable peer
	if len(unreachP) > 0 {
		peers = append(peers, c.randomPeer(unreachP))
	}

	return peers, nil
}

func (c Cycle) randomPeer(items []PeerIdentity) PeerIdentity {
	if len(items) > 1 {
		rnd := rand.Intn(len(items))
		return items[rnd]
	}
	return items[0]
}
