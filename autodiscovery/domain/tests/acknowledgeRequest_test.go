package tests

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Performs the acknowledge after a synchronize request
	Given peers store locally
	When we receive a synchronize request within a list of peers
	Then we return the unknown peers from the initator and the wished unknown peers
*/
func TestAcknowledgeRequest(t *testing.T) {
	repo := GetRepo()

	usecases.LoadSeedPeers(&SeedLoader{}, repo)
	usecases.StartupPeer(repo, &GeolocService{}, GetSecondValidPublicKey(), 3455)

	requestedPeers := []*entities.Peer{
		&entities.Peer{
			PublicKey: GetThirdValidPublicKey(),
			IP:        net.ParseIP("20.10.0.40"),
		},
	}

	ack, err := usecases.AcknowledgeRequest(repo, requestedPeers)
	assert.Nil(t, err)

	assert.NotEmpty(t, ack.UnknownInitiatorPeers)
	assert.Equal(t, 1, len(ack.UnknownInitiatorPeers))
	assert.Equal(t, "127.0.0.1", ack.UnknownInitiatorPeers[0].IP.String())

	assert.NotEmpty(t, ack.WishedUnknownPeers)
	assert.Equal(t, 1, len(ack.WishedUnknownPeers))
	assert.Equal(t, "20.10.0.40", ack.WishedUnknownPeers[0].IP.String())
}

/*
Scenario: Performs the acknowledge after a synchronize request
	Given peers store locally
	When we receive a synchronize request within an empty of peers
	Then we return the known peers and an empty wished unknown peers
*/
func TestAcknowledgeRequestWithEmptyRequest(t *testing.T) {
	repo := GetRepo()

	usecases.LoadSeedPeers(&SeedLoader{}, repo)

	requestedPeers := []*entities.Peer{}

	ack, err := usecases.AcknowledgeRequest(repo, requestedPeers)
	assert.Nil(t, err)

	assert.NotEmpty(t, ack.UnknownInitiatorPeers)
	assert.Equal(t, 1, len(ack.UnknownInitiatorPeers))
	assert.Equal(t, "127.0.0.1", ack.UnknownInitiatorPeers[0].IP.String())

	assert.Empty(t, ack.WishedUnknownPeers)
}
