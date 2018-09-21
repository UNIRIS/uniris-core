package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
	"github.com/uniris/uniris-core/autodiscovery/pkg/storage/redis"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/http"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rabbitmq"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rpc"
	yaml "gopkg.in/yaml.v2"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"google.golang.org/grpc"
)

const (
	defaultConfFile = "../../../default-conf.yml"
)

func main() {
	if err := runApp(); err != nil {
		os.Exit(-1)
	}
}

func runApp() error {
	log.Print("Service starting...")

	log.Print("Configuration loading...")
	//Loads peer's configuration
	conf, err := loadConfiguration()
	if err != nil {
		return err
	}

	log.Print("PEER CONFIGURATION")
	log.Print("=================================")
	log.Printf("Network type: %s", conf.Network)
	log.Printf("Key: %s", conf.PublicKey)
	log.Printf("Port: %d", conf.Discovery.Port)
	log.Printf("Version: %s", conf.Version)
	log.Printf("P2P Factor: %d", conf.Discovery.P2PFactor)

	//Initializes dependencies
	repo, err := redis.NewRepository(conf.Discovery.Redis.Host, conf.Discovery.Redis.Port, conf.Discovery.Redis.Pwd)
	if err != nil {
		log.Print("Cannot connect with redis")
		return err
	}
	var np bootstraping.PeerNetworker
	if conf.Network == "public" {
		np = http.NewPeerNetworker()
	} else {
		if conf.NetworkInterface == "" {
			return errors.New("Missing the network interface configuration when using the private network")
		}
		np = system.NewPeerNetworker(conf.NetworkInterface)
	}
	pos := http.NewPeerPositioner()
	notif := rabbitmq.NewNotifier()
	msg := rpc.NewMessenger()
	mon := monitoring.NewService(repo, system.NewPeerMonitor())
	gos := gossip.NewService(repo, msg, notif, mon)
	boot := bootstraping.NewService(repo, pos, np)

	//Initializes the seeds
	if err := boot.LoadSeeds(conf.Discovery.Seeds); err != nil {
		return err
	}

	//Stores the startup peer
	startPeer, err := boot.Startup([]byte(conf.PublicKey), conf.Discovery.Port, conf.Discovery.P2PFactor, conf.Version)
	if err != nil {
		return err
	}

	log.Printf("Endpoint: %s", startPeer.GetEndpoint())
	log.Print("=================================")

	//Starts server
	go func() {
		if err := startServer(conf.Discovery.Port, repo, gos, mon); err != nil {
			log.Fatal(err)
		}
	}()

	//Starts gossiping
	time.Sleep(1 * time.Second)
	log.Print("Start gossip...")
	if err := gos.ScheduleGossip(startPeer); err != nil {
		log.Print(err)
	}
	return nil
}

type Conf struct {
	Network          string `yaml:"network"`
	NetworkInterface string `yaml:"networkInterface"`
	PublicKey        string `yaml:"publicKey"`
	Version          string `yaml:"version"`
	Discovery        ConfDiscovery
}

type ConfDiscovery struct {
	Port      int              `yaml:"port"`
	P2PFactor int              `yaml:"p2pFactor"`
	Seeds     []discovery.Seed `yaml:"seeds"`
	Redis     ConfRedis
}

type ConfRedis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Pwd  string `yaml:"pwd"`
}

func loadConfiguration() (*Conf, error) {

	ver := os.Getenv("UNIRIS_VERSION")
	pbKey := os.Getenv("UNIRIS_PUBLICKEY")
	network := os.Getenv("UNIRIS_NETWORK")
	netiface := os.Getenv("UNIRIS_NETWORK_INTERFACE")
	port := os.Getenv("UNIRIS_DISCOVERY_PORT")
	p2pFactor := os.Getenv("UNIRIS_DISCOVERY_P2PFACTOR")
	seeds := os.Getenv("UNIRIS_DISCOVERY_SEEDS")
	redisHost := os.Getenv("UNIRIS_DISCOVERY_REDIS_HOST")
	redisPort := os.Getenv("UNIRIS_DISCOVERY_REDIS_PORT")
	redisPwd := os.Getenv("UNIRIS_DISCOVERY_REDIS_PWD")

	//LOAD BY ENV VARIABLE
	if ver != "" {
		_seeds := make([]discovery.Seed, 0)
		ss := strings.Split(seeds, ",")
		for _, s := range ss {
			addr := strings.Split(s, ":")
			sPort, _ := strconv.Atoi(addr[1])

			ips, err := net.LookupIP(addr[0])
			if err != nil {
				return nil, err
			}

			_seeds = append(_seeds, discovery.Seed{
				IP:   ips[0],
				Port: sPort,
			})
		}

		_port, _ := strconv.Atoi(port)
		_p2pFactor, _ := strconv.Atoi(p2pFactor)
		_redisPort, _ := strconv.Atoi(redisPort)

		return &Conf{
			Version:          ver,
			PublicKey:        pbKey,
			Network:          network,
			NetworkInterface: netiface,
			Discovery: ConfDiscovery{
				Port:      _port,
				P2PFactor: _p2pFactor,
				Seeds:     _seeds,
				Redis: ConfRedis{
					Host: redisHost,
					Port: _redisPort,
					Pwd:  redisPwd,
				},
			},
		}, nil
	}

	//LOAD BY CONFIGURATION FILE
	confFile := flag.String("conf-file", defaultConfFile, "Configuration file")
	flag.Parse()

	confFilePath, err := filepath.Abs(*confFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return nil, err
	}

	var c Conf
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func startServer(port int, repo discovery.Repository, gos gossip.Service, mon monitoring.Service) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServer(grpcServer, rpc.NewHandler(repo, gos, mon))
	log.Printf("Server listening on %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
