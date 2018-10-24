package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/leading"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	listing "github.com/uniris/uniris-core/datamining/pkg/listing"
	mem "github.com/uniris/uniris-core/datamining/pkg/storage/mem"
	"github.com/uniris/uniris-core/datamining/pkg/transport/mock"
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

	db := mem.NewDatabase(config.SharedKeys.BiodPublicKey)

	validSrv := validating.NewService(
		crypto.NewSigner(),
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
	)

	poolFinder := mock.NewPoolFinder()
	addingSrv := adding.NewService(db)
	poolDispatcher := mock.NewPoolDispatcher(addingSrv, validSrv)

	leadService := leading.NewService(
		poolFinder,
		poolDispatcher,
		mock.NewNotifier(),
		crypto.NewSigner(),
		db,
		config.SharedKeys.RobotPublicKey,
		config.SharedKeys.RobotPrivateKey,
	)

	listService := listing.NewService(db)

	log.Print("DataMining Service starting...")

	//Starts Internal grpc server
	if err := startInternalServer(listService, leadService, *config); err != nil {
		log.Fatal(err)
	}

}

func startInternalServer(listService listing.Service, leadService leading.Service, config system.UnirisConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Datamining.InternalPort))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	stores := make(chan string)
	handler := internalrpc.NewInternalServerHandler(listService, leadService,
		config.SharedKeys.RobotPrivateKey,
		config.Datamining.Errors,
		stores)

	go func() {
		for tx := range stores {
			log.Printf("Transaction %s stored on the database", tx)
		}
	}()

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
