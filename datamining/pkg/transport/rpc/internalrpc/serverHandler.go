package internalrpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

type internalSrvHandler struct {
	list                  listing.Service
	add                   adding.Service
	sharedRobotPrivateKey string
	errors                system.DataMininingErrors
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(list listing.Service, add adding.Service, sharedRobotPrivateKey string, errors system.DataMininingErrors) api.InternalServer {
	return internalSrvHandler{
		list:                  list,
		add:                   add,
		sharedRobotPrivateKey: sharedRobotPrivateKey,
		errors:                errors,
	}
}

//GetWallet implements the protobuf GetWallet request handler
func (s internalSrvHandler) GetWallet(ctx context.Context, req *api.WalletRequest) (*api.WalletResult, error) {
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

	return BuildWalletResult(wallet, bioWallet), nil
}

//StoreWallet implements the protobuf StoreWallet request handler
func (s internalSrvHandler) StoreWallet(ctx context.Context, req *api.Wallet) (*api.StorageResult, error) {
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

	wData := BuildWalletData(wal, req.SignatureWalletData, clearaddr)
	if err := s.add.AddWallet(wData); err != nil {
		return nil, err
	}

	bData := BuildBioData(bio, req.SignatureBioData)
	if err := s.add.AddBioWallet(bData); err != nil {
		return nil, err
	}

	w, err := s.list.GetWallet(clearaddr)
	if err != nil {
		return nil, err
	}

	if w == nil {
		return nil, fmt.Errorf("Cannot find created wallet")
	}

	return &api.StorageResult{
		TransactionHash: w.Endorsement().TransactionHash(),
	}, nil
}
