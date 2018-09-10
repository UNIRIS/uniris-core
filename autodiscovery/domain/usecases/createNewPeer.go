package usecases

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//CreateNewPeer initializes a new peer
func CreateNewPeer(publicKey []byte, ip string) *entities.Peer {
	return &entities.Peer{
		PublicKey: publicKey,
		IP:        net.ParseIP(ip),
		Heartbeat: entities.PeerHeartbeat{
			GenerationTime: time.Now(),
		},
		Details: entities.PeerDetails{
			State: entities.BootstrapingState,
		},
	}
}
