package tests

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Startup a new peer
        Given a peer
        When it startups
        Then the peer is initialized and store locally
*/
func TestStartupPeer(t *testing.T) {
	repo := GetRepo()
	ts := time.Now().Unix()

	err := usecases.StartupPeer(repo, &GeolocService{}, GetValidPublicKey(), 3545)
	assert.Nil(t, err)

	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
	assert.Equal(t, GetValidPublicKey(), peers[0].PublicKey, "Public key is not %s", string(GetValidPublicKey()))
	assert.Equal(t, net.ParseIP("127.0.0.1"), peers[0].IP, "IP is not 127.0.0.1")
	assert.Equal(t, entities.BootstrapingState, peers[0].AppState.State, "Peer must boostraping")
	assert.Equal(t, ts, peers[0].Heartbeat.GenerationTime.Unix(), "Generation time is not %d", ts)

	time.Sleep(time.Second * 1)
	assert.Equal(t, peers[0].GetElapsedHeartbeats(), int64(1), "Peer heartbeat must be equal to 1 seconds")
}
