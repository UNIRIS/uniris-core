package internalrpc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

type internalSrvHandler struct {
	accountReq account.PoolRequester
	aiClient   AIClient
	hasher     Hasher
	decrypter  rpc.Decrypter
	conf       system.UnirisConfig
}

//NewInternalServerHandler create a new GRPC server handler for account
func NewInternalServerHandler(accountReq account.PoolRequester, aiClient AIClient, h Hasher, d rpc.Decrypter, conf system.UnirisConfig) api.InternalServer {
	return internalSrvHandler{
		accountReq: accountReq,
		aiClient:   aiClient,
		hasher:     h,
		decrypter:  d,
		conf:       conf,
	}
}

//GetAccount implements the protobuf GetAccount request handler
func (s internalSrvHandler) GetAccount(ctx context.Context, req *api.AccountSearchRequest) (*api.AccountSearchResult, error) {
	personHash, err := s.decrypter.DecryptHashPerson(req.EncryptedHashPerson, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	biometricPool, err := s.aiClient.GetBiometricStoragePool(personHash)
	if err != nil {
		return nil, err
	}

	biometric, err := s.accountReq.RequestBiometric(biometricPool, req.EncryptedHashPerson)
	if err != nil {
		return nil, err
	}

	clearAddr, err := s.decrypter.DecryptCipherAddress(biometric.CipherAddrRobot(), s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	keychainPool, err := s.aiClient.GetKeychainStoragePool(clearAddr)
	if err != nil {
		return nil, err
	}

	keychain, err := s.accountReq.RequestKeychain(keychainPool, biometric.CipherAddrRobot())
	if err != nil {
		return nil, err
	}

	return buildAccountSearchResult(keychain, biometric), nil
}

func (s internalSrvHandler) CreateKeychain(ctx context.Context, req *api.KeychainCreationRequest) (*api.CreationResult, error) {

	keychainRawData, err := s.decrypter.DecryptTransactionData(req.EncryptedKeychainData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	var keychain *KeychainDataJSON
	err = json.Unmarshal([]byte(keychainRawData), &keychain)
	if err != nil {
		return nil, err
	}
	txHashKeychain, err := s.hasher.HashKeychainJSON(keychain)
	if err != nil {
		return nil, err
	}

	//TODO: Get elected master and validators ==> contact AI GRPC
	masterIP := "127.0.0.1"
	validators := []string{"127.0.0.1"}

	go func() {
		serverAddr := fmt.Sprintf("%s:%d", masterIP, s.conf.Datamining.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			log.Print(err)
		}

		client := api.NewExternalClient(conn)
		_, err = client.LeadKeychainMining(context.Background(), &api.KeychainLeadRequest{
			EncryptedKeychainData: req.EncryptedKeychainData,
			SignatureKeychainData: req.SignatureKeychainData,
			TransactionHash:       txHashKeychain,
			ValidatorPeerIPs:      validators,
		})
		if err != nil {
			log.Print(err)
		}
	}()

	return &api.CreationResult{
		TransactionHash: txHashKeychain,
		MasterPeerIP:    masterIP,
	}, nil
}

func (s internalSrvHandler) CreateBiometric(ctx context.Context, req *api.BiometricCreationRequest) (*api.CreationResult, error) {
	bioRawData, err := s.decrypter.DecryptTransactionData(req.EncryptedBiometricData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	var bio *BioDataJSON
	err = json.Unmarshal([]byte(bioRawData), &bio)
	if err != nil {
		return nil, err
	}
	txHashBiometric, err := s.hasher.HashBiometricJSON(bio)
	if err != nil {
		return nil, err
	}

	//TODO: Get elected master and validators ==> contact AI GRPC
	masterIP := "127.0.0.1"
	validators := []string{"127.0.0.1"}

	go func() {
		serverAddr := fmt.Sprintf("%s:%d", masterIP, s.conf.Datamining.ExternalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			log.Print(err)
		}

		client := api.NewExternalClient(conn)
		_, err = client.LeadBiometricMining(context.Background(), &api.BiometricLeadRequest{
			EncryptedBioData: req.EncryptedBiometricData,
			SignatureBioData: req.SignatureBiometricData,
			TransactionHash:  txHashBiometric,
			ValidatorPeerIPs: validators,
		})
		if err != nil {
			log.Print(err)
		}
	}()

	return &api.CreationResult{
		TransactionHash: txHashBiometric,
		MasterPeerIP:    masterIP,
	}, nil
}
