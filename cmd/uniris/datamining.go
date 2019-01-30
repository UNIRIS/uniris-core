package main

import (
	"fmt"
	"log"
	"net"

	"github.com/uniris/uniris-core/pkg/shared"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	memstorage "github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/system"
	"github.com/uniris/uniris-core/pkg/transaction"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
)

func startDatamining(conf system.UnirisConfig) {
	var pInfo discovery.PeerInformer
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
	poolRetr := rpc.NewPoolRetriever("", "")

	sharedSrv := shared.NewService(sharedDb)
	poolFinderSrv := transaction.NewPoolFindingService(poolRetr)
	miningSrv := transaction.NewMiningService(poolReq, poolFinderSrv, sharedSrv, ip.String(), conf.PublicKey, conf.PrivateKey)
	storeSrv := transaction.NewStorageService(txDb, miningSrv)
	lockSrv := transaction.NewLockService(lockDb)

	intSrv := rpc.NewInternalServer(poolFinderSrv, miningSrv, "", "")
	txSrv := rpc.NewTransactionServer(storeSrv, lockSrv, miningSrv, "", "")

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

	log.Printf("Transaction GRPC Server listening on %d", conf.Services.Datamining.ExternalPort)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
