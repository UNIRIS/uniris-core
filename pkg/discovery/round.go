package discovery

import (
	"errors"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("Unreachable Peer")

//RoundMessenger is the interface that provides methods to send gossip requests
type RoundMessenger interface {
	SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error)
	SendAck(target PeerIdentity, requested []Peer) error
}

type round struct {
	target PeerIdentity
	peers  []Peer
}

//run starts the gossip round by messenging with the target peer
func (r round) run(msg RoundMessenger) ([]Peer, error) {
	unknowns, news, err := msg.SendSyn(r.target, r.peers)
	if err != nil {
		return nil, err
	}

	if len(unknowns) > 0 {
		reqPeers := make([]Peer, 0)
		mapPeers := r.mapPeers()
		for _, p := range unknowns {
			if k, exist := mapPeers[p.Identity().PublicKey()]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := msg.SendAck(r.target, reqPeers); err != nil {
			return nil, err
		}

		return news, nil
	}

	return nil, nil
}

func (r round) mapPeers() map[string]Peer {
	mPeers := make(map[string]Peer, 0)
	for _, p := range r.peers {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
