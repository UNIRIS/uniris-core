package services

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//GossipService represents the autodiscovery requests
type GossipService interface {
	Synchronize(request *entities.SynchronizationRequest) (*entities.AcknowledgeResponse, error)
}
