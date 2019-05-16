package main

import (
	"encoding/json"
	fmt "fmt"
	"net"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sharedKeyReader interface {
	LastNodeCrossKeypair() (privateKey, publicKey, error)
}

type publicKey interface {
	Verify(data []byte, sig []byte) (bool, error)
	Encrypt(data []byte) ([]byte, error)
}

type privateKey interface {
	Sign(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type chainReader interface {
	LastTransaction(addr []byte, txType int) (transaction, error)
}

type transaction interface {
	Address() []byte
	Type() int
	Data() map[string]interface{}
	Timestamp() time.Time
	PreviousPublicKey() interface{}
	Signature() []byte
	OriginSignature() []byte
	CoordinatorStamp() interface{}
	CrossValidations() []interface{}
	IsCoordinatorStampValid() (bool, string)
}

type txSrv struct {
	nodePubk  publicKey
	nodePvKey privateKey
	sharedR   sharedKeyReader
	chainR    chainReader
}

func (s txSrv) GetLastTransaction(ctx context.Context, req *GetLastTransactionRequest) (*GetLastTransactionResponse, error) {
	reqBytes, err := json.Marshal(GetLastTransactionRequest{
		TransactionAddress: req.TransactionAddress,
		Type:               req.Type,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	lastPv, lastPub, err := s.sharedR.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}

	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	} else if !ok {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	txRes, err := s.chainR.LastTransaction(req.TransactionAddress, int(req.Type))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if txRes == nil {
		return nil, status.New(codes.NotFound, "transaction does not exist").Err()
	}
	tx, ok := txRes.(transaction)
	if !ok {
		return nil, status.New(codes.Internal, "last transaction result type is invalid").Err()
	}

	tvf, err := formatAPITransaction(tx)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &GetLastTransactionResponse{
		Timestamp:   time.Now().Unix(),
		Transaction: tvf,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s txSrv) GetTransactionStatus(ctx context.Context, req *GetTransactionStatusRequest) (*GetTransactionStatusResponse, error) {
	return nil, nil
}

func (s txSrv) StoreTransaction(ctx context.Context, req *StoreTransactionRequest) (*StoreTransactionResponse, error) {
	return nil, nil
}

func (s txSrv) TimeLockTransaction(ctx context.Context, req *TimeLockTransactionRequest) (*TimeLockTransactionResponse, error) {
	return nil, nil
}

func (s txSrv) CoordinateTransaction(ctx context.Context, req *CoordinateTransactionRequest) (*CoordinateTransactionResponse, error) {
	return nil, nil
}

func (s txSrv) CrossValidationTransaction(ctx context.Context, req *CrossTransactionValidationRequest) (*CrossTransactionValidationResponse, error) {
	return nil, nil
}

//StartServer creates and start the transaction service GRPC on the given port
func StartServer(port int, nodePvKey interface{}, nodePubKey interface{}, sharedR interface{}, chainR interface{}) error {

	//TODO: uncomment when the main use implementation of the interfaces
	// if _, ok := nodePubKey.(publicKey); !ok {
	// 	return errors.New("invalid node public key")
	// }

	// if _, ok := nodePvKey.(privateKey); !ok {
	// 	return errors.New("invalid node public key")
	// }

	// if _, ok := sharedR.(sharedKeyReader); !ok {
	// 	return errors.New("invalid shared key reader type")
	// }

	// if _, ok := chainR.(chainReader); !ok {
	// 	return errors.New("invalid chain reader type")
	// }

	grpcServer := grpc.NewServer()
	RegisterTransactionServiceServer(grpcServer, txSrv{
		// nodePubk:  nodePubKey.(publicKey),
		// nodePvKey: nodePvKey.(privateKey),
		// sharedR:   sharedR.(sharedKeyReader),
		// chainR:    chainR.(chainReader),
	})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	fmt.Printf("GRPC server listening on %d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
