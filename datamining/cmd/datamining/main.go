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
	biodAdding "github.com/uniris/uniris-core/datamining/pkg/biod/adding"
	biodListing "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	mocktransport "github.com/uniris/uniris-core/datamining/pkg/transport/mock"
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

	poolFinder := mocktransport.NewPoolFinder()
	aiClient := mocktransport.NewAIClient()

	signer := crypto.NewSigner()
	hasher := crypto.NewHasher()
	decrypter := crypto.NewDecrypter()

	rpcCrypto := rpc.NewCrypto(decrypter, signer, hasher)

	grpcClient := rpc.NewExternalClient(rpcCrypto, *config)
	poolRequester := rpc.NewPoolRequester(grpcClient, *config, rpcCrypto)

	biodAdder := biodAdding.NewService(db)
	biodLister := biodListing.NewService(db)
	lockSrv := lock.NewService(db)
	accountLister := accountListing.NewService(db)
	accountAdder := accountAdding.NewService(aiClient, db, accountLister, signer, hasher)

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction:  accountMining.NewKeychainMiner(signer, hasher, accountLister),
		mining.BiometricTransaction: accountMining.NewBiometricMiner(signer, hasher),
	}

	miningSrv := mining.NewService(
		aiClient,
		mocktransport.NewNotifier(),
		poolFinder,
		poolRequester,
		signer,
		biodLister,
		*config,
		txMiners,
	)

	log.Print("DataMining Service starting...")

	go func() {
		internalHandler := rpc.NewInternalServerHandler(biodAdder, poolRequester, aiClient, rpcCrypto, *config)

		//Starts Internal grpc server
		if err := startInternalServer(internalHandler, config.Datamining.InternalPort); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts Internal grpc server
	rpcServices := rpc.NewExternalServices(lockSrv, miningSrv, accountAdder, accountLister)
	externalHandler := rpc.NewExternalServerHandler(rpcServices, rpcCrypto, *config)
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
