package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/uniris/uniris-core/pkg/consensus"

	"github.com/uniris/uniris-core/pkg/crypto"
	logging "github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"
	"github.com/uniris/uniris-core/pkg/system"

	"github.com/uniris/uniris-core/pkg/transport/amqp"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	memstorage "github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/storage/redis"
	memtransport "github.com/uniris/uniris-core/pkg/transport/mem"
	"github.com/uniris/uniris-core/pkg/transport/rest"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

const (
	defaultConfigurationFile = "./conf.yaml"

	discoveryLogFile = "Discovery.log"
	miningLogFile    = "Mining.log"

	discoveryAppID = "Discovery"
	miningAppID    = "Mining"
)

var appIDLogFile = map[string]string{
	discoveryAppID: discoveryLogFile,
	miningAppID:    miningLogFile,
}

func main() {

	rand.Seed(time.Now().Unix())

	conf := unirisConf{}

	app := cli.NewApp()
	app.Name = "uniris"
	app.Usage = "UNIRIS node"
	app.Version = "0.0.1"
	app.Flags = getCliFlags(&conf)

	app.Before = altsrc.InitInputSourceWithContext(app.Flags, func(c *cli.Context) (altsrc.InputSourceContext, error) {
		confFile := c.String("conf")
		context, err := altsrc.NewYamlSourceFromFile(confFile)
		if err != nil {
			fmt.Println("Load configuration by environment variables")
			return &altsrc.MapInputSource{}, nil
		}
		fmt.Println("Load configuration by file")
		return context, nil
	})
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
		fmt.Println("UNIRIS NODE")
		fmt.Println("----------")
		fmt.Printf("Version: %s\n", conf.version)
		fmt.Printf("Public key: %s\n", conf.publicKey)
		fmt.Printf("Network: %s\n", conf.networkType)
		fmt.Printf("Network interface: %s\n", conf.networkInterface)

		sharedDB := memstorage.NewSharedDatabase()

		var nodeDB consensus.NodeReadWriter
		if conf.nodeDatabase.dbType == "redis" {
			nodeRedisDB, err := redis.NewNodeDatabase(conf.nodeDatabase.host, conf.nodeDatabase.port, conf.nodeDatabase.pwd)
			if err != nil {
				panic(err)
			}
			nodeDB = nodeRedisDB
		} else {
			nodeDB = &memstorage.NodeDatabase{}
		}

		ip := getIP(conf.networkType, conf.networkInterface)
		loggerDiscovery := createLogger(discoveryAppID, conf.logging.logType, conf.logging.LogDir, conf.logging.logLevel, ip)
		loggerMining := createLogger(miningAppID, conf.logging.logType, conf.logging.LogDir, conf.logging.logLevel, ip)

		go startGRPCServer(conf, sharedDB, nodeDB, ip, loggerDiscovery, loggerMining)
		startHTTPServer(conf, sharedDB, nodeDB, loggerMining)

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getCliFlags(conf *unirisConf) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "conf",
			Usage:  "Configuration file",
			EnvVar: "UNIRIS_CONFIGURATION_FILE",
			Value:  defaultConfigurationFile,
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "private-key",
			Usage:       " private key in hexadecimal",
			EnvVar:      "UNIRIS_PRIVATE_KEY",
			Destination: &conf.privateKey,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "public-key",
			Usage:       " public key in hexadecimal",
			EnvVar:      "UNIRIS_PUBLIC_KEY",
			Destination: &conf.publicKey,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "network-type",
			EnvVar:      "UNIRIS_NETWORK_TYPE",
			Value:       "public",
			Usage:       "Type of the network (public or private)",
			Destination: &conf.networkType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "network-interface",
			EnvVar:      "UNIRIS_NETWORK_INTERFACE",
			Usage:       "Name of the network interface when type of network is private",
			Destination: &conf.networkInterface,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-seeds",
			EnvVar:      "UNIRIS_DISCOVERY_SEEDS",
			Usage:       "List of the seeds peers to bootstrap the node `IP:PORT:PUBLIC_KEY;IP:PORT:PUBLIC_KEY`",
			Destination: &conf.discoverySeeds,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-db-type",
			EnvVar:      "UNIRIS_DISCOVERY_DB_TYPE",
			Value:       "mem",
			Usage:       "Discovery database instance type (mem or redis)",
			Destination: &conf.discoveryDatabase.dbType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-db-host",
			EnvVar:      "UNIRIS_DISCOVERY_DB_HOST",
			Value:       "localhost",
			Usage:       "Discovery database instance hostname",
			Destination: &conf.discoveryDatabase.host,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "discovery-db-port",
			Value:       6379,
			EnvVar:      "UNIRIS_DISCOVERY_DB_PORT",
			Usage:       "Discovery database instance port",
			Destination: &conf.discoveryDatabase.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "discovery-db-password",
			EnvVar:      "UNIRIS_DISCOVERY_DB_PWD",
			Usage:       "Discovery database instance password",
			Destination: &conf.discoveryDatabase.pwd,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "node-db-type",
			EnvVar:      "UNIRIS_NODE_DB_TYPE",
			Value:       "mem",
			Usage:       "Node database instance type (mem or redis)",
			Destination: &conf.nodeDatabase.dbType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "node-db-host",
			EnvVar:      "UNIRIS_NODE_DB_HOST",
			Value:       "localhost",
			Usage:       "Node database instance hostname",
			Destination: &conf.nodeDatabase.host,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "node-db-port",
			Value:       6379,
			EnvVar:      "UNIRIS_NODE_DB_PORT",
			Usage:       "Node database instance port",
			Destination: &conf.nodeDatabase.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "node-db-password",
			EnvVar:      "UNIRIS_NODE_DB_PWD",
			Usage:       "Node database instance password",
			Destination: &conf.nodeDatabase.pwd,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "bus-type",
			EnvVar:      "UNIRIS_BUS_TYPE",
			Value:       "mem",
			Usage:       "Bus instance type (mem or amqp)",
			Destination: &conf.bus.busType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "bus-host",
			EnvVar:      "UNIRIS_BUS_HOST",
			Value:       "localhost",
			Usage:       "BUS instance hostname",
			Destination: &conf.bus.host,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "bus-port",
			EnvVar:      "UNIRIS_BUS_PORT",
			Value:       5672,
			Usage:       "Bus instance port",
			Destination: &conf.bus.port,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "bus-user",
			EnvVar:      "UNIRIS_BUS_USER",
			Value:       "guest",
			Usage:       "Bus instance user",
			Destination: &conf.bus.user,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "bus-password",
			EnvVar:      "UNIRIS_BUS_PWD",
			Value:       "guest",
			Usage:       "Bus instance password",
			Destination: &conf.bus.password,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "grpc-port",
			EnvVar:      "UNIRIS_GRPC_PORT",
			Value:       5000,
			Usage:       "GRPC server port used to communicate with other nodes",
			Destination: &conf.grpcPort,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "log-type",
			EnvVar:      "UNIRIS_LOG_TYPE",
			Value:       "file",
			Usage:       "Logging type (stdout/file)",
			Destination: &conf.logging.logType,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "log-level",
			EnvVar:      "UNIRIS_LOG_LEVEL",
			Value:       "info",
			Usage:       "Logging level (info/error/debug)",
			Destination: &conf.logging.logLevel,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "log-dir",
			EnvVar:      "UNIRIS_LOG_DIR",
			Value:       "/var/log/uniris",
			Usage:       "Dir to store log files",
			Destination: &conf.logging.LogDir,
		}),
	}
}

func startHTTPServer(conf unirisConf, sharedKeyReader shared.KeyReader, nodeReader consensus.NodeReader, l logging.Logger) {
	pubB, err := hex.DecodeString(conf.publicKey)
	if err != nil {
		panic(err)
	}
	publicKey, err := crypto.ParsePublicKey(pubB)
	if err != nil {
		panic(err)
	}

	pvB, err := hex.DecodeString(conf.privateKey)
	if err != nil {
		panic(err)
	}
	privateKey, err := crypto.ParsePrivateKey(pvB)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	staticDir, _ := filepath.Abs("../../web/static")
	r.StaticFS("/static/", http.Dir(staticDir))

	rootPage, _ := filepath.Abs("../../web/index.html")
	r.StaticFile("/", rootPage)
	swaggerFile, _ := filepath.Abs("../../api/swagger-spec/swagger.yaml")
	r.StaticFile("/swagger.yaml", swaggerFile)

	r.GET("/api/account/:idHash", rest.GetAccountHandler(sharedKeyReader, nodeReader, l))
	r.POST("/api/account", rest.CreateAccountHandler(sharedKeyReader, nodeReader, publicKey, privateKey, l))
	r.GET("/api/transaction/:txReceipt/status", rest.GetTransactionStatusHandler(sharedKeyReader, nodeReader, l))
	r.GET("/api/sharedkeys", rest.GetSharedKeysHandler(sharedKeyReader, l))

	r.Run(":80")
}

func startGRPCServer(conf unirisConf, sharedKeyRW shared.KeyReadWriter, nodeRW consensus.NodeReadWriter, ip net.IP, ld logging.Logger, lm logging.Logger) {
	pubB, err := hex.DecodeString(conf.publicKey)
	if err != nil {
		panic(err)
	}
	publicKey, err := crypto.ParsePublicKey(pubB)
	if err != nil {
		panic(err)
	}

	pvB, err := hex.DecodeString(conf.privateKey)
	if err != nil {
		panic(err)
	}
	privateKey, err := crypto.ParsePrivateKey(pvB)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	poolR := rpc.NewPoolRequester(sharedKeyRW, lm)
	chainDB := memstorage.NewchainDatabase()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, sharedKeyRW, nodeRW, poolR, publicKey, privateKey, lm))

	var discoveryDB discovery.Database
	if conf.discoveryDatabase.dbType == "redis" {
		redisDB, err := redis.NewDiscoveryDatabase(conf.discoveryDatabase.host, conf.discoveryDatabase.port, conf.discoveryDatabase.pwd)
		if err != nil {
			panic(err)
		}
		discoveryDB = redisDB
	} else {
		discoveryDB = memstorage.NewDiscoveryDatabase()
	}

	var notif discovery.Notifier
	if conf.bus.busType == "amqp" {
		notif = amqp.NewDiscoveryNotifier(conf.bus.host, conf.bus.user, conf.bus.password, conf.bus.port, ld)
		go func() {
			if err := amqp.ConsumeDiscoveryNotifications(conf.bus.host, conf.bus.user, conf.bus.password, conf.bus.port, nodeRW, sharedKeyRW, ld); err != nil {
				panic(err)
			}
		}()
	} else {
		notif = &memtransport.DiscoveryNotifier{
			Logger: ld,
		}
	}

	api.RegisterDiscoveryServiceServer(grpcServer, rpc.NewDiscoveryServer(discoveryDB, notif, ld, publicKey, privateKey, sharedKeyRW))

	go startDiscovery(conf, discoveryDB, notif, ld, publicKey, privateKey)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.grpcPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("GRPC server listening on %d\n", conf.grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startDiscovery(conf unirisConf, db discovery.Database, notif discovery.Notifier, logger logging.Logger, publicKey crypto.PublicKey, privateKey crypto.PrivateKey) {

	pub, err := hex.DecodeString(conf.publicKey)
	if err != nil {
		panic(err)
	}
	pubk, err := crypto.ParsePublicKey(pub)
	if err != nil {
		panic(err)
	}

	netCheck := system.NewNetworkChecker(conf.grpcPort)

	var systemReader discovery.SystemReader
	if conf.networkType == "private" {
		systemReader = system.NewReader(true, conf.networkInterface)
	} else {
		systemReader = system.NewReader(false, "")
	}

	roundMessenger := rpc.NewGossipRoundMessenger(logger, publicKey, privateKey)

	ip, err := systemReader.IP()
	if err != nil {
		panic(err)
	}

	lon, lat, err := systemReader.GeoPosition()
	if err != nil {
		logger.Error(discovery.ErrGeoPosition.Error())
	}
	selfPeer := discovery.NewSelfPeer(pubk, ip, conf.grpcPort, conf.version, lon, lat)

	seeds := getSeeds(conf)

	timer := time.NewTicker(time.Second * 3)
	logger.Info("Gossip running...")

	//Store local peer
	if err := db.WriteDiscoveredPeer(selfPeer); err != nil {
		panic(err)
	}

	for range timer.C {
		go func() {
			if err := discovery.Gossip(selfPeer, seeds, db, netCheck, systemReader, roundMessenger, notif, logger); err != nil {
				timer.Stop()
				panic(err)
			}
		}()
	}

}

func getSeeds(conf unirisConf) (seeds []discovery.PeerIdentity) {
	seedsConf := strings.Split(conf.discoverySeeds, ";")
	for _, s := range seedsConf {
		if s != "" {
			seedProps := strings.Split(s, ":")
			if len(seedProps) == 0 {
				panic("invalid seed")
			}
			ip := net.ParseIP(seedProps[0])
			port, _ := strconv.Atoi(seedProps[1])
			key := seedProps[2]
			pub, err := hex.DecodeString(key)
			if err != nil {
				panic(err)
			}
			pubk, err := crypto.ParsePublicKey(pub)
			if err != nil {
				panic(err)
			}
			seeds = append(seeds, discovery.NewPeerIdentity(ip, port, pubk))
		}
	}
	return
}

type unirisConf struct {
	networkType      string
	networkInterface string
	publicKey        string
	privateKey       string
	version          string
	sharedEmKey      struct { //TODO: to remove once the feature is implemented
		encryptedPrivateKey string
		publicKey           string
	}
	sharedKey struct { //TODO: to remove once the feature is implemented
		privateKey string
		publicKey  string
	}
	grpcPort int
	bus      struct {
		busType  string
		host     string
		port     int
		user     string
		password string
	}
	discoverySeeds    string
	discoveryDatabase struct {
		dbType string
		host   string
		port   int
		pwd    string
	}
	nodeDatabase struct {
		dbType string
		host   string
		port   int
		pwd    string
	}
	logging struct {
		logType  string
		logLevel string
		LogDir   string
	}
}

func createLogger(appID string, ty string, dir string, level string, ip net.IP) logging.Logger {

	//check if appID is in the list
	if _, appExist := appIDLogFile[appID]; !appExist {
		log.Fatal("[Fatal] Unkown application ID")
	}

	//check if log type has the good value
	if ty != "stdout" && ty != "file" {
		log.Fatal("[Fatal] log-type value should be (stdout|file)")
	}

	//check if log-level has the good value
	if level != "info" && level != "error" && level != "debug" {
		log.Fatal("[Fatal] log-level value should be (info|error|debug)")
	}

	//map Loglevel
	var ll logging.Level
	if level == "error" {
		ll = logging.ErrorLogLevel
	} else if level == "info" {
		ll = logging.InfoLogLevel
	} else {
		ll = logging.DebugLogLevel
	}

	if ty == "stdout" {
		stdLog := log.New(os.Stdout, "", 0)
		return logging.NewLogger("stdout", stdLog, appID, ip, ll)
	}

	if ty == "file" {

		src, err := os.Stat(dir)

		//check if logdir exist or not
		if os.IsNotExist(err) {
			log.Println("[Error] log-dir" + dir + "does not exist, please create the adequate directory")
			stdLog := log.New(os.Stdout, "", 0)
			return logging.NewLogger("stdout", stdLog, appID, ip, ll)
		}

		//check if logdir is not a file
		if src.Mode().IsRegular() {
			log.Println("[Erro] log-dir" + dir + "is a file, please create the adequate directory")
			stdLog := log.New(os.Stdout, "", 0)
			return logging.NewLogger("stdout", stdLog, appID, ip, ll)
		}
	}

	f, err := os.OpenFile(dir+"/"+appIDLogFile[appID], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		stdLog := log.New(os.Stdout, "", 0)
		return logging.NewLogger("stdout", stdLog, appID, ip, ll)
	}

	fileLog := log.New(f, "", 0)
	return logging.NewLogger("file", fileLog, appID, ip, ll)
}

func getIP(networkType string, networkInterface string) net.IP {

	var sysReader discovery.SystemReader

	if networkType == "private" {
		sysReader = system.NewReader(true, networkInterface)

	} else {
		sysReader = system.NewReader(false, "")
	}

	ip, err := sysReader.IP()
	if err != nil {
		log.Fatal("[Fatal] Cannot get the hostID")
	}

	return ip
}
