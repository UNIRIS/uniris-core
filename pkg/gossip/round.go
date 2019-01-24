package gossip

import (
	"errors"

	uniris "github.com/uniris/uniris-core/pkg"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("Unreachable Peer")

//RoundMessenger is the interface that provides methods to send gossip requests
type RoundMessenger interface {
	SendSyn(source uniris.Peer, target uniris.Peer, known []uniris.Peer) (unknown []uniris.Peer, new []uniris.Peer, err error)
	SendAck(source uniris.Peer, target uniris.Peer, requested []uniris.Peer) error
}

type round struct {
	initator uniris.Peer
	target   uniris.Peer
	msg      RoundMessenger
}

//run starts the gossip round by messenging with the target peer
func (r round) run(kp []uniris.Peer, discovP chan<- uniris.Peer, reachP chan<- uniris.Peer, unreachP chan<- uniris.Peer) error {
	unknowns, news, err := r.msg.SendSyn(r.initator, r.target, kp)
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
		reqPeers := make([]uniris.Peer, 0)
		mapPeers := r.mapPeers(kp)
		for _, p := range unknowns {
			if k, exist := mapPeers[p.Identity().PublicKey()]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := r.msg.SendAck(r.initator, r.target, reqPeers); err != nil {
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

func (r round) mapPeers(pp []uniris.Peer) map[string]uniris.Peer {
	mPeers := make(map[string]uniris.Peer, 0)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}
