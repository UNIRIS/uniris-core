package internalrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/leading"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

type internalSrvHandler struct {
	list                  listing.Service
	lead                  leading.Service
	sharedRobotPrivateKey string
	errors                system.DataMininingErrors
	store                 chan string
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(list listing.Service, lead leading.Service, sharedRobotPrivateKey string, errors system.DataMininingErrors, storeCh chan string) api.InternalServer {
	return internalSrvHandler{
		list:                  list,
		lead:                  lead,
		sharedRobotPrivateKey: sharedRobotPrivateKey,
		errors:                errors,
		store:                 storeCh,
	}
}

//GetWallet implements the protobuf GetWallet request handler
func (s internalSrvHandler) GetWallet(ctx context.Context, req *api.WalletSearchRequest) (*api.WalletSearchResult, error) {
	bioHash, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedHashPerson)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	bioWallet, err := s.list.GetBioWallet(bioHash)
	if err != nil {
		return nil, err
	}

	if bioWallet == nil {
		return nil, errors.New(s.errors.AccountNotExist)
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPrivateKey, bioWallet.CipherAddrRobot())
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	wallet, err := s.list.GetWallet(clearaddr)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New(s.errors.AccountNotExist)
	}

	return BuildWalletSearchResult(wallet, bioWallet), nil
}

//StoreWallet implements the protobuf StoreWallet request handler
func (s internalSrvHandler) StoreWallet(ctx context.Context, req *api.WalletStorageRequest) (*api.WalletStorageResult, error) {
	bioRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedBioData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the bio data - %s", err.Error())
	}
	walletRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedWalletData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var bio BioDataFromJSON
	err = json.Unmarshal([]byte(bioRawData), &bio)
	if err != nil {
		return nil, err
	}

	var wal *WalletDataFromJSON
	err = json.Unmarshal([]byte(walletRawData), &wal)
	if err != nil {
		return nil, err
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPrivateKey, wal.EncryptedAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	txHashBio := crypto.HashString(bioRawData)
	txHashWal := crypto.HashString(walletRawData)

	go func() {
		wData := BuildWalletData(wal, req.SignatureWalletData, clearaddr)
		if err := s.lead.ComputeWallet(wData, txHashWal); err != nil {
			//TODO: handle errors
			log.Fatal(err)
		}
		s.store <- txHashWal
	}()

	go func() {
		bData := BuildBioData(bio, req.SignatureBioData)
		if err := s.lead.ComputeBio(bData, txHashBio); err != nil {
			//TODO: handle errors
			log.Fatal(err)
		}
		s.store <- txHashBio
	}()

	return &api.WalletStorageResult{
		BioTransactionHash:  txHashBio,
		DataTransactionHash: txHashWal,
	}, nil
}
