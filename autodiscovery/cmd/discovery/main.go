package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
	"github.com/uniris/uniris-core/autodiscovery/pkg/storage/redis"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rabbitmq"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rpc"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"google.golang.org/grpc"
)

const (
	defaultConfFile = "../../../default-conf.yml"
)

func main() {
	log.Print("Service starting...")

	conf, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("PEER CONFIGURATION")
	log.Print("=================================")
	log.Printf("Network type: %s", conf.Network)
	log.Printf("Key: %s", conf.PublicKey)
	log.Printf("Port: %d", conf.Discovery.Port)
	log.Printf("Version: %s", conf.Version)

	//Initializes dependencies
	repo, err := redis.NewRepository(conf.Discovery.Redis)
	if err != nil {
		log.Fatal("Cannot connect with redis")
	}
	var np monitoring.PeerNetworker
	if conf.Network == "public" {
		np = system.NewPublicNetworker()
	} else {
		if conf.NetworkInterface == "" {
			log.Fatal("Missing the network interface configuration when using the private network")
		}
		np = system.NewPrivateNetworker(conf.NetworkInterface)
	}
	pos := system.NewPeerPositioner()
	notif := rabbitmq.NewNotifier()
	msg := rpc.NewMessenger()
	mon := monitoring.NewService(repo, system.NewPeerMonitor(), np, system.NewRobotWatcher())
	gos := gossip.NewService(repo, msg, notif)
	boot := bootstraping.NewService(repo, pos, np)

	//Initializes the seeds
	if err := boot.LoadSeeds(conf.Discovery.Seeds); err != nil {
		log.Fatal(err)
	}

	//Stores the startup peer
	startPeer, err := boot.Startup([]byte(conf.PublicKey), conf.Discovery.Port, conf.Version)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Endpoint: %s", startPeer.Endpoint())
	log.Print("=================================")

	//Starts server
	go func() {
		if err := startServer(conf.Discovery.Port, repo, gos, mon, notif); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts gossiping
	time.Sleep(1 * time.Second)
	log.Print("Start gossip...")
	if err := gos.Spread(startPeer); err != nil {
		log.Print(err)
	}
}

func loadConfiguration() (*system.UnirisConfig, error) {
	confFile := flag.String("conf-file", defaultConfFile, "Configuration file")
	flag.Parse()

	var conf *system.UnirisConfig
	if os.Getenv("UNIRIS_VERSION") != "" {
		conf, err := system.BuildFromEnv()
		if err != nil {
			return nil, err
		}
		return conf, nil
	}
	conf, err := system.BuildFromFile(*confFile)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func startServer(port int, repo discovery.Repository, gos gossip.Service, mon monitoring.Service, notif gossip.Notifier) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServer(grpcServer, rpc.NewHandler(repo, gos, mon, notif))
	log.Printf("Server listening on %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
