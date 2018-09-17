package gossip

import "github.com/uniris/uniris-core/autodiscovery/pkg/discovery"

//SynRequest wraps a SYN gossip request
type SynRequest struct {
	Initiator  discovery.Peer
	Receiver   discovery.Peer
	KnownPeers []discovery.Peer
}

//SynAck wraps a SYN gossip response
type SynAck struct {
	Initiator    discovery.Peer
	Receiver     discovery.Peer
	NewPeers     []discovery.Peer
	UnknownPeers []discovery.Peer
}

//AckRequest wraps an ACK gossip request
type AckRequest struct {
	Initiator      discovery.Peer
	Receiver       discovery.Peer
	RequestedPeers []discovery.Peer
}

//NewSynRequest builds a SYN gossip request
func NewSynRequest(initiator discovery.Peer, receiver discovery.Peer, knownPeers []discovery.Peer) SynRequest {
	return SynRequest{
		Initiator:  initiator,
		Receiver:   receiver,
		KnownPeers: knownPeers,
	}
}

//NewAckRequest builds a ACK gossip request
func NewAckRequest(initiator discovery.Peer, receiver discovery.Peer, requestedPeers []discovery.Peer) AckRequest {
	return AckRequest{
		Initiator:      initiator,
		Receiver:       receiver,
		RequestedPeers: requestedPeers,
	}
}

//NewSynAck builds a SYN-ACK response
func NewSynAck(initiator discovery.Peer, receiver discovery.Peer, newPeers []discovery.Peer, unknownPeers []discovery.Peer) SynAck {
	return SynAck{
		Initiator: initiator,
		Receiver:  receiver,
	}
}
