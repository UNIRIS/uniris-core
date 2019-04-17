package amqp

import (
	"encoding/json"

	"github.com/uniris/uniris-core/pkg/discovery"
)

type node struct {
	Identity nodeIdentity `json:"identity"`
	AppState nodeAppState `json:"app_state"`
}

type nodeIdentity struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	PublicKey []byte `json:"public_key"`
}

type nodeAppState struct {
	Version              string       `json:"version"`
	Status               int          `json:"status"`
	CPULoad              string       `json:"cpu_load"`
	FreeDiskSpace        float64      `json:"free_disk_space"`
	GeoPosition          peerPosition `json:"geo_position"`
	P2PFactor            int          `json:"p2p_factor"`
	ReachablePeersNumber int          `json:"reachable_peers_number"`
}

type peerPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func marshalNode(p discovery.Peer) ([]byte, error) {
	pk, err := p.Identity().PublicKey().Marshal()
	if err != nil {
		return nil, err
	}
	return json.Marshal(node{
		Identity: nodeIdentity{
			IP:        p.Identity().IP().String(),
			Port:      p.Identity().Port(),
			PublicKey: pk,
		},
		AppState: nodeAppState{
			Version:       p.AppState().Version(),
			Status:        int(p.AppState().Status()),
			CPULoad:       p.AppState().CPULoad(),
			FreeDiskSpace: p.AppState().FreeDiskSpace(),
			GeoPosition: peerPosition{
				Latitude:  p.AppState().GeoPosition().Latitude(),
				Longitude: p.AppState().GeoPosition().Longitude(),
			},
			P2PFactor:            p.AppState().P2PFactor(),
			ReachablePeersNumber: p.AppState().ReachablePeersNumber(),
		},
	})
}
