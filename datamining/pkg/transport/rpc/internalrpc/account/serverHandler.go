package internalrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

type internalSrvHandler struct {
	list                  listing.Service
	mine                  mining.Service
	sharedRobotPrivateKey string
	errors                system.DataMininingErrors
}

//NewInternalServerHandler create a new GRPC server handler
func NewInternalServerHandler(list listing.Service, mine mining.Service, sharedRobotPrivateKey string, errors system.DataMininingErrors) api.InternalServer {
	return internalSrvHandler{
		list:                  list,
		mine:                  mine,
		sharedRobotPrivateKey: sharedRobotPrivateKey,
		errors:                errors,
	}
}

//GetAccount implements the protobuf GetAccount request handler
func (s internalSrvHandler) GetAccount(ctx context.Context, req *api.AccountSearchRequest) (*api.AccountSearchResult, error) {
	personHash, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedHashPerson)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	bioWallet, err := s.list.GetBiometric(personHash)
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

	keychain, err := s.list.GetKeychain(clearaddr)
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.errors.AccountNotExist)
	}

	return BuildAccountSearchResult(keychain, bioWallet), nil
}

//CreateAccount implements the protobuf CreateAccount request handler
func (s internalSrvHandler) CreateAccount(ctx context.Context, req *api.AccountCreationRequest) (*api.AccountCreationResult, error) {
	bioRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedBioData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the bio data - %s", err.Error())
	}
	keychainRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedKeychainData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var bioData BioDataFromJSON
	err = json.Unmarshal([]byte(bioRawData), &bioData)
	if err != nil {
		return nil, err
	}

	var keychain *KeychainDataFromJSON
	err = json.Unmarshal([]byte(keychainRawData), &keychain)
	if err != nil {
		return nil, err
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPrivateKey, keychain.EncryptedAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	bData := BuildBioData(bioData, req.SignatureBioData)
	txHashBio, err := crypto.NewHasher().HashTransactionData(bData)
	if err != nil {
		return nil, err
	}

	keychainData := BuildKeychainData(keychain, req.SignatureKeychainData, clearaddr)
	txHashKeychain, err := crypto.NewHasher().HashTransactionData(keychainData)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.mine.LeadMining(txHashKeychain, keychainData.CipherAddrRobot, keychainData.Sigs.BiodSig, keychainData, datamining.CreateKeychainTransaction); err != nil {
			//TODO: handle errors
			log.Fatal(err)
		}
	}()

	go func() {
		if err := s.mine.LeadMining(txHashBio, bData.CipherAddrRobot, bData.Sigs.BiodSig, bData, datamining.CreateBioTransaction); err != nil {
			//TODO: handle errors
			log.Fatal(err)
		}
	}()

	return &api.AccountCreationResult{
		BioTransactionHash:      txHashBio,
		KeychainTransactionHash: txHashKeychain,
	}, nil
}
