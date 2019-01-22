package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"

	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/mining"
	"github.com/uniris/uniris-core/pkg/pooling"
	"github.com/uniris/uniris-core/pkg/storage/mem"
	"github.com/uniris/uniris-core/pkg/transport/http"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
)

func main() {
	go startDiscovery()
	go startDatamining()
	startAPI()
}

func startAPI() {
	r := gin.Default()

	http.NewAccountHandler(r)
	http.NewTransactionHandler(r)

	r.Run(fmt.Sprintf(":%d", 8080))
}

func startDatamining() {
	db := mem.NewDatabase()

	//TODO: implement pool requester

	lister := listing.NewService(db)
	adder := adding.NewService(db, lister)
	pooler := pooling.NewService(nil)
	miner := mining.NewService(pooler, nil, lister, "", "")

	intSrv := rpc.NewInternalServer(lister, pooler)
	txSrv := rpc.NewTransactionServer(adder, lister, miner, "", "")
	accSrv := rpc.NewAccountServer(lister)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 1717))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		api.RegisterInternalServiceServer(grpcServer, intSrv)
		log.Printf("Internal GRPC Server listening on 127.0.0.1:%d", 1717)
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 3535))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	api.RegisterAccountServiceServer(grpcServer, accSrv)

	log.Printf("Transaction and Account GRPC Server listening on 127.0.0.1:%d", 3535)
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func startDiscovery() {

}
