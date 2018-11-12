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
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	res, err := client.GetAccount(context.Background(), &api.AccountSearchRequest{
		EncryptedHashPerson: encHash,
	})
	if err != nil {
		s, _ := status.FromError(err)
		if s.Message() == c.conf.Datamining.Errors.AccountNotExist {
			return nil, listing.ErrAccountNotExist
		}
		return nil, errors.New(s.Message())
	}

	if err := c.sigHandler.VerifyAccountSearchResultSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return nil, err
	}

	resAcc := &listing.AccountResult{
		EncryptedAddress: res.EncryptedAddress,
		EncryptedAESKey:  res.EncryptedAESkey,
		EncryptedWallet:  res.EncryptedWallet,
	}

	if err := c.sigHandler.SignAccountResult(resAcc, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return resAcc, nil
}

func (c robotClient) AddAccount(req adding.AccountCreationRequest) (*adding.AccountCreationResult, error) {
	serverAddr := fmt.Sprintf("localhost:%d", c.conf.Datamining.InternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewInternalClient(conn)

	resBio, err := client.CreateBiometric(context.Background(), &api.BiometricCreationRequest{
		EncryptedBiometricData: req.EncryptedBioData,
		SignatureBiometricData: &api.Signature{
			Person: req.SignaturesBio.PersonSig,
			Biod:   req.SignaturesBio.BiodSig,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := c.sigHandler.VerifyCreationResultSignature(c.conf.SharedKeys.RobotPublicKey, resBio); err != nil {
		return nil, err
	}

	resKeychain, err := client.CreateKeychain(context.Background(), &api.KeychainCreationRequest{
		EncryptedKeychainData: req.EncryptedKeychainData,
		SignatureKeychainData: &api.Signature{
			Person: req.SignaturesKeychain.PersonSig,
			Biod:   req.SignaturesKeychain.BiodSig,
		},
	})
	if err != nil {
		return nil, err
	}

	txs := adding.AccountCreationTransactions{
		Biometric: adding.TransactionResult{
			TransactionHash: resBio.TransactionHash,
			MasterPeerIP:    resBio.MasterPeerIP,
			Signature:       resBio.Signature,
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

	if err := c.sigHandler.SignAccountCreationResult(res, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}
