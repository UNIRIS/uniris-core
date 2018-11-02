package rpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/status"

	"github.com/uniris/uniris-core/api/pkg/crypto"
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
	conf                  system.DataMiningConfiguration
	robotSharedPrivateKey string
}

//NewRobotClient creates a new robot client using GRPC
func NewRobotClient(conf system.DataMiningConfiguration, robotSharedPrivateKey string) RobotClient {
	return robotClient{conf, robotSharedPrivateKey}
}

func (c robotClient) GetAccount(encHash string) (*listing.SignedAccountResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	resGRPC, err := client.GetAccount(context.Background(), &api.AccountSearchRequest{
		EncryptedHashPerson: encHash,
	})
	if err != nil {
		s, _ := status.FromError(err)
		if s.Message() == c.conf.Errors.AccountNotExist {
			return nil, listing.ErrAccountNotExist
		}
		return nil, errors.New(s.Message())
	}

	r := listing.AccountResult{
		EncryptedAESKey:  resGRPC.EncryptedAESkey,
		EncryptedAddress: resGRPC.EncryptedAddress,
		EncryptedWallet:  resGRPC.EncryptedWallet,
	}

	sig, err := crypto.SignData(c.robotSharedPrivateKey, r)
	if err != nil {
		return nil, err
	}

	return &listing.SignedAccountResult{
		EncryptedAddress: r.EncryptedAddress,
		EncryptedAESKey:  r.EncryptedAESKey,
		EncryptedWallet:  r.EncryptedWallet,
		SignatureRequest: sig,
	}, nil
}

func (c robotClient) AddAccount(req adding.AccountCreationRequest) (*adding.AccountCreationResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	w := &api.AccountCreationRequest{
		EncryptedBioData:      req.EncryptedBioData,
		EncryptedKeychainData: req.EncryptedKeychainData,
		SignatureBioData: &api.Signature{
			Person: req.SignaturesBio.PersonSig,
			Biod:   req.SignaturesBio.BiodSig,
		},
		SignatureKeychainData: &api.Signature{
			Person: req.SignaturesKeychain.PersonSig,
			Biod:   req.SignaturesKeychain.BiodSig,
		},
	}

	resGRPC, err := client.CreateAccount(context.Background(), w)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	txs := adding.AccountCreationTransactions{
		Biod:     resGRPC.BioTransactionHash,
		Keychain: resGRPC.KeychainTransactionHash,
	}

	sig, err := crypto.SignData(c.robotSharedPrivateKey, txs)
	if err != nil {
		return nil, err
	}

	return &adding.AccountCreationResult{
		Transactions: txs,
		Signature:    sig,
	}, nil
}

func (c robotClient) GetMasterPeer() (listing.MasterPeer, error) {
	//TODO: implement with AI GRPC
	return listing.MasterPeer{
		IP: "127.0.0.1",
	}, nil
}
