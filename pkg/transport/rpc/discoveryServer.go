package rpc

import (
	"context"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type discoverySrv struct {
	db         discovery.Database
	notif      discovery.Notifier
	logger     logging.Logger
	publickey  crypto.PublicKey
	privatekey crypto.PrivateKey
	keyReader  shared.KeyReader
}

//NewDiscoveryServer creates a new GRPC server for the discovery service
func NewDiscoveryServer(db discovery.Database, n discovery.Notifier, l logging.Logger, pubK crypto.PublicKey, privK crypto.PrivateKey, keyR shared.KeyReader) api.DiscoveryServiceServer {
	return &discoverySrv{
		db:         db,
		notif:      n,
		logger:     l,
		publickey:  pubK,
		privatekey: privK,
		keyReader:  keyR,
	}
}

func (s discoverySrv) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynResponse, error) {
	s.logger.Debug("SYN REQUEST - " + time.Unix(req.Timestamp, 0).String())

	pub, err := crypto.ParsePublicKey(req.PublicKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if !s.keyReader.IsAuthorizedNode(pub) {
		return nil, status.New(codes.Internal, ("Not Authorized")).Err()
	}

	if !pub.Verify(req.PublicKey, req.Signature) {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	receivedPeers := make([]discovery.Peer, 0)
	for _, p := range req.KnownPeers {
		fp, err := formatPeerDigest(p)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		receivedPeers = append(receivedPeers, fp)
	}

	localPeers, err := s.db.DiscoveredPeers()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	requestedPeers := make([]*api.PeerIdentity, 0)
	for _, p := range discovery.ComparePeers(localPeers, receivedPeers) {
		fp, err := formatPeerIdentityAPI(p)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		requestedPeers = append(requestedPeers, fp)
	}

	discoveries := make([]*api.PeerDiscovered, 0)
	for _, p := range discovery.ComparePeers(receivedPeers, localPeers) {
		fp, err := formatPeerDiscoveredAPI(p)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		discoveries = append(discoveries, fp)
	}

	return &api.SynResponse{
		RequestedPeers:  requestedPeers,
		DiscoveredPeers: discoveries,
		Timestamp:       time.Now().Unix(),
	}, nil
}

func (s discoverySrv) Acknowledge(ctx context.Context, req *api.AckRequest) (*api.AckResponse, error) {
	s.logger.Debug("ACK REQUEST - " + time.Unix(req.Timestamp, 0).String())

	pub, err := crypto.ParsePublicKey(req.PublicKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if !s.keyReader.IsAuthorizedNode(pub) {
		return nil, status.New(codes.Internal, ("Not Authorized")).Err()
	}

	if !pub.Verify(req.PublicKey, req.Signature) {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	newPeers := make([]discovery.Peer, 0)
	for _, p := range req.RequestedPeers {
		fp, err := formatPeerDiscovered(p)
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		newPeers = append(newPeers, fp)
	}

	for _, p := range newPeers {
		//TODO verify that the discovered peers are authorized.
		//TODO verify that the discovered peers info are authentified
		if err := s.db.WriteDiscoveredPeer(p); err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		if err := s.notif.NotifyDiscovery(p); err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &api.AckResponse{
		Timestamp: time.Now().Unix(),
	}, nil
}
