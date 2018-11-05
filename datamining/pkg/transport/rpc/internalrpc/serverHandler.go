package internalrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	accountlisting "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

type internalSrvHandler struct {
	accLister             accountlisting.Service
	sharedRobotPrivateKey string
	conf                  system.DataMiningConfiguration
}

//NewInternalServerHandler create a new GRPC server handler for account
func NewInternalServerHandler(accLister accountlisting.Service, mine mining.Service, sharedRobotPrivateKey string, conf system.DataMiningConfiguration) api.InternalServer {
	return internalSrvHandler{
		accLister:             accLister,
		sharedRobotPrivateKey: sharedRobotPrivateKey,
		conf:                  conf,
	}
}

//GetAccount implements the protobuf GetAccount request handler
func (s internalSrvHandler) GetAccount(ctx context.Context, req *api.AccountSearchRequest) (*api.AccountSearchResult, error) {
	personHash, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedHashPerson)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	bioWallet, err := s.accLister.GetBiometric(personHash)
	if err != nil {
		return nil, err
	}

	if bioWallet == nil {
		return nil, errors.New(s.conf.Errors.AccountNotExist)
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPrivateKey, bioWallet.CipherAddrRobot())
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	keychain, err := s.accLister.GetKeychain(clearaddr)
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.conf.Errors.AccountNotExist)
	}

	return BuildAccountSearchResult(keychain, bioWallet), nil
}

func (s internalSrvHandler) CreateKeychain(ctx context.Context, req *api.KeychainCreationRequest) (*api.CreationResult, error) {

	keychainRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedKeychainData)
	if err != nil {
		return nil, err
	}

	var keychain *KeychainDataFromJSON
	err = json.Unmarshal([]byte(keychainRawData), &keychain)
	if err != nil {
		return nil, err
	}
	txHashKeychain, err := crypto.NewHasher().HashTransactionData(keychain)
	if err != nil {
		return nil, err
	}

	//TODO: Get elected master and validators ==> contact AI GRPC
	masterIP := "127.0.0.1"
	validators := []string{"127.0.0.1"}

	go func() {
		serverAddr := fmt.Sprintf("%s:%d", masterIP, s.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			log.Fatal(err)
		}

		client := api.NewExternalClient(conn)
		_, err = client.LeadKeychainMining(context.Background(), &api.KeychainLeadRequest{
			EncryptedKeychainData: req.EncryptedKeychainData,
			SignatureKeychainData: req.SignatureKeychainData,
			TransactionHash:       txHashKeychain,
			ValidatorPeerIPs:      validators,
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &api.CreationResult{
		TransactionHash: txHashKeychain,
		MasterPeerIP:    masterIP,
	}, nil
}

func (s internalSrvHandler) CreateBiometric(ctx context.Context, req *api.BiometricCreationRequest) (*api.CreationResult, error) {
	bioRawData, err := crypto.Decrypt(s.sharedRobotPrivateKey, req.EncryptedBiometricData)
	if err != nil {
		return nil, err
	}

	var bio *BioDataFromJSON
	err = json.Unmarshal([]byte(bioRawData), &bio)
	if err != nil {
		return nil, err
	}
	txHashBiometric, err := crypto.NewHasher().HashTransactionData(bio)
	if err != nil {
		return nil, err
	}

	//TODO: Get elected master and validators ==> contact AI GRPC
	masterIP := "127.0.0.1"
	validators := []string{"127.0.0.1"}

	go func() {
		serverAddr := fmt.Sprintf("%s:%d", masterIP, s.conf.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			log.Fatal(err)
		}

		client := api.NewExternalClient(conn)
		_, err = client.LeadBiometricMining(context.Background(), &api.BiometricLeadRequest{
			EncryptedBioData: req.EncryptedBiometricData,
			SignatureBioData: req.SignatureBiometricData,
			TransactionHash:  txHashBiometric,
			ValidatorPeerIPs: validators,
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &api.CreationResult{
		TransactionHash: txHashBiometric,
		MasterPeerIP:    masterIP,
	}, nil
}
