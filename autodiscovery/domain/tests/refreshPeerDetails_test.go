package tests

import (
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Refresh the peer's details
	Given a peer
	When we refresh its detail
	Then peer's detail are refreshed such as hearbeats
*/
func TestRefreshPeerDetails(t *testing.T) {

	repo := GetRepo()
	usecases.StartupPeer(repo, &GeolocService{}, GetValidPublicKey(), 3545)
	time.Sleep(2 * time.Second)

	usecases.RefreshSelfPeer(repo)
	peer, _ := repo.GetOwnPeer()
	assert.Equal(t, peer.Heartbeat.ElapsedBeats, int64(2), "Elapsed beats must be 2 seconds")

	//TODO: test other peer's details
}
