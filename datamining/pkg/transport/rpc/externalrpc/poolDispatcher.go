package externalrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/any"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/mining/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/pool"
	"github.com/uniris/uniris-core/datamining/pkg/mining/transactions"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"google.golang.org/grpc"
)

type poolD struct {
	conf system.DataMiningConfiguration
}

//NewPoolDispatcher creates a new pool dispatcher using GRPC
func NewPoolDispatcher(conf system.DataMiningConfiguration) mining.PoolDispatcher {
	return poolD{conf}
}

func (pd poolD) RequestLock(lastValidPool pool.PeerCluster, txLock lock.TransactionLock, sig string) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		_, err = client.LockTransaction(context.Background(), &api.LockRequest{
			MasterRobotKey:  txLock.MasterRobotKey,
			TransactionHash: txLock.TxHash,
			Signature:       "", //TODO signature
		})

		if err != nil {
			return err
		}
	}

	return nil
}
func (pd poolD) RequestUnlock(lastValidPool pool.PeerCluster, txLock lock.TransactionLock, sig string) error {

	//TODO: using goroutines
	for _, p := range lastValidPool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		_, err = client.UnlockTransaction(context.Background(), &api.LockRequest{
			MasterRobotKey:  txLock.MasterRobotKey,
			TransactionHash: txLock.TxHash,
			Signature:       "", //TODO signature
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (pd poolD) RequestValidations(validPool pool.PeerCluster, data interface{}, txType transactions.Type) ([]datamining.Validation, error) {

	valids := make([]datamining.Validation, 0)

	//TODO: using goroutines
	for _, p := range validPool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return nil, err
		}

		client := api.NewExternalClient(conn)
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		res, err := client.Validate(context.Background(), &api.ValidationRequest{
			Data:            &any.Any{Value: b},
			TransactionType: api.TransactionType(txType),
		})
		if err != nil {
			return nil, err
		}

		valids = append(valids, BuildDomainValidation(res.Validation))
	}

	return valids, nil
}

func (pd poolD) RequestStorage(sPool pool.PeerCluster, data interface{}, txType transactions.Type) error {

	//TODO: using goroutines
	for _, p := range sPool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)

		b, err := json.Marshal(data)
		if err != nil {
			return err
		}

		_, err = client.Store(context.Background(), &api.StorageRequest{
			Data:            &any.Any{Value: b},
			TransactionType: api.TransactionType(txType),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
