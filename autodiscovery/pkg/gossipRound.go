package discovery

import "encoding/hex"

//GossipRound decribes a gossip round
type GossipRound struct {
	initator        Peer
	target          Peer
	discoveredPeers []Peer
}

//GossipSpreader is the interface that provides methods to spread a gossip
type GossipSpreader interface {
	SendSyn(SynRequest) (*SynAck, error)
	SendAck(AckRequest) error
}

//SynRequest describes a SYN gossip request
type SynRequest struct {
	Initiator  Peer
	Target     Peer
	KnownPeers []Peer
}

//SynAck describes a SYN gossip response
type SynAck struct {
	Initiator    Peer
	Target       Peer
	NewPeers     []Peer
	UnknownPeers []Peer
}

//AckRequest describes an ACK gossip request
type AckRequest struct {
	Initiator      Peer
	Target         Peer
	RequestedPeers []Peer
}

//NewGossipRound creates a gossip round
func NewGossipRound(initator Peer, target Peer) *GossipRound {
	return &GossipRound{
		initator: initator,
		target:   target,
	}
}

//Spread starts messenging with a target peer and communicates known peers
func (r *GossipRound) Spread(kp []Peer, spr GossipSpreader) error {
	res, err := spr.SendSyn(SynRequest{r.initator, r.target, kp})
	if err != nil {
		return err
	}
	r.discoveredPeers = res.NewPeers
	if len(res.UnknownPeers) > 0 {
		reqPeers := make([]Peer, 0)
		mapPeers := r.mapPeers(kp)
		for _, p := range res.UnknownPeers {
			if k, exist := mapPeers[hex.EncodeToString(p.PublicKey())]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := spr.SendAck(AckRequest{r.initator, r.target, reqPeers}); err != nil {
			return err
		}
	}
	return nil
}

func (r GossipRound) mapPeers(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}
	return mPeers
}
