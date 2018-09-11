package gossip

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//FormatPeerToGrpc converts a peer entity to GRPC peer message
func FormatPeerToGrpc(peer *entities.Peer) *Peer {
	return &Peer{
		PublicKey: peer.PublicKey,
		IP:        peer.IP.String(),
		Port:      int32(peer.Port),
		HeartbeatState: &Peer_Heartbeat{
			GenerationTime: peer.Heartbeat.GenerationTime.Unix(),
			ElapsedBeats:   peer.Heartbeat.ElapsedBeats,
		},
		AppState: &Peer_PeerAppState{
			CPULoad:       peer.AppState.CPULoad,
			FreeDiskSpace: float32(peer.AppState.FreeDiskSpace),
			GeoCoordinates: &Peer_PeerAppState_Coordinates{
				Lat: float32(peer.AppState.GeoCoordinates.Lat),
				Lon: float32(peer.AppState.GeoCoordinates.Lon),
			},
			IOWaitRate: float32(peer.AppState.IOWaitRate),
			P2PFactor:  int32(peer.AppState.P2PFactor),
			Version:    peer.AppState.Version,
			State:      Peer_PeerAppState_PeerState(peer.AppState.State),
		},
	}
}

//FormatPeerToDomain converts an GRPC message peer to peer entity
func FormatPeerToDomain(peer Peer) *entities.Peer {
	return &entities.Peer{
		PublicKey: []byte(peer.PublicKey),
		IP:        net.ParseIP(peer.IP),
		Port:      int(peer.Port),
		Heartbeat: entities.PeerHeartbeat{
			GenerationTime: time.Unix(peer.HeartbeatState.GenerationTime, 0),
			ElapsedBeats:   peer.HeartbeatState.ElapsedBeats,
		},
		AppState: entities.PeerAppState{
			CPULoad:       peer.AppState.CPULoad,
			FreeDiskSpace: float64(peer.AppState.FreeDiskSpace),
			GeoCoordinates: entities.Coordinates{
				Lat: float64(peer.AppState.GeoCoordinates.Lat),
				Lon: float64(peer.AppState.GeoCoordinates.Lon),
			},
			IOWaitRate: float64(peer.AppState.IOWaitRate),
			P2PFactor:  int(peer.AppState.P2PFactor),
			Version:    peer.AppState.Version,
			State:      entities.PeerState(peer.AppState.State),
		},
	}
}
