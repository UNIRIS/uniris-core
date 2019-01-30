package discovery

import (
	"errors"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("Unreachable Peer")

//Client is the interface that provides methods to send gossip requests
type Client interface {
	SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error)
	SendAck(target PeerIdentity, requested []Peer) error
}

type round struct {
	target PeerIdentity
	cli    Client
}

//run starts the gossip round by messenging with the target peer
func (r round) run(kp []Peer, discovP chan<- Peer, reachP chan<- PeerIdentity, unreachP chan<- PeerIdentity) error {
	unknowns, news, err := r.cli.SendSyn(r.target, kp)
	if err != nil {
		//We do not throw an error when the peer is unreachable
		//Gossip must continue
		if err.Error() == ErrUnreachablePeer.Error() {
			unreachP <- r.target
			return nil
		}

		return err
	}

	//Notifies the peer's response
	reachP <- r.target

	if len(unknowns) > 0 {
		reqPeers := make([]Peer, 0)
		mapPeers := r.mapPeers(kp)
		for _, p := range unknowns {
			if k, exist := mapPeers[p.Identity().PublicKey()]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := r.cli.SendAck(r.target, reqPeers); err != nil {
			//We do not throw an error when the peer is unreachable
			//Gossip must continue
			//We catch the unreachable peer, store somewhere
			if err.Error() == ErrUnreachablePeer.Error() {
				unreachP <- r.target
				return nil
			}
			return err
		}

		for _, p := range news {
			discovP <- p
		}
	}

	return nil
}

func (r round) mapPeers(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer, 0)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
