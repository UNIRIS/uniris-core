package rpc

import (
	"context"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/listing"
)

type accountSrv struct {
	lister      listing.Service
	decrypt     Decrypter
	sigHandler  SignatureHandler
	sharedPubK  string
	sharedPvKey string
}

//NewAccountServer creates a new GRPC account server
func NewAccountServer(l listing.Service, d Decrypter, s SignatureHandler) api.AccountServiceServer {
	return accountSrv{
		lister:     l,
		decrypt:    d,
		sigHandler: s,
	}
}

func (s accountSrv) GetKeychain(ctx context.Context, req *api.KeychainRequest) (*api.KeychainResponse, error) {
	fmt.Printf("GET KEYCHAIN REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyKeychainRequestSignature(req, s.sharedPubK); err != nil {
		return nil, err
	}

	addr, err := s.decrypt.DecryptString(req.EncryptedAddress, s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	keychain, err := s.lister.GetKeychain(addr)
	if err != nil {
		return nil, err
	}

	res := &api.KeychainResponse{
		EncryptedWallet: keychain.EncryptedWallet(),
		Timestamp:       time.Now().Unix(),
	}
	if err := s.sigHandler.SignKeychainResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s accountSrv) GetID(ctx context.Context, req *api.IDRequest) (*api.IDResponse, error) {
	fmt.Printf("GET ID REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyIDRequestSignature(req, s.sharedPubK); err != nil {
		return nil, err
	}

	addr, err := s.decrypt.DecryptString(req.EncryptedAddress, s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	id, err := s.lister.GetID(addr)
	if err != nil {
		return nil, err
	}

	res := &api.IDResponse{
		EncryptedAesKey:          id.EncryptedAESKey(),
		EncryptedKeychainAddress: id.EncryptedAddrByRobot(),
		Timestamp:                time.Now().Unix(),
	}

	if err := s.sigHandler.SignIDResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}

	return res, nil
}
