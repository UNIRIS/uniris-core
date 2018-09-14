package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/uniris/uniris-core/autodiscovery/adapters"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure/proto"
	"github.com/uniris/uniris-core/autodiscovery/usecases"
	"github.com/uniris/uniris-core/autodiscovery/usecases/ports"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"

	"google.golang.org/grpc"
)

//StartServer initiates an HTTP server with GRPC service
func StartServer(peerRepo repositories.PeerRepository, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	identityChecker := new(adapters.IdentityChecker)
	proto.RegisterGossipServer(grpcServer, GrpcServer{
		IdentityChecker: identityChecker,
		PeerRepo:        peerRepo,
	})
	log.Printf("Server listening on port %d", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

//GrpcServer implements the GRPC server
type GrpcServer struct {
	PeerRepo        repositories.PeerRepository
	IdentityChecker ports.PeerIdentityChecker
}

//Synchronize implements the protobuf Synchronize request handler
func (s GrpcServer) Synchronize(ctx context.Context, req *proto.SynRequest) (*proto.SynAck, error) {
	if !s.IdentityChecker.IsPublicKeyAuthorized(req.Sender.PublicKey) {
		return nil, errors.New("Unauthorized public key")
	}

	receivedPeers := adapters.ToDomainBulk(req.KnownPeers)

	detailedPeersRequested, err := usecases.GetUnknownPeers(s.PeerRepo, receivedPeers)
	if err != nil {
		return nil, err
	}

	//TODO: refresh its own state
	newPeers, err := usecases.ProvideNewPeers(s.PeerRepo, receivedPeers)
	if err != nil {
		return nil, err
	}

	return adapters.BuildProtoSynAckResponse(newPeers, detailedPeersRequested), nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (s GrpcServer) Acknowledge(ctx context.Context, req *proto.AckRequest) (*empty.Empty, error) {
	//Store the peers received
	for _, peer := range req.DetailedKnownPeers {
		if err := usecases.StorePeer(s.PeerRepo, adapters.ToDomain(peer)); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
