package internalrpc

import (
	"encoding/json"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

type internalSrvHandler struct {
	list                  listing.Service
	add                   adding.Service
	sharedRobotPrivateKey []byte
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(list listing.Service, add adding.Service, sharedRobotPublicKey []byte, sharedRobotPrivateKey []byte) api.InternalServer {
	return internalSrvHandler{
		list:                  list,
		add:                   add,
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
	bioRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedBioData)
	if err != nil {
		return nil, err
	}
	walletRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedWalletData)
	if err != nil {
		return nil, err
	}

	var bio BioDataFromJSON
	err = json.Unmarshal(bioRawData, &bio)
	if err != nil {
		return nil, err
	}

	var wal WalletDataFromJSON
	err = json.Unmarshal(walletRawData, &wal)
	if err != nil {
		return nil, err
	}

	if err := s.add.AddWallet(BuildWalletData(wal, req.SignatureWalletData)); err != nil {
		return nil, err
	}

	if err := s.add.AddBioWallet(BuildBioData(bio, req.SignatureBioData)); err != nil {
		return nil, err
	}

	addr, err := crypto.Decrypt(s.sharedRobotPrivateKey, []byte(bio.EncryptedAddrRobot))
	if err != nil {
		return nil, err
	}

	w, err := s.list.GetWallet(addr)
	if err != nil {
		return nil, err
	}

	bWallet, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}
	return &api.StorageResult{
		HashUpdatedWallet: crypto.Hash(bWallet),
	}, nil
}
