package main

import (
	"fmt"
	"log"
	"net"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	adding "github.com/uniris/uniris-core/datamining/pkg/adding"
	file "github.com/uniris/uniris-core/datamining/pkg/file"
	listing "github.com/uniris/uniris-core/datamining/pkg/listing"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/externalrpc"
	internalrpc "github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
	validating "github.com/uniris/uniris-core/datamining/pkg/validating"
)

const (
	dataMiningInternalPort = 3333
)

func main() {

	db := mem.NewDatabase()
	valid := validating.NewService(crypto.NewSigner(), externalrpc.NewValidatorRequest())

	listService := listing.NewService(db)
	addService := adding.NewService(db, valid)

	reader, err := file.NewReader()
	if err != nil {
		log.Fatalf("Error opening key file reader: %s", reader)
	}

	sharedRobotPrivateKey, err := reader.SharedRobotPrivateKey()
	if err != nil {
		log.Fatalf("Error reading shared private key:  %s", err.Error())
	}

	sharedRobotPublicKey, err := reader.SharedRobotPublicKey()
	if err != nil {
		log.Fatalf("Error reading shared public key:  %s", err.Error())
	}

	log.Print("DataMining Service starting...")

	//Starts Internal grpc server
	if err := startInternalServer(dataMiningInternalPort, listService, addService, sharedRobotPublicKey, sharedRobotPrivateKey); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(port int, listService listing.Service, addService adding.Service, sharedRobotPublicKey, sharedRobotPrivateKey []byte) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterInternalServer(grpcServer, internalrpc.NewInternalServerHandler(listService, addService, sharedRobotPublicKey, sharedRobotPrivateKey))
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
