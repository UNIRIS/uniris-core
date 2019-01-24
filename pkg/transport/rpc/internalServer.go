package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/electing"
	"github.com/uniris/uniris-core/pkg/mining"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/listing"
)

type internalSrv struct {
	lister          listing.Service
	sharedRobotPubK string
	sharedRobotPvk  string
}

//NewInternalServer creates a new GRPC internal server
func NewInternalServer(l listing.Service) api.InternalServiceServer {
	return internalSrv{
		lister: l,
	}
}

func (s internalSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {

	pool, err := electing.FindStoragePool(req.TransactionHash)
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

	var txRaw transaction
	if err := json.Unmarshal([]byte(txJSON), &txRaw); err != nil {
		return nil, err
	}

	masterPeerIP := electing.FindTransactionMasterPeer(txHash)
	minValidations := mining.GetMinimumTransactionValidation(txHash)
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

func (s internalSrv) getID(idAddress string) (*api.IDResponse, error) {
	address, err := crypto.Decrypt(idAddress, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}

	idPool, err := electing.FindStoragePool(address)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:%d", idPool[0].IP().String(), idPool[0].Port())
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
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(reqBytes), s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig
	res, err := cli.GetID(context.Background(), req)
	if err != nil {
		return nil, err
	}

	resBytes, err := json.Marshal(&api.IDResponse{
		EncryptedAesKey: res.EncryptedAesKey,
		Timestamp:       res.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(resBytes), s.sharedRobotPubK, res.SignatureResponse); err != nil {
		return nil, err
	}

	fmt.Printf("GET ID RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
	return res, nil
}

func (s internalSrv) getKeychain(keychainAddr string) (*api.KeychainResponse, error) {
	keychainAddress, err := crypto.Decrypt(keychainAddr, s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	keychainPool, err := electing.FindStoragePool(keychainAddress)
	if err != nil {
		return nil, err
	}

	//Select storage master peer
	serverAddr := fmt.Sprintf("%s:%d", keychainPool[0].IP().String(), keychainPool[0].Port())
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
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(reqBytes), s.sharedRobotPvk)
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig
	res, err := cli.GetKeychain(context.Background(), req)
	if err != nil {
		return nil, err
	}

	resBytes, err := json.Marshal(&api.KeychainResponse{
		EncryptedWallet: res.EncryptedWallet,
		Timestamp:       res.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(resBytes), s.sharedRobotPubK, res.SignatureResponse); err != nil {
		return nil, err
	}

	fmt.Printf("GET KEYCHAIN RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	return res, nil
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
