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

	resGRPC, err := client.GetWallet(context.Background(), &api.WalletRequest{
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
		EncryptedAESKey:     string(resGRPC.EncryptedAESkey),
		EncryptedAddrPerson: string(resGRPC.EncryptedWalletAddress),
		EncryptedWallet:     string(resGRPC.EncryptedWallet),
	}

	sig, err := crypto.SignData(c.robotSharedPrivateKey, r)
	if err != nil {
		return nil, err
	}

	return &listing.SignedAccountResult{
		EncryptedAddrPerson: r.EncryptedAddrPerson,
		EncryptedAESKey:     r.EncryptedAESKey,
		EncryptedWallet:     r.EncryptedWallet,
		SignatureRequest:    sig,
	}, nil
}

func (c robotClient) AddAccount(req adding.EnrollmentRequest) (*adding.EnrollmentResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	w := &api.Wallet{
		EncryptedBioData:    req.EncryptedBioData,
		EncryptedWalletData: req.EncryptedWalletData,
		SignatureBioData: &api.Signature{
			Person: req.SignaturesBio.PersonSig,
			Biod:   req.SignaturesBio.BiodSig,
		},
		SignatureWalletData: &api.Signature{
			Person: req.SignaturesWallet.PersonSig,
			Biod:   req.SignaturesWallet.BiodSig,
		},
	}

	resGRPC, err := client.StoreWallet(context.Background(), w)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	sig, err := crypto.SignRaw(c.robotSharedPrivateKey, resGRPC.TransactionHash)
	if err != nil {
		return nil, err
	}

	return &adding.EnrollmentResult{
		TransactionHash:  string(resGRPC.TransactionHash),
		SignatureRequest: string(sig),
	}, nil
}
