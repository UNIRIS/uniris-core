package gossip

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

//SynRequest describes a SYN gossip request
type SynRequest struct {

	//Initiator discovery.peer of the request
	Initiator discovery.Peer

	//Receiver discovery.peer of the request
	Receiver discovery.Peer

	//Known peers by the initiator
	KnownPeers []discovery.Peer
}

//SynAck describes a SYN gossip response
type SynAck struct {

	//Initiator discovery.peer of the response
	Initiator discovery.Peer

	//Receiver discovery.peer of the response
	Receiver discovery.Peer

	//Peers unknown from the SYN request initator
	NewPeers []discovery.Peer

	//Peers unknown from the SYN request receiver
	UnknownPeers []discovery.Peer
}

//AckRequest describes an ACK gossip request
type AckRequest struct {

	//Initiator discovery.peer of the ACK request
	Initiator discovery.Peer

	//Receiver discovery.peer of the ACK request
	Receiver discovery.Peer

	//Detailed peers requested by the SYN request receiver
	RequestedPeers []discovery.Peer
}

//NewSynRequest creates a new SYN gossip request
func NewSynRequest(initiator discovery.Peer, receiver discovery.Peer, knownPeers []discovery.Peer) SynRequest {
	return SynRequest{
		Initiator:  initiator,
		Receiver:   receiver,
		KnownPeers: knownPeers,
	}
}

//NewAckRequest creates a new ACK gossip request
func NewAckRequest(initiator discovery.Peer, receiver discovery.Peer, requestedPeers []discovery.Peer) AckRequest {
	return AckRequest{
		Initiator:      initiator,
		Receiver:       receiver,
		RequestedPeers: requestedPeers,
	}
}

//NewSynAck builds a new response from a SYN request (aka SYN-ACK)
func NewSynAck(initiator discovery.Peer, receiver discovery.Peer, newPeers []discovery.Peer, unknownPeers []discovery.Peer) SynAck {
	return SynAck{
		Initiator:    initiator,
		Receiver:     receiver,
		UnknownPeers: newPeers,
		NewPeers:     unknownPeers,
	}
}
