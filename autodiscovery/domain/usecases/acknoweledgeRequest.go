package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//AcknowledgeRequest handles a synchronize request and returns the unknown peers diff and the wished peers
func AcknowledgeRequest(peerRepo repositories.PeerRepository, requestedPeers []*entities.Peer) (*entities.Acknowledge, error) {
	knownPeers, err := peerRepo.ListPeers()
	if err != nil {
		return nil, err
	}
	mySelf, err := peerRepo.GetLocalPeer()
	if err != nil {
		return nil, err
	}

	if len(requestedPeers) == 0 {
		return &entities.Acknowledge{
			UnknownInitiatorPeers: knownPeers,
			WishedUnknownPeers:    []*entities.Peer{},
		}, nil
	}

	return &entities.Acknowledge{
		UnknownInitiatorPeers: GetUnknownPeers(knownPeers, requestedPeers, mySelf),
		WishedUnknownPeers:    GetUnknownPeers(requestedPeers, knownPeers, mySelf),
	}, nil
}
