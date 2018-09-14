package ports

import "github.com/uniris/uniris-core/autodiscovery/domain"

//GossipMessenger wraps the gossip messager
type GossipMessenger interface {
	SendSynchro(req domain.SynRequest) (*domain.SynAck, error)
	SendAcknowledgement(req domain.AckRequest) error
}
