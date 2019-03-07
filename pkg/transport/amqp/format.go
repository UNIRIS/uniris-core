package amqp

import (
	"encoding/json"
	"net"
	"time"

	"github.com/uniris/uniris-core/pkg/discovery"
)

type discoveredPeer struct {
	Identity       peerIdentity      `json:"identity"`
	HeartbeatState peerHearbeatState `json:"heartbeat_state"`
	AppState       peerAppState      `json:"app_state"`
}

type peerIdentity struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	PublicKey string `json:"public_key"`
}

type peerHearbeatState struct {
	GenerationTime    time.Time `json:"generation_time"`
	ElapsedHeartbeats int64     `json:"elapsed_heartbeats"`
}

type peerAppState struct {
	Version              string               `json:"version"`
	Status               discovery.PeerStatus `json:"status"`
	CPULoad              string               `json:"cpu_load"`
	FreeDiskSpace        float64              `json:"free_disk_space"`
	GeoPosition          peerPosition         `json:"geo_position"`
	P2PFactor            int                  `json:"p2p_factor"`
	ReachablePeersNumber int                  `json:"reachable_peers_number"`
}

type peerPosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func marshalDiscoveredPeer(p discovery.Peer) ([]byte, error) {
	return json.Marshal(discoveredPeer{
		Identity: peerIdentity{
			IP:        p.Identity().IP().String(),
			Port:      p.Identity().Port(),
			PublicKey: p.Identity().PublicKey(),
		},
		HeartbeatState: peerHearbeatState{
			GenerationTime:    p.HeartbeatState().GenerationTime(),
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
		},
		AppState: peerAppState{
			Version:       p.AppState().Version(),
			Status:        p.AppState().Status(),
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

func unmarshalPeer(b []byte) (discovery.Peer, error) {
	var p discoveredPeer
	if err := json.Unmarshal(b, &p); err != nil {
		return discovery.Peer{}, err
	}

	return discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP(p.Identity.IP), p.Identity.Port, p.Identity.PublicKey),
		discovery.NewPeerHeartbeatState(p.HeartbeatState.GenerationTime, p.HeartbeatState.ElapsedHeartbeats),
		discovery.NewPeerAppState(p.AppState.Version, p.AppState.Status, p.AppState.GeoPosition.Longitude, p.AppState.GeoPosition.Latitude, p.AppState.CPULoad, p.AppState.FreeDiskSpace, p.AppState.P2PFactor, p.AppState.ReachablePeersNumber),
	), nil
}

func unmarshalPeerIdentity(b []byte) (discovery.PeerIdentity, error) {
	var pi peerIdentity
	if err := json.Unmarshal(b, &pi); err != nil {
		return discovery.PeerIdentity{}, err
	}

	return discovery.NewPeerIdentity(net.ParseIP(pi.IP), pi.Port, pi.PublicKey), nil
}

func marshalPeerIdentity(p discovery.PeerIdentity) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"ip":         p.IP().String(),
		"port":       p.Port(),
		"public_key": p.PublicKey(),
	})
}
