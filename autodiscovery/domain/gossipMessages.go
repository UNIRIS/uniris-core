package domain

import "net"

//SynRequest represents the gossip request to synchronize peers
type SynRequest struct {
	Sender     Peer
	TargetIP   net.IP
	TargetPort int
	KnownPeers []Peer
}

//SynAck represents the gossip acknowledge from the synchronization request
type SynAck struct {
	NewPeers             []Peer
	DetailPeersRequested []Peer
}

//AckRequest represents the gossip acknownledge request for a gossip cycle
type AckRequest struct {
	TargetIP           net.IP
	TargetPort         int
	DetailedKnownPeers []Peer
}

//NewSynRequest builds a synchronization request
func NewSynRequest(sender Peer, targetIP net.IP, targetPort int, knownPeers []Peer) SynRequest {
	return SynRequest{
		Sender:     sender,
		TargetIP:   targetIP,
		TargetPort: targetPort,
		KnownPeers: knownPeers,
	}
}

//NewAckRequest builds a acknowledge request
func NewAckRequest(detailedKnownPeers []Peer) AckRequest {
	return AckRequest{
		DetailedKnownPeers: detailedKnownPeers,
	}
}
