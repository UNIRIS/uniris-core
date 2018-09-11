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
	publicKey := GetValidPublicKey()
	ts := time.Now().Unix()

	err := usecases.StartupPeer(repo, &GeolocService{}, GetValidPublicKey(), 3545)
	assert.Nil(t, err)

	peer, _ := repo.GetLocalPeer()

	assert.Equal(t, net.ParseIP("127.0.0.1"), peer.IP, "IP is not 127.0.0.1")
	assert.Equal(t, 3545, peer.Port, "Port is not 3535")
	assert.Equal(t, publicKey, peer.PublicKey, "Public key is not %s", publicKey)
	assert.Equal(t, entities.BootstrapingState, peer.AppState.State, "Peer must boostraping")
	assert.Equal(t, ts, peer.Heartbeat.GenerationTime.Unix(), "Generation time is not %d", ts)
	assert.Equal(t, entities.DiscoveredCategory, peer.Category, "Peer category must be discovered")

	time.Sleep(time.Second * 1)
	assert.Equal(t, peer.GetElapsedHeartbeats(), int64(1), "Peer heartbeat must be equal to 1 seconds")
}
