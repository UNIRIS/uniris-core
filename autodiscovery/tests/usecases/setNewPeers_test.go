package tests

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenerio: Set new peers
	Given a set of new peers
	When we add the peers
	Then the repository contains the new peers
*/
func TestSetNewPeers(t *testing.T) {
	repo := &PeerRepository{}
	peers := make([]entities.Peer, 0)
	for i := 0; i < 10; i++ {
		newPeer := entities.Peer{
			IP: net.ParseIP(fmt.Sprintf("35.165.78.20%d", i)),
		}
		peers = append(peers, newPeer)
	}
	err := usecases.SetNewPeers(repo, peers)
	assert.Nil(t, err)

	knownPeers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)
	assert.NotEmpty(t, knownPeers)
	assert.Equal(t, 10, len(knownPeers))
}

/*
Scenario: Update an older peer with a yougest peer
	Given a known peer
	When received update for a yougest peer
	Then the peer is updated in the database
*/
func TestUpdatePeerWithYoungestPeer(t *testing.T) {
	repo := &PeerRepository{}

	peer := usecases.CreateNewPeer(net.ParseIP("127.0.0.1"), GetValidPublicKey())
	err := usecases.SetNewPeers(repo, []entities.Peer{peer})
	assert.Nil(t, err)

	time.Sleep(time.Second * 2)
	peer.RefreshHearbeat()
	peer.Details.State = entities.FaultyState

	err = usecases.SetNewPeers(repo, []entities.Peer{peer})
	assert.Nil(t, err)
	peers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)

	assert.Equal(t, entities.FaultyState, peers[0].Details.State)
}

/*
Scenario: Update an peer with an older peer
	Given a known peer
	When received update for an older peer
	Then the peer is not updated in the database
*/
func TestUpdatePeerWithOlderPeer(t *testing.T) {
	repo := &PeerRepository{}

	olderPeer := usecases.CreateNewPeer(net.ParseIP("127.0.0.1"), GetValidPublicKey())
	olderPeer.Details.State = entities.FaultyState
	time.Sleep(time.Second * 2)

	peer := usecases.CreateNewPeer(net.ParseIP("127.0.0.1"), GetValidPublicKey())
	err := usecases.SetNewPeers(repo, []entities.Peer{peer})
	assert.Nil(t, err)

	err = usecases.SetNewPeers(repo, []entities.Peer{olderPeer})
	assert.Nil(t, err)
	peers, err := usecases.ListKnownPeers(repo)
	assert.Nil(t, err)

	assert.Equal(t, entities.BootstrapingState, peers[0].Details.State)
}
