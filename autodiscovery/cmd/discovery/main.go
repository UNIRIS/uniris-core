package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/http"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rabbitmq"
	"github.com/uniris/uniris-core/autodiscovery/pkg/transport/rpc"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/bootstraping"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/storage/mem"
	"google.golang.org/grpc"
)

const (
	seedsFile        = "../../configs/seeds.json"
	versionFile      = "../../configs/version"
	defaultPbKeyFile = "../../configs/id.pub"
)

func main() {

	//Loads peer's configuration
	network, pbKey, port, ver, seedsFile, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("PEER CONFIGURATION")
	log.Print("=================================")
	log.Printf("Network type: %s", network)
	log.Printf("Key: %s", pbKey)
	log.Printf("Port: %d", port)
	log.Printf("Version: %s", ver)

	//Initializes dependencies
	repo := new(mem.Repository)
	var np bootstraping.PeerNetworker
	if network == "public" {
		np = http.NewPeerNetworker()
	} else {
		np = system.NewPeerNetworker()
	}
	pos := http.NewPeerPositioner()
	monit := monitoring.NewService(repo, system.NewSystemWatcher(repo))
	notif := rabbitmq.NewNotifier()
	msg := rpc.NewMessenger()

	//Store the startup peer
	boot := bootstraping.NewService(repo, pos, np)
	startPeer, err := boot.Startup(pbKey, port, ver)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Endpoint: %s", startPeer.GetEndpoint())
	log.Print("=================================")

	//Starts server
	go func() {
		if err := startServer(port, repo, notif); err != nil {
			log.Fatal(err)
		}
	}()

	//Initializes the seeds
	if err := loadSeeds(seedsFile, boot); err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	//Starts gossiping
	log.Print("Start gossip...")
	g := gossip.NewService(repo, msg, notif, monit)
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if err := g.Spread(startPeer); err != nil {
			log.Printf("Gossip failure: %s", err.Error())
		}
		selfp, _ := repo.GetOwnedPeer()
		log.Printf("DEBUG: cpu: %s, freedisk: %b, status: %d, discoveredPeersNumber: %d", selfp.CPULoad(), selfp.FreeDiskSpace(), selfp.Status(), selfp.DiscoveredPeersNumber())
	}
}

func loadConfiguration() (string, []byte, int, string, string, error) {
	network := flag.String("network", "public", "Network type: public, private")
	port := flag.Int("port", 3545, "Discovery port")
	pbKeyFile := flag.String("key-file", defaultPbKeyFile, "Public key file")
	seedsFile := flag.String("seeds-file", seedsFile, "Seeds listing file")

	flag.Parse()

	pbKeyPath, err := filepath.Abs(*pbKeyFile)
	if err != nil {
		return "", nil, 0, "", "", err
	}

	pbKey, err := ioutil.ReadFile(pbKeyPath)
	if err != nil {
		return "", nil, 0, "", "", err
	}

	verPath, err := filepath.Abs(versionFile)
	if err != nil {
		return "", nil, 0, "", "", err
	}
	verBytes, err := ioutil.ReadFile(verPath)
	if err != nil {
		return "", nil, 0, "", "", err
	}
	version := string(verBytes)

	return *network, pbKey, *port, version, *seedsFile, nil
}

func startServer(port int, repo discovery.Repository, notif gossip.Notifier) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServer(grpcServer, rpc.NewHandler(repo, notif))
	log.Printf("Server listening on %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func loadSeeds(seedsFile string, boot bootstraping.Service) error {
	seedPath, err := filepath.Abs(seedsFile)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(seedPath)
	if err != nil {
		return err
	}

	seeds := make([]discovery.Seed, 0)
	if err := json.Unmarshal(bytes, &seeds); err != nil {
		return err
	}

	if err := boot.LoadSeeds(seeds); err != nil {
		return err
	}
	return nil
}
