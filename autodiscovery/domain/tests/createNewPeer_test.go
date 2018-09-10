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
Scenario: Initialize a new peer
        Given new peer creation order
        When it public key is valid
        Then the peer is initialized
*/
func TestInitializePeer(t *testing.T) {
	publicKey := GetValidPublicKey()
	ts := time.Now().Unix()

	peer := usecases.CreateNewPeer(publicKey, "127.0.0.1")

	assert.Equal(t, net.ParseIP("127.0.0.1"), peer.IP, "IP is not 127.0.0.1")
	assert.Equal(t, publicKey, peer.PublicKey, "Public key is not %s", publicKey)
	assert.Equal(t, entities.BootstrapingState, peer.Details.State, "Peer must boostraping")
	assert.Equal(t, ts, peer.Heartbeat.GenerationTime.Unix(), "Generation time is not %d", ts)

	time.Sleep(time.Second * 1)
	assert.True(t, peer.Heartbeat.GenerationTime.Second() < time.Now().Second(), "Peer heartbeat generation seconds must be lower than current seconds")
}
