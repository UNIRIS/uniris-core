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

	_, err = client.IsEmitterAuthorized(context.Background(), &api.AuthorizationRequest{PublicKey: emPubKey})
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	return nil
}

func (c robotClient) GetSharedKeys() (*listing.SharedKeysResult, error) {
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
		emKeys = append(emKeys, listing.SharedKeyPair{
			EncryptedPrivateKey: kp.EncryptedPrivateKey,
			PublicKey:           kp.PublicKey,
		})
	}
	return &listing.SharedKeysResult{
		RobotPublicKey:  res.RobotPublicKey,
		RobotPrivateKey: res.RobotPrivateKey,
		EmitterKeys:     emKeys,
	}, nil
}

func (c robotClient) GetAccount(encHash string) (*listing.AccountResult, error) {
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

	resAcc := &listing.AccountResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESKey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
		Signature:        res.Signature,
	}

	return resAcc, nil
}

func (c robotClient) AddAccount(req *adding.AccountCreationRequest) (*adding.AccountCreationResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Services.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	resID, err := client.CreateID(context.Background(), &api.IDCreationRequest{
		EncryptedID: req.EncryptedID,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	resKeychain, err := client.CreateKeychain(context.Background(), &api.KeychainCreationRequest{
		EncryptedKeychain: req.EncryptedKeychain,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	res := &adding.AccountCreationResult{
		Transactions: adding.AccountCreationTransactionsResult{
			ID: adding.TransactionResult{
				TransactionHash: resID.TransactionHash,
				MasterPeerIP:    resID.MasterPeerIP,
				Signature:       resID.Signature,
			},
			Keychain: adding.TransactionResult{
				TransactionHash: resKeychain.TransactionHash,
				MasterPeerIP:    resKeychain.MasterPeerIP,
				Signature:       resKeychain.Signature,
			},
		},
	}

	return res, nil
}
