// package main

// import (
// 	"log"
// 	"time"
// )

// func main() {

// 	log.Print("Autodiscovery starting...")

// 	// repo := new(adapters.InMemoryPeerRepository)

// 	time.Sleep(2 * time.Second)
// 	ticker := time.NewTicker(1 * time.Second)
// 	for range ticker.C {
// 		// rnd, err := ucBuilder.StartGossipRound()
// 		// if err != nil {
// 		// 	log.Print(err)
// 		// 	return
// 		// }
// 		// newPeers, err := ucBuilder.StartGossip(rnd)
// 		// if err != nil {
// 		// 	log.Print(err)
// 		// 	return
// 		// }
// 		// err := ucBuilder.FinishGossipRound(newPeers)
// 		// if err != nil {
// 		// 	log.Print(err)
// 		// 	return
// 		// }

// 	}

// }

// // 	httpClient := new(infrastructure.HTTPClient)
// // 	fileReader := new(infrastructure.FileReader)
// // 	repo := new(adapters.InMemoryPeerRepository)
// // 	metric := new(adapters.MetricReader)
// // 	conf := adapters.NewConfigurationReader(*fileReader, *httpClient)

// // 	if err := usecases.StartPeer(repo, conf); err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	go func() {
// // 		if err := startServer(repo, conf, metric); err != nil {
// // 			log.Fatal(err)
// // 		}
// // 	}()

// // 	time.Sleep(2 * time.Second)

// // 	broker := new(adapters.GrpcGossipBroker)
// // 	notifier := new(adapters.InMemoryDiscoveryNotifier)

// // 	log.Print("Seed loading...")
// // 	if err := usecases.LoadSeeds(repo, conf); err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	log.Print("Gossip starting...")

// // 	ticker := time.NewTicker(1 * time.Second)
// // 	for range ticker.C {
// // 		go func() {
// // 			if err := usecases.StartGossipRound(repo, broker, notifier); err != nil {
// // 				log.Printf("Gossip failure %s", err.Error())
// // 			}
// // 		}()
// // 	}
// // }

// // func startServer(r ports.PeerRepository, c ports.ConfigurationReader, m ports.MetricReader) error {
// // 	port, err := c.GetPort()
// // 	if err != nil {
// // 		return err
// // 	}

// // 	lis, err := infrastructure.NewNetListener(port)
// // 	if err != nil {
// // 		return err
// // 	}

// // 	grpcServer := adapters.NewGRPC(r, c, m)
// // 	log.Printf("Server listening on port %d", port)
// // 	if err := grpcServer.Serve(lis); err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }
