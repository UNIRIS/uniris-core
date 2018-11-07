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
	"github.com/uniris/uniris-core/datamining/pkg/mock"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
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
	decrypter := crypto.NewDecrypter()
	aiClient := mock.NewAIClient()

	accountLister := accountListing.NewService(db)
	accountAdder := accountAdding.NewService(db)

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction:  accountMining.NewKeychainMiner(signer, hasher, accountLister),
		mining.BiometricTransaction: accountMining.NewBiometricMiner(signer, hasher),
	}

	miningSrv := mining.NewService(
		mock.NewNotifier(),
		poolFinder,
		poolRequester,
		signer,
		biodLister,
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
		txMiners,
	)

	log.Print("DataMining Service starting...")

	go func() {
		internalHandler := internalrpc.NewInternalServerHandler(poolRequester, aiClient, hasher, decrypter, *config)

		//Starts Internal grpc server
		if err := startInternalServer(internalHandler, config.Datamining.InternalPort); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts Internal grpc server
	externalHandler := externalrpc.NewExternalServerHandler(lockSrv, miningSrv, accountAdder, accountLister, decrypter, signer, *config)
	if err := startExternalServer(externalHandler, config.Datamining.ExternalPort); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(handler api.InternalServer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	api.RegisterInternalServer(grpcServer, handler)
	log.Printf("Internal grpc Server listening on 127.0.0.1:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func startExternalServer(handler api.ExternalServer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	api.RegisterExternalServer(grpcServer, handler)
	log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
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
