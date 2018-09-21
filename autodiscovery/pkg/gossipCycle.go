package discovery

import (
	"encoding/hex"
	"errors"
)

var ErrPeerUnreachable = errors.New("Cannot reach the peer %s")

//GossipCycle defines a gossip cycle
type GossipCycle struct {
	initator   Peer
	receiver   Peer
	knownPeers []Peer
	msg        GossipCycleMessenger
}

//GossipCycleMessenger is the interface that provides methods to gossip during a cycle
type GossipCycleMessenger interface {

	//Sends a SYN request
	SendSyn(SynRequest) (*SynAck, error)

	//Sends a ACK request after receipt of the SYN request
	SendAck(AckRequest) error
}

//SynRequest describes a SYN gossip request
type SynRequest struct {

	//Initiator discovery.peer of the request
	Initiator Peer

	//Receiver discovery.peer of the request
	Receiver Peer

	//Known peers by the initiator
	KnownPeers []Peer
}

//SynAck describes a SYN gossip response
type SynAck struct {

	//Initiator discovery.peer of the response
	Initiator Peer

	//Receiver discovery.peer of the response
	Receiver Peer

	//Peers unknown from the SYN request initator
	NewPeers []Peer

	//Peers unknown from the SYN request receiver
	UnknownPeers []Peer
}

//AckRequest describes an ACK gossip request
type AckRequest struct {

	//Initiator discovery.peer of the ACK request
	Initiator Peer

	//Receiver discovery.peer of the ACK request
	Receiver Peer

	//Detailed peers requested by the SYN request receiver
	RequestedPeers []Peer
}

//Run the gossip with the selected peeer
func (c GossipCycle) Run() ([]Peer, error) {

	//Send SYN request
	synAck, err := c.msg.SendSyn(SynRequest{c.initator, c.receiver, c.knownPeers})
	if err != nil {
		return nil, err
	}

	//Provide details about the requested peers from the SYN-ACK initator (aka SYN receiver)
	if len(synAck.UnknownPeers) > 0 {
		reqPeers := make([]Peer, 0)
		mapPeers := c.mapPeers(c.knownPeers)
		for _, p := range synAck.UnknownPeers {
			if k, exist := mapPeers[hex.EncodeToString(p.PublicKey())]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := c.msg.SendAck(AckRequest{c.initator, c.receiver, reqPeers}); err != nil {
			return nil, err
		}
	}

	//We get the new discovered peers to be stored and notified
	return synAck.NewPeers, nil
}

func (c GossipCycle) mapPeers(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}
	return mPeers
}

//NewGossipCycle creates a new gossip cycle with its dependencies
func NewGossipCycle(init Peer, rec Peer, kp []Peer, msg GossipCycleMessenger) GossipCycle {
	return GossipCycle{init, rec, kp, msg}
}
