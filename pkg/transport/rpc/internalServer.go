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
	sharedRobotPubK string
	sharedRobotPvk  string
	poolFinding     transaction.PoolFindingService
	mining          transaction.MiningService
}

//NewInternalServer creates a new GRPC internal server
func NewInternalServer(poolFinding transaction.PoolFindingService, mining transaction.MiningService) api.InternalServiceServer {
	return internalSrv{
		poolFinding: poolFinding,
		mining:      mining,
	}
}

func (s internalSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {

	pool, err := s.poolFinding.FindStoragePool(req.TransactionHash)
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
	res, err := cli.GetTransactionStatus(context.Background(), req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("GET TRANSACTION STATUS RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
	resBytes, err := json.Marshal(&api.TransactionStatusResponse{
		Status:    res.Status,
		Timestamp: res.Timestamp,
	})
	if err != nil {
		return nil, err
	}

	if err := crypto.VerifySignature(string(resBytes), s.sharedRobotPubK, res.SignatureResponse); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrv) HandleTransaction(ctx context.Context, req *api.IncomingTransaction) (*api.TransactionResult, error) {
	fmt.Printf("HANDLING TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: check emitter is authorized

	txJSON, err := crypto.Decrypt(req.EncryptedTransaction, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	txHash := crypto.HashString(txJSON)

	var txRaw tx
	if err := json.Unmarshal([]byte(txJSON), &txRaw); err != nil {
		return nil, err
	}

	masterPeerIP := s.poolFinding.FindTransactionMasterPeer(txHash)
	minValidations := s.mining.GetMinimumTransactionValidation(txHash)
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
	fmt.Printf("PRE VALIDATE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	txRes := &api.TransactionResult{
		MasterPeerIp:    masterPeerIP,
		Timestamp:       time.Now().Unix(),
		TransactionHash: txHash,
	}
	txResBytes, err := json.Marshal(txRes)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(txResBytes), s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}

	txRes.Signature = sig
	return txRes, nil
}

func (s internalSrv) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	fmt.Printf("GET ACCOUNT REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.GetAccountRequest{
		EncryptedIdAddress: req.EncryptedIdAddress,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), s.sharedRobotPubK, req.SignatureRequest); err != nil {
		return nil, err
	}

	idAddr, err := crypto.Decrypt(req.EncryptedIdAddress, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	idTx, err := s.poolFinding.RequestLastTransaction(idAddr, transaction.IDType)
	if err != nil {
		return nil, err
	}
	if idTx == nil {
		return nil, status.New(codes.NotFound, "ID does not exist").Err()
	}

	id, err := transaction.NewID(*idTx)
	if err != nil {
		return nil, err
	}

	keychainAddr, err := crypto.Decrypt(id.EncryptedAddrByRobot(), s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	keychainTx, err := s.poolFinding.RequestLastTransaction(keychainAddr, transaction.KeychainType)
	if err != nil {
		return nil, err
	}
	keychain, err := transaction.NewKeychain(*keychainTx)
	if err != nil {
		return nil, err
	}
	if keychainTx == nil {
		return nil, status.New(codes.NotFound, "Keychain does not exist").Err()
	}

	res := &api.GetAccountResponse{
		EncryptedAesKey: id.EncryptedAESKey(),
		EncryptedWallet: keychain.EncryptedWallet(),
		Timestamp:       time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(string(resBytes), s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig
	return res, nil
}

type tx struct {
	Address          string
	Data             string
	Type             int
	Timestamp        int64
	PublicKey        string
	Signature        string
	EmitterSignature string
	Proposal         txProp
}

type txProp struct {
	SharedEmitterKeys txSharedKeys
}

type txSharedKeys struct {
	EncryptedPrivateKey string
	PublicKey           string
}

func formatPreValidationRequest(txRaw tx, txType int, txHash string, minValidations int) *api.PreValidationRequest {
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
