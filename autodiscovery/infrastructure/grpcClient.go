package infrastructure

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/adapters"
	"github.com/uniris/uniris-core/autodiscovery/domain"

	"github.com/uniris/uniris-core/autodiscovery/infrastructure/proto"
	"google.golang.org/grpc"
)

//GrpcClient implements the protobuf gossip requests
type GrpcClient struct{}

//SendSynchro call the Synchronize protobuf method to retrieve unknown peers (SYN handshake)
func (s GrpcClient) SendSynchro(synReq domain.SynRequest) (*domain.SynAck, error) {
	serverAddr := fmt.Sprintf("%s:%d", synReq.TargetIP.String(), synReq.TargetPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := proto.NewGossipClient(conn)
	resp, err := client.Synchronize(context.Background(), adapters.BuildSynRequest(synReq.Sender, synReq.KnownPeers))
	if err != nil {
		return nil, err
	}

	return adapters.BuildDomainSynAckResponse(resp), nil
}

//SendAcknowledgement call the Acknoweledge protobuf method to send detailed peers requested
func (s GrpcClient) SendAcknowledgement(ackReq domain.AckRequest) error {
	serverAddr := fmt.Sprintf("%s:%d", ackReq.TargetIP.String(), ackReq.TargetPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := proto.NewGossipClient(conn)
	_, err = client.Acknowledge(context.Background(), adapters.BuildAckRequest(ackReq.DetailedKnownPeers))
	if err != nil {
		return err
	}
	return nil
}
