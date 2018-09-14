package infrastructure

import (
	"fmt"
	"net"
)

//NewNetListener creates HTTP listener
func NewNetListener(port int) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}
	return lis, nil

	// grpcServer := grpc.NewServer()

	// log.Printf("Server listening on port %d", port)
	// err = grpcServer.Serve(lis)
	// if err != nil {
	// 	return err
	// }
	// return nil
}
