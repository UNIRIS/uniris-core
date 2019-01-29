package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	memstorage "github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/system"
	memtransport "github.com/uniris/uniris-core/pkg/transport/mem"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
)

func startDiscovery(conf system.UnirisConfig) {
	log.Print("------------------------------")
	log.Print("DISCOVERY SERVICE STARTING...")
	log.Print("------------------------------")
	log.Printf("Port: %d", conf.Services.Discovery.Port)

	db := memstorage.NewDiscoveryDatabase()
	pnet := system.NewPeerNetworker()

	var pInfo discovery.PeerInformer
	if conf.Network.Type == "private" {
		pInfo = system.NewPeerInformer(true, conf.Network.Interface)
	} else {
		pInfo = system.NewPeerInformer(false, "")
	}

	msg := rpc.NewGossipRoundMessenger()
	notif := memtransport.NewGossipNotifier()
	discoverySrv := discovery.NewService(db, msg, notif, pnet, pInfo)

	go startDiscoveryServer(discoverySrv, conf.Services.Discovery.Port)

	peer, err := discoverySrv.StoreLocalPeer(conf.PublicKey, conf.Services.Discovery.Port, conf.Version)
	if err != nil {
		panic(err)
	}
	log.Print("Local peer stored")

	startGossip(peer, discoverySrv, conf)
}

func getSeeds(conf system.UnirisConfig) (seeds []discovery.Seed) {
	seedsConf := strings.Split(conf.Services.Discovery.Seeds, ";")
	for _, s := range seedsConf {
		seedItems := strings.Split(s, ":")
		ip := net.ParseIP(seedItems[0])
		port, _ := strconv.Atoi(seedItems[1])
		key := seedItems[2]
		seeds = append(seeds, discovery.Seed{
			PeerIdentity: discovery.NewPeerIdentity(ip, port, key),
		})
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

func startGossip(p discovery.Peer, discoverySrv discovery.Service, conf system.UnirisConfig) {
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
