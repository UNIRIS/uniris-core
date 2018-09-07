package repositories

//PeerRepository implements the IPeerRepository on memory
type PeerRepository struct {
	// db database.KeyValueStorage
}

//GetPeers retrieves the peers stored locally as map
// func (ps PeerRepository) GetPeers() (map[string]entities.Peer, error) {

// }

// //AddPeer stores a peer locally
// func (ps *PeerRepository) AddPeer(p entities.Peer) error {
// 	ps.peerStorage = append(ps.peerStorage, p)
// 	return nil
// }

// //UpdatePeer changes a peer locally
// func (ps *PeerRepository) UpdatePeer(p entities.Peer) error {
// 	for _, peer := range ps.peerStorage {
// 		if peer.IP.Equal(p.IP) {
// 			peer = p
// 		}
// 	}
// 	return nil
// }
