package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	adding "github.com/uniris/uniris-core/datamining/pkg/adding"
	listing "github.com/uniris/uniris-core/datamining/pkg/listing"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/externalrpc"
	internalrpc "github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
	validating "github.com/uniris/uniris-core/datamining/pkg/validating"
)

const (
	defaultConfFile = "../../../default-conf.yml"
)

func main() {

	config, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	db := mem.NewDatabase()
	valid := validating.NewService(crypto.NewSigner(), externalrpc.NewValidatorRequest())

	listService := listing.NewService(db)
	addService := adding.NewService(db, valid)

	log.Print("DataMining Service starting...")

	//Starts Internal grpc server
	if err := startInternalServer(listService, addService, *config); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(listService listing.Service, addService adding.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.InternalPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	handler := internalrpc.NewInternalServerHandler(listService, addService,
		config.SharedKeys.RobotPrivateKey,
		config.Datamining.Errors)
	api.RegisterInternalServer(grpcServer, handler)
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", config.Datamining.InternalPort)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func loadConfiguration() (*system.UnirisConfig, error) {
	confFile := flag.String("config", defaultConfFile, "Configuration file")
	flag.Parse()

	confFilePath, err := filepath.Abs(*confFile)
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		conf, err := system.BuildFromEnv()
		if err != nil {
			return nil, err
		}
		return conf, nil
	}

	conf, err := system.BuildFromFile(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
