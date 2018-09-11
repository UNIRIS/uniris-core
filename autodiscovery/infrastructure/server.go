package infrastructure

import (
	"fmt"
	"log"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/adapters/transport"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"

	"google.golang.org/grpc"
)

//StartServer initiates an HTTP server with GRPC service
func StartServer(peerRepo repositories.PeerRepository, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	transport.RegisterGossipGrpcServiceServer(grpcServer, &transport.GrpcServer{
		PeerRepo: peerRepo,
	})
	log.Println(fmt.Sprintf("Server listening on port %d", port))
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
