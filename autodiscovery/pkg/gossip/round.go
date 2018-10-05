package gossip

import (
	"encoding/hex"
	"errors"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("Unreachable Peer")

//Round decribes a gossip round
type Round struct {
	initator discovery.Peer
	target   discovery.Peer
	msg      Messenger
}

//Messenger is the interface that provides methods to send gossip requests
type Messenger interface {
	SendSyn(SynRequest) (*SynAck, error)
	SendAck(AckRequest) error
}

//SynRequest describes a SYN gossip request
type SynRequest struct {
	Initiator  discovery.Peer
	Target     discovery.Peer
	KnownPeers []discovery.Peer
}

//SynAck describes a SYN gossip response
type SynAck struct {
	Initiator    discovery.Peer
	Target       discovery.Peer
	UnknownPeers []discovery.Peer
	NewPeers     []discovery.Peer
}

//AckRequest describes an ACK gossip request
type AckRequest struct {
	Initiator      discovery.Peer
	Target         discovery.Peer
	RequestedPeers []discovery.Peer
}

//Spread starts messenging with a target peer
func (r *Round) Spread(kp []discovery.Peer, discovP chan<- discovery.Peer, reachP chan<- discovery.Peer, unreachP chan<- discovery.Peer) error {
	res, err := r.msg.SendSyn(SynRequest{r.initator, r.target, kp})
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

	if len(res.UnknownPeers) > 0 {
		reqPeers := make([]discovery.Peer, 0)
		mapPeers := r.mapPeers(kp)
		for _, p := range res.UnknownPeers {
			if k, exist := mapPeers[hex.EncodeToString(p.Identity().PublicKey())]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := r.msg.SendAck(AckRequest{r.initator, r.target, reqPeers}); err != nil {
			//We do not throw an error when the peer is unreachable
			//Gossip must continue
			//We catch the unreachable peer, store somewhere
			if err.Error() == ErrUnreachablePeer.Error() {
				unreachP <- r.target
				return nil
			}
			return err
		}

		for _, p := range res.NewPeers {
			discovP <- p
		}
	}

	return nil
}

func (r Round) mapPeers(pp []discovery.Peer) map[string]discovery.Peer {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}
	return mPeers
}

//NewGossipRound creates a new gossip round
func NewGossipRound(init discovery.Peer, target discovery.Peer, msg Messenger) *Round {
	return &Round{
		initator: init,
		target:   target,
		msg:      msg,
	}
}
