package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/gossip"
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

	var pInfo gossip.PeerInformer
	if conf.Network.Type == "private" {
		pInfo = system.NewPeerInformer(true, conf.Network.Interface)
	} else {
		pInfo = system.NewPeerInformer(false, "")
	}

	msg := rpc.NewGossipRoundMessenger()
	notif := memtransport.NewGossipNotifier()
	gossipSrv := gossip.NewService(db, msg, notif, pnet, pInfo)

	go startDiscoveryServer(gossipSrv, conf.Services.Discovery.Port)

	peer, err := gossipSrv.StoreLocalPeer(conf.PublicKey, conf.Services.Discovery.Port, conf.Version)
	if err != nil {
		panic(err)
	}
	log.Print("Local peer stored")

	startGossip(peer, gossipSrv, conf)
}

func getSeeds(conf system.UnirisConfig) (seeds []uniris.Seed) {
	seedsConf := strings.Split(conf.Services.Discovery.Seeds, ";")
	for _, s := range seedsConf {
		seedItems := strings.Split(s, ":")
		ip := net.ParseIP(seedItems[0])
		port, _ := strconv.Atoi(seedItems[1])
		key := seedItems[2]
		seeds = append(seeds, uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(ip, port, key),
		})
	}
	return
}

func startDiscoveryServer(gossipSrv gossip.Service, discoveryPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", discoveryPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServiceServer(grpcServer, rpc.NewDiscoveryServer(gossipSrv))
	log.Printf("Discovery GRPC server listening on %d", discoveryPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startGossip(p uniris.Peer, gossipSrv gossip.Service, conf system.UnirisConfig) {
	timer := time.NewTicker(time.Second * 3)
	log.Print("Gossip running...")
	seeds := getSeeds(conf)
	abortChan, err := gossipSrv.Run(p, seeds, timer)
	if err != nil {
		panic(err)
	}

	for err := range abortChan {
		log.Fatalf("Gossip aborted - Error: %s", err.Error())
	}
}
