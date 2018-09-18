package main

import (
	"fmt"
	"net"
)

func startServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprint("localhost:%d", port))
	if err != nil {
		return err
	}

	grpcHandler := grpc.NewGrpcHandler()

	log.Printf("Server listening on %d", port)
	if err := grpcHandler.Serve(lis); err != nil {
		return err
	}
	return nil
}