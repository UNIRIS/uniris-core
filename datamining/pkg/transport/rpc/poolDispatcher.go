package rpc

import (
	"context"
	"fmt"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/leading"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc/externalrpc"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
	"google.golang.org/grpc"
)

type poolD struct {
	conf system.DataMiningConfiguration
}

//NewPoolDispatcher creates a new pool dispatcher using GRPC
func NewPoolDispatcher(conf system.DataMiningConfiguration) leading.PoolDispatcher {
	return poolD{conf}
}

func (pd poolD) RequestLock(lastValidPool leading.Pool, txLock validating.TransactionLock, sig string) error {

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
func (pd poolD) RequestUnlock(lastValidPool leading.Pool, txLock validating.TransactionLock, sig string) error {

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

func (pd poolD) RequestLastWallet(storagePool leading.Pool, addr string) ([]*datamining.Wallet, error) {

	wallets := make([]*datamining.Wallet, 0)

	//TODO: using goroutines
	for _, p := range storagePool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return nil, err
		}

		client := api.NewExternalClient(conn)
		res, err := client.GetLastWallet(context.Background(), &api.LastWalletRequest{
			Address: addr,
		})
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, externalrpc.BuillDomainWallet(res.Wallet))
	}

	return wallets, nil
}
func (pd poolD) RequestWalletValidation(validPool leading.Pool, d *datamining.WalletData) ([]datamining.Validation, error) {

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
		res, err := client.ValidateWallet(context.Background(), &api.WalletValidationRequest{
			BiodPublicKey:      d.BiodPubk,
			EncryptedAddrRobot: d.CipherAddrRobot,
			EncryptedWallet:    d.CipherWallet,
			PersonPubKey:       d.EmPubk,
			Signatures: &api.Signature{
				Biod:   d.Sigs.BiodSig,
				Person: d.Sigs.EmSig,
			},
		})
		if err != nil {
			return nil, err
		}

		valids = append(valids, externalrpc.BuildDomainValidation(res.Validation))
	}

	return valids, nil
}
func (pd poolD) RequestBioValidation(validPool leading.Pool, b *datamining.BioData) ([]datamining.Validation, error) {
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
		res, err := client.ValidateBio(context.Background(), &api.BioValidationRequest{
			BiodPubKey:         b.BiodPubk,
			BiometricHash:      b.BHash,
			EncryptedAddrBiod:  b.CipherAddrBio,
			EncryptedAESKey:    b.CipherAESKey,
			EncryptedAddrRobot: b.CipherAddrRobot,
			PersonPubKey:       b.EmPubk,
			Signatures: &api.Signature{
				Biod:   b.Sigs.BiodSig,
				Person: b.Sigs.EmSig,
			},
		})
		if err != nil {
			return nil, err
		}

		valids = append(valids, externalrpc.BuildDomainValidation(res.Validation))
	}

	return valids, nil
}
func (pd poolD) RequestWalletStorage(storagePool leading.Pool, wallet *datamining.Wallet) error {

	//TODO: using goroutines
	for _, p := range storagePool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)
		_, err = client.StoreWallet(context.Background(), externalrpc.BuildAPIWalletStoreRequest(wallet))
		if err != nil {
			return err
		}
	}

	return nil
}
func (pd poolD) RequestBioStorage(storagePool leading.Pool, bio *datamining.BioWallet) error {

	//TODO: using goroutines
	for _, p := range storagePool.Peers {
		serverAddr := fmt.Sprintf("%s:%d", p.IP.String(), pd.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			return err
		}

		client := api.NewExternalClient(conn)
		_, err = client.StoreBio(context.Background(), externalrpc.BuildAPIBioStoreRequest(bio))
		if err != nil {
			return err
		}
	}

	return nil
}
