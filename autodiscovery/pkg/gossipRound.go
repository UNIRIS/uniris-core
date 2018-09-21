package discovery

import (
	"errors"
	"math/rand"
)

//ErrEmptySeed is returns when no seeds has been provided
var ErrEmptySeed = errors.New("Cannot start a gossip round without a list seeds")

//GossipRound describes a round in a gossip protocol
type GossipRound struct {
	initiator  Peer
	knownPeers []Peer
	seedPeers  []Seed
	msg        GossipCycleMessenger
}

//GossipRoundNotifier is the interface that provides methods to dispatch events during the round
type GossipRoundNotifier interface {

	//DisptachNewPeer notifies a new peers has been discovered
	DisptachNewPeer(peer Peer)
}

//Run starts a gossip round
func (r GossipRound) Run(init Peer, kp []Peer) ([]Peer, error) {
	newPeers := make([]Peer, 0)

	//Pick random known peer and seed
	pp, err := r.selectPeers()
	if err != nil {
		return nil, err
	}

	for _, p := range pp {
		pp, err := NewGossipCycle(init, p, kp, r.msg).Run()
		if err != nil {
			//We do not throw an error when the peer is unreachable
			//Gossip must continue
			if err == ErrPeerUnreachable {
				return nil, nil
			}
			return nil, err
		}
		for _, p := range pp {
			newPeers = append(newPeers, p)
		}
	}
	return newPeers, nil
}

//SelectPeers returns a seed and a known peer randomly
func (r GossipRound) selectPeers() ([]Peer, error) {
	peers := make([]Peer, 0)

	//We pick a random seed peer
	s := r.randomSeed().ToPeer()

	//Exclude ourself if we are of inside our list seed (impossible in reality, useful for testing)
	if s.GetEndpoint() != r.initiator.GetEndpoint() {
		peers = append(peers, s)
	}

	//Filter ourself (we don't want gossip with ourself)
	filtered := make([]Peer, 0)
	for _, p := range r.knownPeers {
		if p.GetEndpoint() != r.initiator.GetEndpoint() {
			filtered = append(filtered, p)
		}
	}
	r.knownPeers = filtered

	//We pick a random known peer
	if len(r.knownPeers) > 0 {
		peers = append(peers, r.randomPeer())
	}
	return peers, nil
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
func NewGossipRound(init Peer, kp []Peer, sp []Seed, msg GossipCycleMessenger) (*GossipRound, error) {

	if sp == nil || len(sp) == 0 {
		return nil, ErrEmptySeed
	}

	return &GossipRound{
		initiator:  init,
		knownPeers: kp,
		seedPeers:  sp,
		msg:        msg,
	}, nil
}
