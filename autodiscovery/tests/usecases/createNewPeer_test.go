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
	ip := net.ParseIP("127.0.0.1")
	publicKey := GetValidPublicKey()

	peer := usecases.CreateNewPeer(ip, publicKey)

	assert.Equal(t, net.ParseIP("127.0.0.1"), peer.IP, "IP is not 127.0.0.1")
	assert.Equal(t, publicKey, peer.PublicKey, "Public key is not %s", publicKey)
	assert.True(t, peer.Heartbeat.GenerationTime.Second() < time.Now().Second(), "Peer heartbeat generation seconds must be lower than current seconds")
	assert.Equal(t, entities.BootstrapingState, peer.Details.State, "Peer must boostraping")
}

func GetValidPublicKey() []byte {
	return []byte("0448fe7dde9ce2151991abfba8f07ccfbd153419e3fd218357b2166d9811b02e5ad9cdfb6dba299e92dfcb954f57fb9188c5835b22c6b48d708f873c9e61da50ca")
}
