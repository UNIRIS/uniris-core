package adapters

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/adapters/protobuf"
	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"google.golang.org/grpc"
)

//GrpcGossipBroker implements the gossip broker using grpc
type GrpcGossipBroker struct{}

//SendSyn callS the Synchronize protobuf method to retrieve unknown peers (SYN handshake)
func (s GrpcGossipBroker) SendSyn(synReq domain.SynRequest) (synAck domain.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s:%d", synReq.Receiver.IP.String(), synReq.Receiver.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return
	}

	client := protobuf.NewGossipClient(conn)
	resp, err := client.Synchronize(context.Background(), protobuf.BuildSynRequest(synReq))
	if err != nil {
		return
	}

	synAck = protobuf.BuildDomainSynAckResponse(resp)
	return
}

//SendAck calls the Acknoweledge protobuf method to send detailed peers requested
func (s GrpcGossipBroker) SendAck(ackReq domain.AckRequest) error {
	serverAddr := fmt.Sprintf("%s:%d", ackReq.Receiver.IP.String(), ackReq.Receiver.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := protobuf.NewGossipClient(conn)
	_, err = client.Acknowledge(context.Background(), protobuf.BuildAckRequest(ackReq))
	if err != nil {
		return err
	}
	return nil
}
