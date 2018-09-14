package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Check if an slice contains a peer
	Given a list of peers
	When we provide a peer to compare
	Then we retrieve the existing peer and the flag the peer exist
*/
func TestContainsPeer(t *testing.T) {

	peers := make([]domain.Peer, 0)
	peers = append(peers, domain.Peer{PublicKey: []byte("key"), IP: net.ParseIP("127.0.0.1")})

	exist, knownPeer := ContainsPeer(peers, domain.Peer{PublicKey: []byte("key")})
	assert.NotNil(t, knownPeer)
	assert.True(t, exist)
	assert.Equal(t, "127.0.0.1", knownPeer.IP.String())
}

/*
Scenario: Check if an slice not contains a peer
	Given a list of peers
	When we provide a peer to compare
	Then we not retrieve because the peer does not exist the provided list
*/
func TestNotContainsPeer(t *testing.T) {

	peers := make([]domain.Peer, 0)
	peers = append(peers, domain.Peer{PublicKey: []byte("key"), IP: net.ParseIP("127.0.0.1")})

	exist, _ := ContainsPeer(peers, domain.Peer{PublicKey: []byte("other")})
	assert.False(t, exist)
}
