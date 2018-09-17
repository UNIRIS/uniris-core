package mock

import "github.com/uniris/uniris-core/autodiscovery/pkg/discovery"

//MockPeerRepository mocks the peer storage
type MockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.SeedPeer
}

func NewRepository() MockPeerRepository {
	return MockPeerRepository{}
}

func (r MockPeerRepository) AddPeer(p discovery.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r MockPeerRepository) AddSeed(s discovery.SeedPeer) error {
	r.seeds = append(r.seeds, s)
	return nil
}
