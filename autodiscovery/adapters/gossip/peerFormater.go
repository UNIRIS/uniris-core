package gossip

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

func FormatPeerToGrpc(peer entities.Peer) *Peer {
	return &Peer{
		PublicKey: peer.PublicKey,
		IP:        peer.IP.String(),
		HeartbeatState: &Heartbeat{
			GenerationTime: peer.Heartbeat.GenerationTime.Unix(),
			ElapsedBeats:   peer.Heartbeat.ElapsedBeats,
		},
		Details: &PeerDetail{
			CPULoad:       peer.Details.CPULoad,
			FreeDiskSpace: float32(peer.Details.FreeDiskSpace),
			GeoCoordinates: &PeerDetail_Coordinates{
				Lat: float32(peer.Details.GeoCoordinates.Lat),
				Lon: float32(peer.Details.GeoCoordinates.Lon),
			},
			IOWaitRate: float32(peer.Details.IOWaitRate),
			P2PFactor:  int32(peer.Details.P2PFactor),
			Version:    peer.Details.Version,
			State:      PeerDetail_PeerState(peer.Details.State),
		},
	}
}

func FormatPeerToDomain(peer Peer) *entities.Peer {
	return &entities.Peer{
		PublicKey: []byte(peer.PublicKey),
		IP:        net.ParseIP(peer.IP),
		Heartbeat: entities.PeerHeartbeat{
			GenerationTime: time.Unix(peer.HeartbeatState.GenerationTime, 0),
			ElapsedBeats:   peer.HeartbeatState.ElapsedBeats,
		},
		Details: entities.PeerDetails{
			CPULoad:       peer.Details.CPULoad,
			FreeDiskSpace: float64(peer.Details.FreeDiskSpace),
			GeoCoordinates: entities.Coordinates{
				Lat: float64(peer.Details.GeoCoordinates.Lat),
				Lon: float64(peer.Details.GeoCoordinates.Lon),
			},
			IOWaitRate: float64(peer.Details.IOWaitRate),
			P2PFactor:  int(peer.Details.P2PFactor),
			Version:    peer.Details.Version,
			State:      entities.PeerState(peer.Details.State),
		},
	}
}
