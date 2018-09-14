package ports

import "github.com/uniris/uniris-core/autodiscovery/core/domain"

//GossipBroker wraps the gossip message communication
type GossipBroker interface {
	SendSyn(req domain.SynRequest) (domain.SynAck, error)
	SendAck(req domain.AckRequest) error
}
