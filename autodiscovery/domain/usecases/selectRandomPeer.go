package usecases

import (
	"math/rand"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//SelectRandomPeer pick a known random peer from a list of peers
func SelectRandomPeer(peers []*entities.Peer) *entities.Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return peers[0]
}
