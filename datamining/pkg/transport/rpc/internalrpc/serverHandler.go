package internalrpc

import (
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

type internalSrvHandler struct {
	list                  listing.Service
	add                   adding.Service
	sharedRobotPrivateKey []byte
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(listservice listing.Service, addservice adding.Service, sharedRobotPrivateKey []byte) api.InternalServer {
	return internalSrvHandler{
		list:                  listservice,
		add:                   addservice,
		sharedRobotPrivateKey: sharedRobotPrivateKey,
	}
}

//GetWallet implements the protobuf GetWallet request handler
func (s internalSrvHandler) GetWallet(ctx context.Context, req *api.WalletRequest) (*api.WalletResult, error) {
	bioWallet, err := s.list.GetBioWallet(req.EncryptedHashPerson)
	if err != nil {
		return nil, err
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPrivateKey, bioWallet.CipherAddrRobot())
	if err != nil {
		return nil, err
	}
	wallet, err := s.list.GetWallet(clearaddr)
	if err != nil {
		return nil, err
	}

	return BuildWalletResult(wallet, bioWallet), nil
}

//StoreWallet implements the protobuf StoreWallet request handler
func (s internalSrvHandler) StoreWallet(ctx context.Context, req *api.Wallet) (*api.StorageResult, error) {
	bioRawData, walletRawData, err := DecryptWallet(req, s.sharedRobotPrivateKey)
	if err != nil {
		return nil, err
	}

	jsonBio, err := DecodeBioData(bioRawData, req.SignatureBioData)
	if err != nil {
		return nil, err
	}

	jsonWal, err := DecodeWalletData(walletRawData, req.SignatureWalletData)
	if err != nil {
		return nil, err
	}

	if err := s.add.AddWallet(BuildWalletData(jsonWal, req.SignatureWalletData)); err != nil {
		return nil, err
	}

	if err := s.add.AddBioWallet(BuildBioData(jsonBio, req.SignatureBioData)); err != nil {
		return nil, err
	}

	//TODO: calculate the updated hash wallet

	return &api.StorageResult{}, nil
}
