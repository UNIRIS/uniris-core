package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/shared"
	memstorage "github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/system"
	"github.com/uniris/uniris-core/pkg/transaction"
	memtransport "github.com/uniris/uniris-core/pkg/transport/mem"
	"github.com/uniris/uniris-core/pkg/transport/rest"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

const (
	defaultConfigurationFile = "./conf.yaml"
)

func main() {

	conf := UnirisConf{}

	app := cli.NewApp()
	app.Name = "uniris-miner"
	app.Usage = "UNIRIS miner"
	app.Version = "0.0.1"
	app.Flags = getCliFlags(&conf)

	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("conf"))
	app.Action = func(c *cli.Context) error {

		if c.String("private-key") == "" {
			fmt.Printf("Error: missing private key\n\n")
			return cli.ShowAppHelp(c)
		}

		if c.String("private-key") == "" {
			fmt.Printf("Error: missing public key\n\n")
			return cli.ShowAppHelp(c)
		}

		if c.String("discovery-seeds") == "" {
			fmt.Printf("Error: missing seeds\n\n")
			return cli.ShowAppHelp(c)
		}

		conf.version = app.Version

		fmt.Println("----------")
		fmt.Println("UNIRIS MINER")
		fmt.Println("----------")
		fmt.Printf("Version: %s\n", conf.version)
		fmt.Printf("Public key: %s\n", conf.publicKey)
		fmt.Printf("Network: %s\n", conf.networkType)
		fmt.Printf("Network interface: %s\n", conf.networkInterface)

		go startDiscovery(conf)
		go startDatamining(conf)
		startAPI(conf)

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCliFlags(conf *UnirisConf) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "conf",
			Usage:  "Configuration file",
			EnvVar: "UNIRIS_CONFIGURATION_FILE",
			Value:  defaultConfigurationFile,
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "private-key",
			Usage:       "Miner private key in hexadecimal",
			EnvVar:      "UNIRIS_PRIVATE_KEY",
			Destination: &conf.privateKey,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "public-key",
			Usage:       "Miner public key in hexadecimal",
			EnvVar:      "UNIRIS_PUBLIC_KEY",
			Destination: &conf.publicKey,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "network-type",
			EnvVar:      "UNIRIS_NETWORK_TYPE",
			Value:       "public",
			Usage:       "Type of the blockchain network (public or private) - Help to identify the IP address",
			Destination: &conf.networkType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "network-interface",
			EnvVar:      "UNIRIS_NETWORK_INTERFACE",
			Usage:       "Name of the network interface when type of network is private",
			Destination: &conf.networkInterface,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "discovery-port",
			EnvVar:      "UNIRIS_DISCOVERY_PORT",
			Value:       4000,
			Usage:       "Discovery service port",
			Destination: &conf.services.discovery.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-seeds",
			EnvVar:      "UNIRIS_DISCOVERY_SEEDS",
			Usage:       "List of the seeds peers to bootstrap the miner `IP:PORT:PUBLIC_KEY;IP:PORT:PUBLIC_KEY`",
			Destination: &conf.services.discovery.seeds,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-db-host",
			EnvVar:      "UNIRIS_DISCOVERY_DB_PORT",
			Value:       "localhost",
			Usage:       "Database instance hostname",
			Destination: &conf.services.discovery.db.host,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "discovery-db-port",
			Value:       6379,
			EnvVar:      "UNIRIS_DISCOVERY_DB_PORT",
			Usage:       "Redis instance port",
			Destination: &conf.services.discovery.db.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-db-password",
			EnvVar:      "UNIRIS_DISCOVERY_DB_PWD",
			Usage:       "Redis instance password",
			Destination: &conf.services.discovery.db.pwd,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-amqp-host",
			EnvVar:      "UNIRIS_DISCOVERY_AMQP_HOST",
			Value:       "localhost",
			Usage:       "AMQP instance hostname",
			Destination: &conf.services.discovery.amqp.host,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "discovery-amqp-port",
			EnvVar:      "UNIRIS_DISCOVERY_AMQP_PORT",
			Value:       5672,
			Usage:       "AMQP instance port",
			Destination: &conf.services.discovery.amqp.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-amqp-user",
			EnvVar:      "UNIRIS_DISCOVERY_AMQP_USER",
			Value:       "guest",
			Usage:       "AMQP instance user",
			Destination: &conf.services.discovery.amqp.user,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-amqp-password",
			EnvVar:      "UNIRIS_DISCOVERY_AMQP_PWD",
			Value:       "guest",
			Usage:       "AMQP instance password",
			Destination: &conf.services.discovery.amqp.password,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "datamining-port",
			EnvVar:      "UNIRIS_DATAMINING_PORT",
			Value:       5000,
			Usage:       "Datamining port",
			Destination: &conf.services.datamining.externalPort,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "datamining-internal-port",
			Value:       3009,
			Usage:       "Datamining internal port",
			Hidden:      true,
			Destination: &conf.services.datamining.internalPort,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "api-port",
			Value:       8080,
			Usage:       "API port",
			Destination: &conf.services.api.port,
		}),
	}
}

func startAPI(conf UnirisConf) {
	r := gin.Default()

	staticDir, _ := filepath.Abs("../../web/static")
	r.StaticFS("/static/", http.Dir(staticDir))

	rootPage, _ := filepath.Abs("../../web/index.html")
	r.StaticFile("/", rootPage)
	swaggerFile, _ := filepath.Abs("../../api/swagger-spec/swagger.yaml")
	r.StaticFile("/swagger.yaml", swaggerFile)

	apiRouter := r.Group("/api")
	{
		rest.NewAccountHandler(apiRouter, conf.services.datamining.internalPort, "")
		rest.NewTransactionHandler(apiRouter, conf.services.datamining.internalPort)
	}

	r.Run(fmt.Sprintf(":%d", conf.services.api.port))
}

func startDatamining(conf UnirisConf) {
	var pInfo discovery.PeerInformer
	if conf.networkType == "private" {
		pInfo = system.NewPeerInformer(true, conf.networkInterface)
	} else {
		pInfo = system.NewPeerInformer(false, "")
	}

	ip, err := pInfo.IP()
	if err != nil {
		panic(err)
	}

	txDb := memstorage.NewTransactionDatabase()
	lockDb := memstorage.NewLockDatabase()
	sharedDb := memstorage.NewSharedDatabase()

	poolReq := rpc.NewPoolRequester("", "")
	poolRetr := rpc.NewPoolRetriever("", "")

	sharedSrv := shared.NewService(sharedDb)
	poolFinderSrv := transaction.NewPoolFindingService(poolRetr)
	miningSrv := transaction.NewMiningService(poolReq, poolFinderSrv, sharedSrv, ip.String(), conf.publicKey, conf.privateKey)
	storeSrv := transaction.NewStorageService(txDb, miningSrv)
	lockSrv := transaction.NewLockService(lockDb)

	go startDataminingInternalServer(conf, poolFinderSrv, miningSrv)
	startDataminingExternalServer(conf, storeSrv, lockSrv, miningSrv)
}

func startDataminingInternalServer(conf UnirisConf, poolFinderSrv transaction.PoolFindingService, miningSrv transaction.MiningService) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conf.services.datamining.internalPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	intSrv := rpc.NewInternalServer(poolFinderSrv, miningSrv, "", "")
	api.RegisterInternalServiceServer(grpcServer, intSrv)
	log.Printf("Internal GRPC Server listening on %d", conf.services.datamining.internalPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startDataminingExternalServer(conf UnirisConf, storeSrv transaction.StorageService, lockSrv transaction.LockService, miningSrv transaction.MiningService) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conf.services.datamining.externalPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	txSrv := rpc.NewTransactionServer(storeSrv, lockSrv, miningSrv, "", "")
	api.RegisterTransactionServiceServer(grpcServer, txSrv)

	log.Printf("Transaction GRPC Server listening on %d", conf.services.datamining.externalPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startDiscovery(conf UnirisConf) {
	log.Print("------------------------------")
	log.Print("DISCOVERY SERVICE STARTING...")
	log.Print("------------------------------")
	log.Printf("Port: %d", conf.services.discovery.port)

	db := memstorage.NewDiscoveryDatabase()
	pnet := system.NewPeerNetworker()

	var pInfo discovery.PeerInformer
	if conf.networkType == "private" {
		pInfo = system.NewPeerInformer(true, conf.networkInterface)
	} else {
		pInfo = system.NewPeerInformer(false, "")
	}

	cli := rpc.NewDiscoveryClient()
	notif := memtransport.NewGossipNotifier()
	discoverySrv := discovery.NewService(db, cli, notif, pnet, pInfo)

	go startDiscoveryServer(discoverySrv, conf.services.discovery.port)

	peer, err := discoverySrv.StoreLocalPeer(conf.publicKey, conf.services.discovery.port, conf.version)
	if err != nil {
		panic(err)
	}
	log.Print("Local peer stored")

	startGossip(peer, discoverySrv, conf)
}

func getSeeds(conf UnirisConf) (seeds []discovery.PeerIdentity) {
	seedsConf := strings.Split(conf.services.discovery.seeds, ";")
	for _, s := range seedsConf {
		seedItems := strings.Split(s, ":")
		ip := net.ParseIP(seedItems[0])
		port, _ := strconv.Atoi(seedItems[1])
		key := seedItems[2]
		seeds = append(seeds, discovery.NewPeerIdentity(ip, port, key))
	}
	return
}

func startDiscoveryServer(discoverySrv discovery.Service, discoveryPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", discoveryPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServiceServer(grpcServer, rpc.NewDiscoveryServer(discoverySrv))
	log.Printf("Discovery GRPC server listening on %d", discoveryPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startGossip(p discovery.Peer, discoverySrv discovery.Service, conf UnirisConf) {
	timer := time.NewTicker(time.Second * 3)
	log.Print("Gossip running...")
	seeds := getSeeds(conf)
	abortChan, err := discoverySrv.Gossip(p, seeds, timer)
	if err != nil {
		panic(err)
	}

	for err := range abortChan {
		log.Fatalf("Gossip aborted - Error: %s", err.Error())
	}
}

type UnirisConf struct {
	networkType      string
	networkInterface string
	publicKey        string
	privateKey       string
	version          string
	services         struct {
		api struct {
			port int
		}
		discovery struct {
			port  int
			seeds string
			db    struct {
				host string
				port int
				pwd  string
			}
			amqp struct {
				host     string
				port     int
				user     string
				password string
			}
		}
		datamining struct {
			internalPort int
			externalPort int
		}
	}
}
