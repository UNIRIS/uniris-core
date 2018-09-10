package tests

import (
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Select random peer from an empty known peers
	Given no known peers
	When we want to pick a random peer
	Then a error is returned telling there is no peers
*/
func TestSelectRandomPeerWithoutPeers(t *testing.T) {
	peerRepo := &PeerRepository{}
	_, err := usecases.SelectRandomPeer(peerRepo)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Random peer selection cannot be done with no known peers")
}

/*
Scenario: Select random peer from an known peer's list with 1 peer
	Given known peers which contains one peer
	When we want to pick a random peer
	Then we get the unique peer
*/
func TestSelectRandomWithOneKnownPeer(t *testing.T) {
	repo := &PeerRepository{}
	err := usecases.LoadSeedPeers(&SeedLoader{}, repo)
	assert.Nil(t, err)
	randomPeer, err := usecases.SelectRandomPeer(repo)
	assert.Nil(t, err)
	assert.NotNil(t, randomPeer, "Random peer cannot be null")
	assert.Equal(t, "127.0.0.1", randomPeer.IP.String(), "Random peer IP must be 127.0.0.1")
}

/*
Scenario: Select random peer from an known peer's list with 3 peers
	Given known peers which contains one peer
	When we want to pick a random peer
	Then we get a random peer
*/
func TestSelectRandomWithThreeKnownPeers(t *testing.T) {
	repo := &PeerRepository{}

	peer1 := usecases.CreateNewPeer(GetValidPublicKey(), "127.0.0.1")
	peer2 := usecases.CreateNewPeer(GetSecondValidPublicKey(), "30.50.230.10")
	peer3 := usecases.CreateNewPeer(GetThirdValidPublicKey(), "50.10.200.220")

	err := usecases.SetNewPeers(repo, []*entities.Peer{peer1, peer2, peer3})
	assert.Nil(t, err)

	randomPeer, err := usecases.SelectRandomPeer(repo)
	assert.Nil(t, err)
	assert.NotNil(t, randomPeer, "Random peer cannot be null")
}
