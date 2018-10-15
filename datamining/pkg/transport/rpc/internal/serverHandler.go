package internalrpc

import (
	"github.com/uniris/uniris-core/datamining/pkg/validating"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/walletAdding"
	"github.com/uniris/uniris-core/datamining/pkg/walletListing"
)

type internalSrvHandler struct {
	list                  walletlisting.Service
	add                   walletadding.Service
	sharedRobotPrivateKey []byte
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(repo walletlisting.Repository, addRepo walletadding.Repository, sharedRobotPrivateKey []byte) api.InternalServer {
	return internalSrvHandler{
		list:                  walletlisting.NewService(repo),
		add:                   walletadding.NewService(addRepo, validating.NewService()),
		sharedRobotPrivateKey: sharedRobotPrivateKey,
	}
}

//GetWallet implements the protobuf GetWallet request handler
func (s internalSrvHandler) GetWallet(ctx context.Context, req *api.WalletRequest) (*api.WalletResult, error) {
	decrypter := NewDecrypter(s.sharedRobotPrivateKey)
	b := DataBuilder{decrypter}

	bioWallet, err := s.list.GetBioWallet(req.EncryptedHashPerson)
	if err != nil {
		return nil, err
	}

	clearaddr, err := decrypter.Decipher(bioWallet.CipherAddrRobot())
	if err != nil {
		return nil, err
	}
	wallet, err := s.list.GetWallet(clearaddr)
	if err != nil {
		return nil, err
	}

	return b.BuildWalletResult(wallet, bioWallet), nil
}

//StoreWallet implements the protobuf StoreWallet request handler
func (s internalSrvHandler) StoreWallet(ctx context.Context, req *api.Wallet) (*api.StorageResult, error) {
	decrypter := NewDecrypter(s.sharedRobotPrivateKey)
	b := DataBuilder{decrypter}

	walletData, bioData, err := b.BuildWallet(req)
	if err != nil {
		return nil, err
	}

	if err := s.add.AddWallet(walletData); err != nil {
		return nil, err
	}

	if err := s.add.AddBioWallet(bioData); err != nil {
		return nil, err
	}

	//TODO: find the updated hash wallet

	return &api.StorageResult{}, nil
}
