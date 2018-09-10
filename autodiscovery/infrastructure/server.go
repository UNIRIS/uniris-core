package infrastructure

import (
	"fmt"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/adapters/gossip"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"

	"google.golang.org/grpc"
)

const (
	GrpcPort = 3545
)

//StartServer initiate an HTTP server with GRPC service
func StartServer(peerRepo repositories.PeerRepository) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", GrpcPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	gossip.RegisterGossipServiceServer(grpcServer, &gossip.Server{
		PeerRepo: peerRepo,
	})
	grpcServer.Serve(lis)
	return nil
}
