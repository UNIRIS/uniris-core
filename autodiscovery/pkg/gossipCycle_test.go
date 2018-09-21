package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Run a cycle run and gossip with a peer
	Given a initator, a receiver peer and a list of known peers
	When we run a cycle to gossip with receiver peer
	Then we get some new peers from the receiver
*/
func TestRunCycle(t *testing.T) {

	init := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{}, 1)
	rec := NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)
	kp := []Peer{init, NewPeerDigest([]byte("key3"), net.ParseIP("20.0.0.1"), 3000)}
	msg := mockMessenger{}

	c := NewGossipCycle(init, rec, kp, msg)
	np, err := c.Run()
	assert.Nil(t, err)
	assert.NotEmpty(t, np)
	assert.Equal(t, "dkey", string(np[0].PublicKey()))
}

/*
Scenario: Run a cycle run and gossip with a peer and send back some unknown peers
	Given a run cycle started
	When we received a SYN ACK, we got some unknown peers from the receiver
	Then we send details from theses peers and returns without error
*/
func TestRunAckRequest(t *testing.T) {

	init := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{}, 1)
	rec := NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)

	kp1 := NewPeerDetailed(
		[]byte("key3"), net.ParseIP("20.0.0.1"), 3000, time.Now(), false,
		NewState("1.1", OkStatus, PeerPosition{}, "200.10.000", 500.20, 200.10, 1),
	)

	kp := []Peer{init, kp1}
	msg := mockAckMessenger{}

	c := NewGossipCycle(init, rec, kp, msg)
	np, err := c.Run()
	assert.Nil(t, err)
	assert.NotEmpty(t, np)
	assert.Equal(t, "dkey", string(np[0].PublicKey()))
}

type mockAckMessenger struct{}

func (m mockAckMessenger) SendSyn(r SynRequest) (*SynAck, error) {
	return &SynAck{
		NewPeers: []Peer{
			NewPeerDetailed([]byte("dkey"), net.ParseIP("10.0.0.1"), 3000, time.Now(), false, nil),
		},
		UnknownPeers: []Peer{
			NewPeerDigest([]byte("key3"), net.ParseIP("20.0.0.1"), 3000),
		},
	}, nil
}

func (m mockAckMessenger) SendAck(r AckRequest) error {
	return nil
}
