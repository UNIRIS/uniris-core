package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	accountAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accountListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
	biodlisting "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

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

	db := mem.NewDatabase()

	poolFinder := mock.NewPoolFinder()
	poolRequester := externalrpc.NewPoolRequester(config.Datamining)

	biodLister := biodlisting.NewService(db)
	lockSrv := lock.NewService(db)
	signer := crypto.NewSigner()
	hasher := crypto.NewHasher()

	accountLister := accountListing.NewService(db)
	accountAdder := accountAdding.NewService(db)

	checks := map[mining.TransactionType]mining.Checker{
		mining.CreateKeychainTransaction: accountMining.NewKeychainChecker(signer, hasher),
		mining.CreateBioTransaction:      accountMining.NewBiometricChecker(signer, hasher),
	}

	miningSrv := mining.NewService(
		mock.NewNotifier(),
		poolFinder,
		poolRequester,
		signer,
		biodLister,
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
		checks,
	)

	log.Print("DataMining Service starting...")

	go func() {
		//Starts Internal grpc server
		if err := startInternalServer(accountLister, miningSrv, *config); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts Internal grpc server
	if err := startExternalServer(lockSrv, miningSrv, accountAdder, *config); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(accountLister accountListing.Service, mineService mining.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.InternalPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	handler := internalrpc.NewInternalServerHandler(accountLister, mineService,
		config.SharedKeys.RobotPrivateKey,
		config.Datamining)

	api.RegisterInternalServer(grpcServer, handler)
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", config.Datamining.InternalPort)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func startExternalServer(lockSrv lock.Service, mineSrv mining.Service, accountAdder accountAdding.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.ExternalPort))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	handler := externalrpc.NewExternalServerHandler(lockSrv, mineSrv, accountAdder,
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
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
