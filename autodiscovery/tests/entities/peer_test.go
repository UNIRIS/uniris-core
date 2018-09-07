package tests

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

/*
Scenarion: Reset the peer heartbeat
	Given an initialized peer
	When the peer must send its state
	Then the peer resets its heartbeat generation time
*/
func TestResetHeartBeat(t *testing.T) {
	peer := &entities.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: GetValidPublicKey(),
	}

	peer.RefreshHearbeat()
	ts := time.Now().Unix()
	assert.Equal(t, ts, peer.Heartbeat.GenerationTime.Unix(), "Reset heartbeat generation time is not %d", ts)
}

/*
Scenario: Get the elapsed peer's heartbeats
        Given an initialized peer
        When 2 seconds elapsed
        Then peer's heartbeats must equal to 2
*/
func TestGetPeerElapsedHeartBeats(t *testing.T) {
	peer := &entities.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: GetValidPublicKey(),
	}
	peer.RefreshHearbeat()
	time.Sleep(time.Second * 2)

	heartbeats := peer.GetElapsedHeartbeats()
	assert.Equal(t, int64(2), heartbeats, "Peer heartbeats must be 2")
}

func GetValidPublicKey() []byte {
	return []byte("0448fe7dde9ce2151991abfba8f07ccfbd153419e3fd218357b2166d9811b02e5ad9cdfb6dba299e92dfcb954f57fb9188c5835b22c6b48d708f873c9e61da50ca")
}
