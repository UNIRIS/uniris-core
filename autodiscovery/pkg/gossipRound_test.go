package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Execute a syn request
	Given a initator, a receiver peer and a list of known peers
	When we run a gossip round to gossip with targeted peer
	Then we get some new peers from the receiver
*/
func TestRoundSpread(t *testing.T) {
	init := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{}, 1)
	tar := NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)
	kp := []Peer{init, NewPeerDigest([]byte("key3"), net.ParseIP("20.0.0.1"), 3000)}
	spr := mockSpreader{}

	r := NewGossipRound(init, tar)
	err := r.Spread(kp, spr)
	assert.Nil(t, err)
	assert.NotEmpty(t, r.discoveredPeers)
	assert.Equal(t, "dkey", string(r.discoveredPeers[0].PublicKey()))
}

type mockSpreader struct{}

func (m mockSpreader) SendSyn(r SynRequest) (*SynAck, error) {
	return &SynAck{
		NewPeers: []Peer{
			NewPeerDetailed([]byte("dkey"), net.ParseIP("10.0.0.1"), 3000, time.Now(), false, nil),
		},
		UnknownPeers: []Peer{
			NewPeerDigest([]byte("key3"), net.ParseIP("20.0.0.1"), 3000),
		},
	}, nil
}

func (m mockSpreader) SendAck(r AckRequest) error {
	return nil
}
