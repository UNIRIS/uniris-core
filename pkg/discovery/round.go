package discovery

import (
	"errors"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("Unreachable Peer")

//RoundMessenger is the interface that provides methods to send gossip requests
type RoundMessenger interface {
	SendSyn(target PeerIdentity, known []Peer) (localDiscoveries []Peer, remoteDiscoveries []Peer, err error)
	SendAck(target PeerIdentity, requested []Peer) error
}

type round struct {
	target PeerIdentity
	peers  []Peer
}

//run starts the gossip round by messenging with the target peer
func (r round) run(msg RoundMessenger) ([]Peer, error) {
	localDiscoveries, remoteDiscoveries, err := msg.SendSyn(r.target, r.peers)
	if err != nil {
		return nil, err
	}

	if len(localDiscoveries) > 0 {
		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := msg.SendAck(r.target, localDiscoveries); err != nil {
			return nil, err
		}
	}
	return remoteDiscoveries, nil
}

func (r round) mapPeers() map[string]Peer {
	mPeers := make(map[string]Peer, 0)
	for _, p := range r.peers {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
