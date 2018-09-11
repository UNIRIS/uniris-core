package entities

type Acknowledge struct {
	UnknownInitiatorPeers []*Peer
	WishedUnknownPeers    []*Peer
}
