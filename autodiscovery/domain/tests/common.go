package tests

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

func GetValidPublicKey() []byte {
	return []byte("0448fe7dde9ce2151991abfba8f07ccfbd153419e3fd218357b2166d9811b02e5ad9cdfb6dba299e92dfcb954f57fb9188c5835b22c6b48d708f873c9e61da50ca")
}
func GetSecondValidPublicKey() []byte {
	return []byte("0448fe7dde9ce2151991abfba8f07ccfbd153419e3fd218357b2166d9811b02e5ad9cdfb6dba299e92dfcb954f57fb9188c5835b22c6b48d708f873c9e61da50cb")
}

func GetThirdValidPublicKey() []byte {
	return []byte("0448fe7dde9ce2151991abfba8f07ccfbd153419e3fd218357b2166d9811b02e5ad9cdfb6dba299e92dfcb954f57fb9188c5835b22c6b48d708f873c9e61da50cc")
}

func GetRepo() repositories.PeerRepository {
	return &PeerRepository{}
}

type SeedLoader struct{}

func (s SeedLoader) GetSeedPeers() ([]*entities.Peer, error) {
	seeds := make([]*entities.Peer, 0)
	seeds = append(seeds, usecases.CreateNewPeer(GetValidPublicKey(), "127.0.0.1"))
	return seeds, nil
}

type GeolocService struct{}

func (geo *GeolocService) Lookup() (services.GeoLoc, error) {
	return services.GeoLoc{
		IP:  net.ParseIP("127.0.0.1"),
		Lat: 2.33,
		Lon: 64.20,
	}, nil
}

type FullGossipService struct{}

func (s *FullGossipService) DiscoverPeers(destPeer entities.Peer, knownPeers []*entities.Peer) ([]*entities.Peer, error) {
	return []*entities.Peer{
		usecases.CreateNewPeer(GetSecondValidPublicKey(), "30.10.200.50"),
		usecases.CreateNewPeer(GetThirdValidPublicKey(), "50.250.111.32"),
	}, nil
}

type SameGossipService struct{}

func (s *SameGossipService) DiscoverPeers(destPeer entities.Peer, knownPeers []*entities.Peer) ([]*entities.Peer, error) {
	return []*entities.Peer{
		usecases.CreateNewPeer(GetValidPublicKey(), "127.0.0.1"),
	}, nil
}
