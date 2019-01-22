package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/listing"
)

type accountSrv struct {
	lister      listing.Service
	sharedPubK  string
	sharedPvKey string
}

//NewAccountServer creates a new GRPC account server
func NewAccountServer(l listing.Service) api.AccountServiceServer {
	return accountSrv{
		lister: l,
	}
}

func (s accountSrv) GetKeychain(ctx context.Context, req *api.KeychainRequest) (*api.KeychainResponse, error) {
	fmt.Printf("GET KEYCHAIN REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.KeychainRequest{
		EncryptedAddress: req.EncryptedAddress,
		Timestamp:        req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubK, req.SignatureRequest); err != nil {
		return nil, err
	}

	addr, err := crypto.Decrypt(req.EncryptedAddress, s.sharedPvKey)
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

	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}

	res.SignatureResponse = sig
	return res, nil
}

func (s accountSrv) GetID(ctx context.Context, req *api.IDRequest) (*api.IDResponse, error) {
	fmt.Printf("GET ID REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.IDRequest{
		EncryptedAddress: req.EncryptedAddress,
		Timestamp:        req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubK, req.SignatureRequest); err != nil {
		return nil, err
	}

	addr, err := crypto.Decrypt(req.EncryptedAddress, s.sharedPvKey)
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
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig

	return res, nil
}
