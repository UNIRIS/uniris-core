package main

import (
	"fmt"
	"log"
	"net"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/gossip"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/mining"
	memstorage "github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/system"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
)

func startDatamining(conf system.UnirisConfig) {
	var pInfo gossip.PeerInformer
	if conf.Network.Type == "private" {
		pInfo = system.NewPeerInformer(true, conf.Network.Interface)
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

	lister := listing.NewService(txDb, lockDb, sharedDb)
	adder := adding.NewService(txDb, lockDb, sharedDb, lister)
	miner := mining.NewService(lister, poolReq, "", "", ip.String())

	intSrv := rpc.NewInternalServer(lister)
	txSrv := rpc.NewTransactionServer(adder, lister, miner, "", "")
	accSrv := rpc.NewAccountServer(lister)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conf.Services.Datamining.InternalPort))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		api.RegisterInternalServiceServer(grpcServer, intSrv)
		log.Printf("Internal GRPC Server listening on %d", conf.Services.Datamining.InternalPort)
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", conf.Services.Datamining.ExternalPort))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	api.RegisterAccountServiceServer(grpcServer, accSrv)

	log.Printf("Transaction,Account GRPC Server listening on %d", conf.Services.Datamining.ExternalPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
