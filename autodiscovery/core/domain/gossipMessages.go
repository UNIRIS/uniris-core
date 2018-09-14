package domain

//SynRequest wraps a SYN gossip request
type SynRequest struct {
	Initiator  Peer
	Receiver   Peer
	KnownPeers []Peer
}

//SynAck wraps a SYN gossip response
type SynAck struct {
	Initiator    Peer
	Receiver     Peer
	NewPeers     []Peer
	UnknownPeers []Peer
}

//AckRequest wraps an ACK gossip request
type AckRequest struct {
	Initiator      Peer
	Receiver       Peer
	RequestedPeers []Peer
}

//NewSynRequest builds a SYN gossip request
func NewSynRequest(initiator Peer, receiver Peer, knownPeers []Peer) SynRequest {
	return SynRequest{
		Initiator:  initiator,
		Receiver:   receiver,
		KnownPeers: knownPeers,
	}
}

//NewAckRequest builds a ACK gossip request
func NewAckRequest(initiator Peer, receiver Peer, requestedPeers []Peer) AckRequest {
	return AckRequest{
		Initiator:      initiator,
		Receiver:       receiver,
		RequestedPeers: requestedPeers,
	}
}

//NewSynAck builds a SYN-ACK response
func NewSynAck(initiator Peer, receiver Peer, newPeers []Peer, unknownPeers []Peer) SynAck {
	return SynAck{
		Initiator: initiator,
		Receiver:  receiver,
	}
}
