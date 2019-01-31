package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type internalSrv struct {
	sharedPubKey string
	sharedPvKey  string
	poolFinding  transaction.PoolFindingService
	mining       transaction.MiningService
}

//NewInternalServer creates a new GRPC internal server
func NewInternalServer(poolFinding transaction.PoolFindingService, mining transaction.MiningService, sharedPubk, sharedPvk string) api.InternalServiceServer {
	return internalSrv{
		poolFinding:  poolFinding,
		mining:       mining,
		sharedPubKey: sharedPubk,
		sharedPvKey:  sharedPvk,
	}
}

func (s internalSrv) GetTransactionStatus(ctx context.Context, req *api.InternalTransactionStatusRequest) (*api.TransactionStatusResponse, error) {

	pool, err := s.poolFinding.FindStoragePool(req.TransactionAddress)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:%d", pool[0].IP().String(), pool[0].Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli := api.NewTransactionServiceClient(conn)

	reqStatus := &api.TransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(reqStatus)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(reqBytes), s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	reqStatus.SignatureRequest = sig

	res, err := cli.GetTransactionStatus(context.Background(), reqStatus)
	if err != nil {
		grpcErr, _ := status.FromError(err)
		return nil, status.New(grpcErr.Code(), grpcErr.Message()).Err()
	}

	fmt.Printf("GET TRANSACTION STATUS RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
	resBytes, err := json.Marshal(&api.TransactionStatusResponse{
		Status:    res.Status,
		Timestamp: res.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if err := crypto.VerifySignature(string(resBytes), s.sharedPubKey, res.SignatureResponse); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return res, nil
}

func (s internalSrv) HandleTransaction(ctx context.Context, req *api.IncomingTransaction) (*api.TransactionResult, error) {
	fmt.Printf("HANDLING TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: check emitter is authorized

	txJSON, err := crypto.Decrypt(req.EncryptedTransaction, s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	txHash := crypto.HashString(txJSON)

	var tx txSigned
	if err := json.Unmarshal([]byte(txJSON), &tx); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	masterPeerIP, masterPeerPort := s.poolFinding.FindTransactionMasterPeer(txHash)
	minValidations := s.mining.GetMinimumTransactionValidation(txHash)

	//Building the request to the master miner
	leadReq := formatLeadMiningRequest(tx, txHash, minValidations)
	leadRBytes, err := json.Marshal(leadReq)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	reqSig, err := crypto.Sign(string(leadRBytes), s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	leadReq.SignatureRequest = reqSig

	//Send the request
	serverAddr := fmt.Sprintf("%s:%d", masterPeerIP, masterPeerPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	defer conn.Close()
	cli := api.NewTransactionServiceClient(conn)
	res, err := cli.LeadTransactionMining(context.Background(), leadReq)
	if err != nil {
		statusCodes, _ := status.FromError(err)
		return nil, status.New(statusCodes.Code(), err.Error()).Err()
	}
	fmt.Printf("PRE VALIDATE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	txRes := &api.TransactionResult{
		Timestamp:       time.Now().Unix(),
		TransactionHash: txHash,
	}
	txResBytes, err := json.Marshal(txRes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(txResBytes), s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	txRes.Signature = sig
	return txRes, nil
}

func (s internalSrv) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	fmt.Printf("GET ACCOUNT REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	idAddr, err := crypto.Decrypt(req.EncryptedIdAddress, s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	idTx, err := s.poolFinding.RequestLastTransaction(idAddr, transaction.IDType)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if idTx == nil {
		return nil, status.New(codes.NotFound, "ID does not exist").Err()
	}

	id, err := transaction.NewID(*idTx)
	if err != nil {
		return nil, err
	}

	keychainAddr, err := crypto.Decrypt(id.EncryptedAddrByRobot(), s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	keychainTx, err := s.poolFinding.RequestLastTransaction(keychainAddr, transaction.KeychainType)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if keychainTx == nil {
		return nil, status.New(codes.NotFound, "Keychain does not exist").Err()
	}
	keychain, err := transaction.NewKeychain(*keychainTx)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.GetAccountResponse{
		EncryptedAesKey: id.EncryptedAESKey(),
		EncryptedWallet: keychain.EncryptedWallet(),
		Timestamp:       time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

type txSigned struct {
	Address          string            `json:"address"`
	Data             map[string]string `json:"data"`
	Timestamp        int64             `json:"timestamp"`
	Type             int               `json:"type"`
	PublicKey        string            `json:"public_key"`
	Proposal         txProp            `json:"proposal"`
	Signature        string            `json:"signature"`
	EmitterSignature string            `json:"em_signature"`
}

type txRaw struct {
	Address   string            `json:"address"`
	Data      map[string]string `json:"data"`
	Timestamp int64             `json:"timestamp"`
	Type      int               `json:"type"`
	PublicKey string            `json:"public_key"`
	Proposal  txProp            `json:"proposal"`
}

type txProp struct {
	SharedEmitterKeys txSharedKeys `json:"shared_emitter_keys"`
}

type txSharedKeys struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}

func formatLeadMiningRequest(tx txSigned, txHash string, minValidations int) *api.LeadMiningRequest {
	return &api.LeadMiningRequest{
		MinimumValidations: int32(minValidations),
		Timestamp:          time.Now().Unix(),
		Transaction: &api.Transaction{
			Address:          tx.Address,
			Data:             tx.Data,
			Type:             api.TransactionType(tx.Type),
			Timestamp:        tx.Timestamp,
			PublicKey:        tx.PublicKey,
			Signature:        tx.Signature,
			EmitterSignature: tx.EmitterSignature,
			Proposal: &api.TransactionProposal{
				SharedEmitterKeys: &api.SharedKeys{
					EncryptedPrivateKey: tx.Proposal.SharedEmitterKeys.EncryptedPrivateKey,
					PublicKey:           tx.Proposal.SharedEmitterKeys.PublicKey,
				},
			},
			TransactionHash: txHash,
		},
	}
}
