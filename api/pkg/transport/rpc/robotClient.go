package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	adding "github.com/uniris/uniris-core/api/pkg/adding"
	listing "github.com/uniris/uniris-core/api/pkg/listing"
	system "github.com/uniris/uniris-core/api/pkg/system"
	"google.golang.org/grpc"
)

//RobotClient defines wrapper of robot client methods
type RobotClient interface {
	adding.RobotClient
	listing.RobotClient
}

type robotClient struct {
	conf system.UnirisConfig
}

//NewRobotClient creates a new robot client using GRPC
func NewRobotClient(conf system.UnirisConfig) RobotClient {
	return robotClient{conf}
}

func (c robotClient) IsEmitterAuthorized(emPubKey string) error {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewInternalClient(conn)

	res, err := client.IsEmitterAuthorized(context.Background(), &api.AuthorizationRequest{PublicKey: emPubKey})
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if res.Status == false {
		return listing.ErrUnauthorized
	}

	return nil
}

func (c robotClient) GetSharedKeys() (listing.SharedKeys, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	res, err := client.GetSharedKeys(context.Background(), &empty.Empty{})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	emKeys := make([]listing.SharedKeyPair, 0)
	for _, kp := range res.EmitterKeys {
		emKeys = append(emKeys, listing.NewSharedKeyPair(kp.EncryptedPrivateKey, kp.PublicKey))
	}

	return listing.NewSharedKeys(res.RobotPrivateKey, res.RobotPublicKey, emKeys), nil
}

func (c robotClient) GetAccount(encHash string) (listing.AccountResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	res, err := client.GetAccount(context.Background(), &api.AccountSearchRequest{
		EncryptedIDHash: encHash,
	})
	if err != nil {
		s, _ := status.FromError(err)
		if s.Message() == c.conf.Services.Datamining.Errors.AccountNotExist {
			return nil, listing.ErrAccountNotExist
		}
		return nil, errors.New(s.Message())
	}

	return listing.NewAccountResult(res.EncryptedAESkey, res.EncryptedWallet, res.EncryptedAddress, res.Signature), nil
}

func (c robotClient) AddAccount(req adding.AccountCreationRequest) (adding.AccountCreationResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	resID, err := client.CreateID(context.Background(), &api.IDCreationRequest{
		EncryptedID: req.EncryptedID(),
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	resKeychain, err := client.CreateKeychain(context.Background(), &api.KeychainCreationRequest{
		EncryptedKeychain: req.EncryptedKeychain(),
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	txID := adding.NewTransactionResult(resID.TransactionHash, resID.MasterPeerIP, resID.Signature)
	txKeychain := adding.NewTransactionResult(resKeychain.TransactionHash, resKeychain.MasterPeerIP, resKeychain.Signature)

	resTx := adding.NewAccountCreationTransactionResult(txID, txKeychain)
	return adding.NewAccountCreationResult(resTx, ""), nil
}

func (c robotClient) GetTransactionStatus(addr string, txHash string) (listing.TransactionStatus, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return listing.TransactionFailure, err
	}

	client := api.NewInternalClient(conn)
	res, err := client.GetTransactionStatus(context.Background(), &api.TransactionStatusRequest{
		Address: addr,
		Hash:    txHash,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return listing.TransactionFailure, errors.New(s.Message())
	}

	return listing.TransactionStatus(res.Status), nil
}
