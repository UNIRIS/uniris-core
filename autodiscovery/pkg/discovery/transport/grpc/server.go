package grpc

import (
	"fmt"
	"net"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"google.golang.org/grpc"
)

func StartServer(ip net.IP, port int) error {
	lis, err := net.Listen("tcp", fmt.Printf("%s:d"), ip.String(), port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServer(grpcServer, GrpcHandler{})
}
