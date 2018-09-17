package discovery

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/gossip"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/boostraping"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/file"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/mock"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/seeding"
)

func main() {
	pbKey, port, ver, p2pfactor := loadConfiguration()
	repo := mock.NewRepository()
	loc := mock.NewPeerLocalizer()

	if err := boostraping.NewService(repo, loc).Startup(pbKey, port, ver, p2pfactor); err != nil {
		log.Fatal(err)
	}

	seedReader := file.SeedReader{}
	if err := seeding.NewService(seedReader, repo).LoadSeeds(); err != nil {
		log.Fatal(err)
	}

	msg := mock.NewGossipMessenger()
	notif := mock.NewGossipNotifier()

	gossip.NewService(repo, msg, notif)

}

// func startGossip(peer discovery.Peer, repo discovery.PeerRepository) {
// 	msg := mock.NewGossipMessenger()
// 	notif := mock.NewGossipNotifier()
// 	gs := discovery.NewGossiper(repo, msg, notif)

// 	ticker := time.NewTicker(1 * time.Second)
// 	for range ticker.C {
// 		knownPeers, err := repo.ListPeers()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		round := discovery.NewGossipRound(peer, seeds, knownPeers)
// 		if err := gs.Gossip(round); err != nil {
// 			log.Print("Gossip failure: %s", err.Error())
// 		}
// 	}
// }

func loadConfiguration() ([]byte, int, string, int) {
	port := flag.Int("port", 3545, "Discovery port")
	p2pFactor := flag.Int("p2p-factor", 1, "P2P replication factor")
	pbKeyFile := flag.String("public-key-file", "id.pub", "Public key file")

	version, err := getVersion()
	if err != nil {
		log.Panic(err)
	}

	pbKey, err := getPublicKey(*pbKeyFile)
	if err != nil {
		log.Panic(err)
	}

	flag.Parse()

	return pbKey, *port, version, *p2pFactor
}

func getVersion() (string, error) {
	bytes, err := ioutil.ReadFile("version")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getPublicKey(file string) ([]byte, error) {

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
