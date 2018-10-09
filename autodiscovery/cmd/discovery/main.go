package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
	"github.com/uniris/uniris-core/autodiscovery/pkg/storage/redis"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/amqp"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rpc"
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

	//Initialize middlewares dependencies
	repo, err := redis.NewRepository(conf.Discovery.Redis)
	if err != nil {
		log.Fatal("Cannot connect with redis")
	}
	notif := amqp.NewNotifier(conf.Discovery.AMQP)

	//Initializes the infrastructure implementations
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
	msg := rpc.NewMessenger()

	//Setup services
	mon := monitoring.NewService(repo, system.NewPeerMonitor(), np, system.NewRobotWatcher())
	gos := gossip.NewService(repo, msg, notif, mon)
	boot := bootstraping.NewService(repo, pos, np)

	//Initializes the seeds
	seeds := make([]discovery.Seed, 0)
	for _, s := range conf.Discovery.Seeds {
		seeds = append(seeds, discovery.Seed{
			IP:        net.ParseIP(s.IP),
			Port:      s.Port,
			PublicKey: []byte(s.PublicKey),
		})
	}
	if err := boot.LoadSeeds(seeds); err != nil {
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
		if err := startServer(conf.Discovery.Port, repo, notif); err != nil {
			log.Fatal(err)
		}
	}()

	//Waiting server to start up
	time.Sleep(1 * time.Second)

	//We want to gossip every seconds
	ticker := time.NewTicker(1 * time.Second)

	//Starts the gossip
	log.Print("Gossip started...")
	res, err := gos.Start(startPeer, ticker)
	if err != nil {
		log.Fatal(err)
	}

	//When an unexpected error from the gossip is returned we crash
	for range res.Finish {
		err := <-res.Errors
		res.CloseChannels()
		log.Fatal(fmt.Errorf("Unexpected error %s", err.Error()))
	}
}

func loadConfiguration() (*system.UnirisConfig, error) {
	confFile := flag.String("conf-file", defaultConfFile, "Configuration file")
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

func startServer(port int, repo discovery.Repository, notif gossip.Notifier) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServer(grpcServer, rpc.NewServerHandler(repo, notif))
	log.Printf("Server listening on %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
