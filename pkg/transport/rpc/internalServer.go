package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/uniris/uniris-core/pkg/inspecting"
	"github.com/uniris/uniris-core/pkg/pooling"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/listing"
)

type internalSrv struct {
	lister          listing.Service
	pooler          pooling.Service
	decrypt         Decrypter
	hasher          Hasher
	sigHandler      SignatureHandler
	sharedRobotPubK string
	sharedRobotPvk  string
}

//NewInternalServer creates a new GRPC internal server
func NewInternalServer(l listing.Service, p pooling.Service, d Decrypter) api.InternalServiceServer {
	return internalSrv{
		lister:  l,
		decrypt: d,
	}
}

func (s internalSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	pool, err := s.pooler.FindStoragePool(req.TransactionHash)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:1717", pool.Peers()[0])
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := api.NewTransactionServiceClient(conn)
	res, err := cli.GetTransactionStatus(context.Background(), req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("GET TRANSACTION STATUS RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
	if err := s.sigHandler.VerifyTransactionStatusResponseSignature(res, s.sharedRobotPubK); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrv) HandleTransaction(ctx context.Context, req *api.IncomingTransaction) (*api.TransactionResult, error) {
	fmt.Printf("HANDLING TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: check emitter is authorized

	txJSON, err := s.decrypt.DecryptString(req.EncryptedTransaction, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}

	var txRaw transaction
	if err := json.Unmarshal([]byte(txJSON), &txRaw); err != nil {
		return nil, err
	}

	txHash := s.hasher.HashString(txJSON)
	masterPeerIP := inspecting.FindTransactionMasterPeer(txHash)
	minValidations := inspecting.GetMinimumTransactionValidation(txHash)
	preValidReq := formatPreValidationRequest(txRaw, int(req.Type), txHash, minValidations)

	serverAddr := fmt.Sprintf("%s:1717", masterPeerIP)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	cli := api.NewTransactionServiceClient(conn)
	res, err := cli.PreValidateTransaction(context.Background(), preValidReq)
	if err != nil {
		return nil, err
	}
	log.Printf("PRE VALIDATE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	txRes := &api.TransactionResult{
		MasterPeerIp:    masterPeerIP,
		Timestamp:       time.Now().Unix(),
		TransactionHash: txHash,
	}
	if err := s.sigHandler.SignTransactionResult(txRes, s.sharedRobotPvk); err != nil {
		return nil, err
	}
	return txRes, nil
}

func (s internalSrv) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	fmt.Printf("GET ACCOUNT REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyAccountRequestSignature(req, s.sharedRobotPubK); err != nil {
		return nil, err
	}

	//TODO: check emitter is authorized

	id, err := s.getID(req.EncryptedIdAddress)
	if err != nil {
		return nil, err
	}

	keychain, err := s.getKeychain(id.EncryptedKeychainAddress)
	if err != nil {
		return nil, err
	}

	res := &api.GetAccountResponse{
		EncryptedAesKey: id.EncryptedAesKey,
		EncryptedWallet: keychain.EncryptedWallet,
		Timestamp:       time.Now().Unix(),
	}
	if err := s.sigHandler.SignAccountResponse(res, s.sharedRobotPvk); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrv) getID(idAddress string) (*api.IDResponse, error) {
	address, err := s.decrypt.DecryptString(idAddress, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}

	idPool, err := s.pooler.FindStoragePool(address)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:1717", idPool.Peers()[0])
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := api.NewAccountServiceClient(conn)
	req := &api.IDRequest{
		EncryptedAddress: idAddress,
		Timestamp:        time.Now().Unix(),
	}
	if err := s.sigHandler.SignIDRequest(req, s.sharedRobotPvk); err != nil {
		return nil, err
	}
	id, err := cli.GetID(context.Background(), req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("GET ID RESPONSE - %s\n", time.Unix(id.Timestamp, 0).String())

	return id, nil
}

func (s internalSrv) getKeychain(keychainAddr string) (*api.KeychainResponse, error) {
	keychainAddress, err := s.decrypt.DecryptString(keychainAddr, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	keychainPool, err := s.pooler.FindStoragePool(keychainAddress)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:1717", keychainPool.Peers()[0])
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := api.NewAccountServiceClient(conn)
	req := &api.KeychainRequest{
		EncryptedAddress: keychainAddr,
		Timestamp:        time.Now().Unix(),
	}
	if err := s.sigHandler.SignKeychainRequest(req, s.sharedRobotPvk); err != nil {
		return nil, err
	}
	keychain, err := cli.GetKeychain(context.Background(), req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("GET KEYCHAIN RESPONSE - %s\n", time.Unix(keychain.Timestamp, 0).String())

	return keychain, nil
}

type transaction struct {
	Address          string
	Data             string
	Type             int
	Timestamp        int64
	PublicKey        string
	Signature        string
	EmitterSignature string
	Proposal         transactionProposal
}

type transactionProposal struct {
	SharedEmitterKeys transactionSharedKeys
}

type transactionSharedKeys struct {
	EncryptedPrivateKey string
	PublicKey           string
}

func formatPreValidationRequest(txRaw transaction, txType int, txHash string, minValidations int) *api.PreValidationRequest {
	return &api.PreValidationRequest{
		MinimumValidations: int32(minValidations),
		Timestamp:          time.Now().Unix(),
		Transaction: &api.Transaction{
			Address:          txRaw.Address,
			Data:             txRaw.Data,
			Type:             api.TransactionType(txType),
			Timestamp:        txRaw.Timestamp,
			PublicKey:        txRaw.PublicKey,
			Signature:        txRaw.Signature,
			EmitterSignature: txRaw.EmitterSignature,
			Proposal: &api.TransactionProposal{
				SharedEmitterKeys: &api.SharedKeys{
					EncryptedPrivateKey: txRaw.Proposal.SharedEmitterKeys.EncryptedPrivateKey,
					PublicKey:           txRaw.Proposal.SharedEmitterKeys.PublicKey,
				},
			},
			TransactionHash: txHash,
		},
	}
}
