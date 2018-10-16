package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	adding "github.com/uniris/uniris-core/datamining/pkg/adding"
	file "github.com/uniris/uniris-core/datamining/pkg/file"
	listing "github.com/uniris/uniris-core/datamining/pkg/listing"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	internalrpc "github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
	validating "github.com/uniris/uniris-core/datamining/pkg/validating"
)

const (
	dataMiningInternalPort = 3333
)

func main() {

	db := mem.NewDatabase()
	valid := validating.NewService()

	listService := listing.NewService(db)
	addService := adding.NewService(db, valid)

	reader, err := file.NewReader()
	if err != nil {
		log.Panic("Error reading keys...")
	}

	sharedRobotPrivateKey, err := reader.SharedRobotPrivateKey()
	if err != nil {
		log.Panic("Error reading keys...")
	}
	log.Print("DataMining Service starting...")

	//Starts Internal grpc server
	if err := startInternalServer(dataMiningInternalPort, listService, addService, sharedRobotPrivateKey); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(port int, listService listing.Service, addService adding.Service, sharedRobotPrivateKey []byte) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterInternalServer(grpcServer, internalrpc.NewInternalServerHandler(listService, addService, sharedRobotPrivateKey))
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
