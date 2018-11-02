package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/datamining/pkg/account/adding"
	"github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/locking"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	"github.com/uniris/uniris-core/datamining/pkg/transport/mock"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/externalrpc"
	internalrpc "github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
)

const (
	defaultConfFile = "../../../default-conf.yml"
)

func main() {

	config, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	db := mem.NewDatabase(config.SharedKeys.BiodPublicKey)

	poolFinder := mock.NewPoolFinder()
	addingSrv := adding.NewService(db)
	poolDispatcher := externalrpc.NewPoolDispatcher(config.Datamining)

	txLocker := mockstorage.NewTransactionLocker()

	listService := listing.NewService(db)
	lockSrv := locking.NewService(txLocker)

	masterMiningSrv := master.NewService(
		poolFinder,
		poolDispatcher,
		mock.NewNotifier(),
		crypto.NewSigner(),
		crypto.NewHasher(),
		listService,
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
	)

	slaveMiningSrv := slave.NewService(
		crypto.NewSigner(),
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
	)

	log.Print("DataMining Service starting...")

	go func() {

		//Starts Internal grpc server
		if err := startInternalServer(listService, masterMiningSrv, *config); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts Internal grpc server
	if err := startExternalServer(listService, addingSrv, slaveMiningSrv, lockSrv, *config); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(listService listing.Service, mineService master.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.InternalPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	handler := internalrpc.NewInternalServerHandler(listService, mineService,
		config.SharedKeys.RobotPrivateKey,
		config.Datamining.Errors)

	api.RegisterInternalServer(grpcServer, handler)
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", config.Datamining.InternalPort)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func startExternalServer(listService listing.Service, add adding.Service, mineService slave.Service, lockSrv locking.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.ExternalPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	handler := externalrpc.NewExternalServerHandler(listService, add, mineService, lockSrv,
		config.SharedKeys.RobotPublicKey,
		config.Datamining.Errors)

	api.RegisterExternalServer(grpcServer, handler)
	log.Printf("External grpc Server listening on 127.0.0.1:%d", config.Datamining.ExternalPort)
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
