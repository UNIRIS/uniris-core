package entities

//SynchronizationRequest represents the gossip request to synchronize peers
type SynchronizationRequest struct {
	//Represents the peer which will receive the SYN request
	PeerReceiver *Peer

	//Represents the peers that the sender knowns
	KnownSenderPeers []*Peer
}

//AcknowledgeResponse represents the gossip acknowledge from the synchronization request
type AcknowledgeResponse struct {
	//Represents the unknown peers from the SYN initiator
	UnknownSenderPeers []*Peer
	//Represens the unknown peers from the SYN receiver
	UnknownReceiverPeers []*Peer
}
