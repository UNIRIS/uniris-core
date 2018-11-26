package rpc

import (
	"context"
	"errors"
	"fmt"

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
	conf       system.UnirisConfig
	sigHandler SignatureHandler
}

//NewRobotClient creates a new robot client using GRPC
func NewRobotClient(conf system.UnirisConfig, sigHandler SignatureHandler) RobotClient {
	return robotClient{conf, sigHandler}
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

	if err := c.sigHandler.VerifyAccountSearchResultSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	resAcc := &listing.AccountResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESKey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
	}

	if err := c.sigHandler.SignAccountResult(resAcc, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return resAcc, nil
}

func (c robotClient) AddAccount(req adding.AccountCreationRequest) (*adding.AccountCreationResult, error) {
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

	if err := c.sigHandler.VerifyCreationResultSignature(c.conf.SharedKeys.Robot.PublicKey, resID); err != nil {
		return nil, err
	}

	resKeychain, err := client.CreateKeychain(context.Background(), &api.KeychainCreationRequest{
		EncryptedKeychain: req.EncryptedKeychain,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	txs := adding.AccountCreationTransactionsResult{
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
	}

	res := &adding.AccountCreationResult{
		Transactions: txs,
	}

	if err := c.sigHandler.SignAccountCreationResult(res, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}
