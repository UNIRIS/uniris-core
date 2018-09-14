package adapters

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/uniris/uniris-core/autodiscovery/adapters/protobuf"
	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"github.com/uniris/uniris-core/autodiscovery/core/ports"
	"github.com/uniris/uniris-core/autodiscovery/core/usecases"
	"google.golang.org/grpc"
)

//GrpcHandlers implements the GRPC server
type GrpcHandlers struct {
	PeerRepo            ports.PeerRepository
	ConfigurationReader ports.ConfigurationReader
	MetricReader        ports.MetricReader
}

//NewGRPC create grpc interceptor for protobuf messages
func NewGRPC(r ports.PeerRepository, c ports.ConfigurationReader, m ports.MetricReader) *grpc.Server {
	grpcServer := grpc.NewServer()
	protobuf.RegisterGossipServer(grpcServer, GrpcHandlers{
		PeerRepo:            r,
		ConfigurationReader: c,
		MetricReader:        m,
	})
	return grpcServer
}

//Synchronize implements the protobuf Synchronize request handler
func (s GrpcHandlers) Synchronize(ctx context.Context, req *protobuf.SynRequest) (*protobuf.SynAck, error) {

	initiator := protobuf.ToDomain(req.Initiator)
	receiver := protobuf.ToDomain(req.Receiver)
	receivedPeers := protobuf.ToDomainBulk(req.KnownPeers)

	//TODO: check the identity of the req.Initiator peer

	unknownPeers, err := usecases.GetUnknownPeers(s.PeerRepo, receivedPeers)
	if err != nil {
		return nil, err
	}

	newPeers, err := usecases.ProvideNewPeers(s.PeerRepo, receivedPeers)
	if err != nil {
		return nil, err
	}

	//Refresh owned peer to send it back
	ownedPeer, err := s.PeerRepo.GetOwnedPeer()
	if err != nil {
		return nil, err
	}
	usecases.RefreshPeer(&ownedPeer, s.PeerRepo, s.ConfigurationReader, s.MetricReader)
	newPeers = append(newPeers, ownedPeer)

	synAck := domain.NewSynAck(initiator, receiver, newPeers, unknownPeers)
	return protobuf.BuildProtoSynAckResponse(synAck), nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (s GrpcHandlers) Acknowledge(ctx context.Context, req *protobuf.AckRequest) (*empty.Empty, error) {
	//Store the peers requested
	for _, peer := range req.RequestedPeers {
		domainPeer := protobuf.ToDomain(peer)
		exist, err := s.PeerRepo.ContainsPeer(domainPeer)
		if err != nil {
			return nil, err
		}
		if exist {
			s.PeerRepo.UpdatePeer(domainPeer)
		} else {
			s.PeerRepo.InsertPeer(domainPeer)
		}
	}
	return nil, nil
}
