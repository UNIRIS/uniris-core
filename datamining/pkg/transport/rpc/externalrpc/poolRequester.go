package externalrpc

import (
	"context"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"google.golang.org/grpc"
)

type poolR struct {
	conf system.DataMiningConfiguration
}

//NewPoolRequester creates a new pool requester using GRPC
func NewPoolRequester(conf system.DataMiningConfiguration) mining.PoolRequester {
	return poolR{conf}
}

func (pR poolR) RequestLock(lastValidPool mining.Pool, txLock lock.TransactionLock, sig string) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers() {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pR.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		_, err = client.LockTransaction(context.Background(), &api.LockRequest{
			MasterRobotKey:  txLock.MasterRobotKey,
			TransactionHash: txLock.TxHash,
			Address:         txLock.Address,
			Signature:       "", //TODO signature
		})

		if err != nil {
			return err
		}
	}

	return nil
}
func (pR poolR) RequestUnlock(lastValidPool mining.Pool, txLock lock.TransactionLock, sig string) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers() {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pR.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		_, err = client.UnlockTransaction(context.Background(), &api.LockRequest{
			MasterRobotKey:  txLock.MasterRobotKey,
			TransactionHash: txLock.TxHash,
			Address:         txLock.Address,
			Signature:       "", //TODO signature
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (pR poolR) RequestValidations(validPool mining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]datamining.Validation, error) {

	valids := make([]datamining.Validation, 0)

	//TODO: using goroutines
	for _, p := range validPool.Peers() {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pR.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return nil, err
		}

		client := api.NewExternalClient(conn)
		var res *api.ValidationResponse

		switch txType {
		case mining.KeychainTransaction:
			res, err = client.ValidateKeychain(context.Background(), &api.KeychainValidationRequest{
				Data:            createKeychainData(data.(*account.KeyChainData)),
				TransactionHash: txHash,
			})
		case mining.BiometricTransaction:
			res, err = client.ValidateBiometric(context.Background(), &api.BiometricValidationRequest{
				Data:            createBiometricData(data.(*account.BioData)),
				TransactionHash: txHash,
			})
		}

		if err != nil {
			return nil, err
		}

		valids = append(valids, datamining.NewValidation(
			datamining.ValidationStatus(res.Validation.Status),
			time.Unix(res.Validation.Timestamp, 0),
			res.Validation.PublicKey,
			res.Validation.Signature))
	}

	return valids, nil
}

func (pR poolR) RequestStorage(sPool mining.Pool, data interface{}, end datamining.Endorsement, txType mining.TransactionType) error {

	//TODO: using goroutines
	for _, p := range sPool.Peers() {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pR.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		switch txType {
		case mining.KeychainTransaction:
			_, err = client.StoreKeychain(context.Background(), &api.KeychainStorageRequest{
				Data:        createKeychainData(data.(*account.KeyChainData)),
				Endorsement: createEndorsement(end.(datamining.Endorsement)),
			})
		case mining.BiometricTransaction:
			_, err = client.StoreBiometric(context.Background(), &api.BiometricStorageRequest{
				Data:        createBiometricData(data.(*account.BioData)),
				Endorsement: createEndorsement(end.(datamining.Endorsement)),
			})
		}

		if err != nil {
			return err
		}
	}

	return nil
}
