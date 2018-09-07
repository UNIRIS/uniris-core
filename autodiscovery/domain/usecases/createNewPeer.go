package usecases

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//CreateNewPeer initializes a new peer
func CreateNewPeer(ip net.IP, publicKey []byte) entities.Peer {
	peer := entities.Peer{
		IP:        ip,
		PublicKey: publicKey,
		Details: entities.PeerDetails{
			State: entities.BootstrapingState,
		},
	}
	return peer
}
