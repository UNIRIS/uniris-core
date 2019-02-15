package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/chain"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type intSrv struct {
	techDB shared.TechDatabaseReader
	poolR  consensus.PoolRequester
}

//NewInternalServer creates a new GRPC server for internal communication
func NewInternalServer(tDB shared.TechDatabaseReader, pr consensus.PoolRequester) api.InternalServiceServer {
	return &intSrv{
		techDB: tDB,
		poolR:  pr,
	}
}

func (s intSrv) GetTransactionStatus(ctx context.Context, req *api.InternalTransactionStatusRequest) (*api.TransactionStatusResponse, error) {

	pool, err := consensus.FindStoragePool(req.TransactionAddress)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:%d", pool[0].IP().String(), pool[0].Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	defer conn.Close()

	cli := api.NewChainServiceClient(conn)

	reqStatus := &api.TransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(reqStatus)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	lastMinersKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(reqBytes), lastMinersKeys.PrivateKey())
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

	if err := crypto.VerifySignature(string(resBytes), lastMinersKeys.PublicKey(), res.SignatureResponse); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return res, nil
}

func (s intSrv) HandleTransaction(ctx context.Context, req *api.IncomingTransaction) (*api.TransactionResult, error) {
	fmt.Printf("HANDLING TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	lastMinersKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	txJSON, err := crypto.Decrypt(req.EncryptedTransaction, lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	txHash := crypto.HashString(txJSON)

	var tx txSigned
	if err := json.Unmarshal([]byte(txJSON), &tx); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	masterPeers, err := consensus.FindMasterPeers(txHash)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	minValidations := consensus.GetMinimumValidation(txHash)

	//Building the request to the master miner
	leadReq := formatLeadMiningRequest(tx, txHash, minValidations)
	leadRBytes, err := json.Marshal(leadReq)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	reqSig, err := crypto.Sign(string(leadRBytes), lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	leadReq.SignatureRequest = reqSig

	//TODO: send to all the elected master peers

	//Send the request
	serverAddr := fmt.Sprintf("%s:%d", masterPeers[0].IP(), masterPeers[0].Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	defer conn.Close()
	cli := api.NewMiningServiceClient(conn)
	res, err := cli.LeadTransactionMining(context.Background(), leadReq)
	if err != nil {
		statusCodes, _ := status.FromError(err)
		return nil, status.New(statusCodes.Code(), statusCodes.Message()).Err()
	}
	fmt.Printf("PRE VALIDATE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	txRes := &api.TransactionResult{
		Timestamp:          time.Now().Unix(),
		TransactionReceipt: tx.Address + txHash,
	}
	txResBytes, err := json.Marshal(txRes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(txResBytes), lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	txRes.Signature = sig
	return txRes, nil

}

func (s intSrv) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	fmt.Printf("GET ACCOUNT REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	lastMinersKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	idAddr, err := crypto.Decrypt(req.EncryptedIdAddress, lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	idPool, err := consensus.FindStoragePool(idAddr)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	idTx, err := s.poolR.RequestLastTransaction(idPool, idAddr, chain.IDTransactionType)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if idTx == nil {
		return nil, status.New(codes.NotFound, "ID does not exist").Err()
	}

	id, err := chain.NewID(*idTx)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	keychainAddr, err := crypto.Decrypt(id.EncryptedAddrByMiner(), lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	keychainPool, err := consensus.FindStoragePool(keychainAddr)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	keychainTx, err := s.poolR.RequestLastTransaction(keychainPool, keychainAddr, chain.KeychainTransactionType)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if keychainTx == nil {
		return nil, status.New(codes.NotFound, "Keychain does not exist").Err()
	}
	keychain, err := chain.NewKeychain(*keychainTx)
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

	sig, err := crypto.Sign(string(resBytes), lastMinersKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s intSrv) GetLastSharedKeys(ctx context.Context, req *api.LastSharedKeysRequest) (*api.LastSharedKeys, error) {
	fmt.Printf("GET LAST SHARED KEYS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	authorized, err := shared.IsEmitterKeyAuthorized(req.EmitterPublicKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if !authorized {
		return nil, status.New(codes.PermissionDenied, "emitter not authorized").Err()
	}

	emKeys, err := s.techDB.EmitterKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	emKeyPairs := make([]*api.SharedKeyPair, 0)
	for _, keys := range emKeys {
		emKeyPairs = append(emKeyPairs, &api.SharedKeyPair{
			EncryptedPrivateKey: keys.EncryptedPrivateKey(),
			PublicKey:           keys.PublicKey(),
		})
	}

	lastMinersKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	return &api.LastSharedKeys{
		EmitterKeys:    emKeyPairs,
		MinerPublicKey: lastMinersKeys.PublicKey(),
		Timestamp:      time.Now().Unix(),
	}, nil
}

type txSigned struct {
	Address                   string            `json:"addr"`
	Data                      map[string]string `json:"data"`
	Timestamp                 int64             `json:"timestamp"`
	Type                      int               `json:"type"`
	PublicKey                 string            `json:"public_key"`
	SharedKeysEmitterProposal txSharedKeys      `json:"em_shared_keys_proposal"`
	Signature                 string            `json:"signature"`
	EmitterSignature          string            `json:"em_signature"`
}

type txRaw struct {
	Address                   string            `json:"address"`
	Data                      map[string]string `json:"data"`
	Timestamp                 int64             `json:"timestamp"`
	Type                      int               `json:"type"`
	PublicKey                 string            `json:"public_key"`
	SharedKeysEmitterProposal txSharedKeys      `json:"em_shared_keys_proposal"`
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
			SharedKeysEmitterProposal: &api.SharedKeyPair{
				EncryptedPrivateKey: tx.SharedKeysEmitterProposal.EncryptedPrivateKey,
				PublicKey:           tx.SharedKeysEmitterProposal.PublicKey,
			},
			TransactionHash: txHash,
		},
	}
}
