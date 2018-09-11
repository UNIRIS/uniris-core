package tests

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

/*
Scenario: Refresh the peer's details
	Given a peer
	When we refresh its detail
	Then peer's detail are refreshed such as hearbeats
*/
func TestRefreshPeerDetails(t *testing.T) {
	peer := &entities.Peer{
		IP: net.ParseIP("127.0.0.1"),
		Heartbeat: entities.PeerHeartbeat{
			GenerationTime: time.Now(),
		},
	}
	time.Sleep(2 * time.Second)
	peer.UpdateElapsedHeartbeats()

	assert.Equal(t, peer.Heartbeat.ElapsedBeats, int64(2), "Elapsed beats must be 2 seconds")

	//TODO: test other peer's details
}
