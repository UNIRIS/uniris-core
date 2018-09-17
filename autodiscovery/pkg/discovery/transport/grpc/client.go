package grpc

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/gossip"
	"google.golang.org/grpc"
)

type GrpcClient struct{}

//SendSyn callS the Synchronize protobuf method to retrieve unknown peers (SYN handshake)
func (g GrpcClient) SendSyn(synReq gossip.SynRequest) (synAck gossip.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s:%d", synReq.Receiver.IP.String(), synReq.Receiver.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return
	}

	client := protobuf.NewGossipClient(conn)
	resp, err := client.Synchronize(context.Background(), BuildSynRequest(synReq))
	if err != nil {
		return
	}

	synAck = protobuf.BuildDomainSynAckResponse(resp)
	return
}

//SendAck calls the Acknoweledge protobuf method to send detailed peers requested
func (g GrpcClient) SendAck(ackReq gossip.AckRequest) error {
	serverAddr := fmt.Sprintf("%s:%d", ackReq.Receiver.IP.String(), ackReq.Receiver.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := protobuf.NewGossipClient(conn)
	_, err = client.Acknowledge(context.Background(), BuildAckRequest(ackReq))
	if err != nil {
		return err
	}
	return nil
}
