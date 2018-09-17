package mock

import (
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/gossip"
)

type MockGossipMessenger struct{}

func NewGossipMessenger() gossip.GossipMessenger {
	return MockGossipMessenger{}
}

func (m MockGossipMessenger) SendSyn(req gossip.SynRequest) (gossip.SynAck, error) {
	newPeers := make([]discovery.Peer, 0)
	unknownPeers := make([]discovery.Peer, 0)

	return gossip.SynAck{
		Initiator:    discovery.Peer{},
		Receiver:     discovery.Peer{},
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m MockGossipMessenger) SendAck(req gossip.AckRequest) error {
	return nil
}
