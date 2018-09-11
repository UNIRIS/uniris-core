package tests

import (
	"net"
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Select random peer from a list of peer with only one peer
	Given a list of peer within one peer
	When we want to pick a random peer
	Then we returns the only peer
*/
func TestSelectRandomPeerWithOnePeer(t *testing.T) {
	peers := []*entities.Peer{
		&entities.Peer{
			IP: net.ParseIP("127.0.0.1"),
		},
	}

	peer := usecases.SelectRandomPeer(peers)
	assert.NotNil(t, peer)
	assert.Equal(t, "127.0.0.1", peer.IP.String(), "IP address must be 127.0.0.1")
}

/*
Scenario: Select random peer from a list of peer within many peers
	Given a list of peer within many peers
	When we want to pick a random peer
	Then we returns a random peer
*/
func TestSelectRandomPeerWithManyPeers(t *testing.T) {
	peers := []*entities.Peer{
		&entities.Peer{
			IP: net.ParseIP("127.0.0.1"),
		},
		&entities.Peer{
			IP: net.ParseIP("30.0.0.1"),
		},
	}

	peer := usecases.SelectRandomPeer(peers)
	assert.NotNil(t, peer)
	assert.True(t, peer.IP.String() == "127.0.0.1" || peer.IP.String() == "30.0.0.1")
}
