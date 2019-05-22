package rpc

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"plugin"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type txSrv struct {
	chainDB         chainDB
	indexDB         indexDB
	sharedKeyReader sharedKeyReader
	nodeReader      nodeReader
	poolR           PoolRequester
	nodePublicKey   publicKey
	nodePrivateKey  privateKey
	logger          logging.Logger
}

// NewTransactionService creates service handler for the GRPC Transaction service
func NewTransactionService(cDB chainDB, iDB indexDB, skr sharedKeyReader, nr nodeReader, pR PoolRequester, pubk publicKey, pvk privateKey, l logging.Logger) api.TransactionServiceServer {
	return txSrv{
		chainDB:         cDB,
		indexDB:         iDB,
		sharedKeyReader: skr,
		nodeReader:      nr,
		poolR:           pR,
		nodePublicKey:   pubk,
		nodePrivateKey:  pvk,
		logger:          l,
	}
}

func (s txSrv) GetLastTransaction(ctx context.Context, req *api.GetLastTransactionRequest) (*api.GetLastTransactionResponse, error) {
	s.logger.Debug("GET LAST TRANSACTION REQUEST - " + time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.GetLastTransactionRequest{
		TransactionAddress: req.TransactionAddress,
		Type:               req.Type,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "chain/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sym, err := p.Lookup("GetLastTransaction")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	tx, err := sym.(func([]byte, interface{}, interface{}) (interface{}, error))(req.TransactionAddress, s.chainDB, s.indexDB)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if tx == nil {
		return nil, status.New(codes.NotFound, "transaction does not exist").Err()
	}

	tvf, err := formatAPITransaction(tx.(transaction))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	coordStmp, err := formatAPICoordinatorStamp(tx.(transaction).CoordinatorStamp().(coordinatorStamp))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	crossV := make([]*api.ValidationStamp, 0)
	for _, v := range tx.(transaction).CrossValidations() {
		vStamp, err := formatAPIValidation(v.(validationStamp))
		if err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		crossV = append(crossV, vStamp)
	}

	res := &api.GetLastTransactionResponse{
		Timestamp: time.Now().Unix(),
		Transaction: &api.MinedTransaction{
			Transaction:      tvf,
			CoordinatorStamp: coordStmp,
			CrossValidations: crossV,
		},
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s txSrv) GetTransactionStatus(ctx context.Context, req *api.GetTransactionStatusRequest) (*api.GetTransactionStatusResponse, error) {
	s.logger.Debug("GET TRANSACTION STATUS REQUEST - " + time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.GetTransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "chain/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sym, err := p.Lookup("GetTransactionStatus")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	txStatus, err := sym.(func([]byte) (int, error))(req.TransactionHash)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.GetTransactionStatusResponse{
		Status:    api.TransactionStatus(txStatus),
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s txSrv) StoreTransaction(ctx context.Context, req *api.StoreTransactionRequest) (*api.StoreTransactionResponse, error) {
	s.logger.Debug("STORE TRANSACTION REQUEST - " + time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.StoreTransactionRequest{
		MinedTransaction: req.MinedTransaction,
		Timestamp:        req.Timestamp,
	})
	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.CoordinatorStamp, req.MinedTransaction.CrossValidations)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "chain/plugin.so"))
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	sym, err := p.Lookup("StoreTransaction")
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	f := sym.(func(tx interface{}, minV int, chainWriter interface{}) error)

	if err := f(tx, int(req.MinimumValidations), s.chainDB); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.StoreTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) TimeLockTransaction(ctx context.Context, req *api.TimeLockTransactionRequest) (*api.TimeLockTransactionResponse, error) {
	s.logger.Debug("TIMELOCK TRANSACTION REQUEST - " + time.Unix(req.Timestamp, 0).String())

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sym, err := p.Lookup("ParsePublicKey")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	parsePub := sym.(func([]byte) (interface{}, error))

	reqBytes, err := json.Marshal(&api.TimeLockTransactionRequest{
		TransactionHash:     req.TransactionHash,
		Address:             req.Address,
		MasterNodePublicKey: req.MasterNodePublicKey,
		Timestamp:           req.Timestamp,
	})
	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	masterKey, err := parsePub(req.MasterNodePublicKey)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid master public key").Err()
	}

	pTimelock, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "timelock/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	pTimelockSym, err := pTimelock.Lookup("TimeLockTransaction")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	if err := pTimelockSym.(func([]byte, []byte, []byte) error)(req.TransactionHash, req.Address, masterKey.(publicKey).Marshal()); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.TimeLockTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) CoordinateTransaction(ctx context.Context, req *api.CoordinateTransactionRequest) (*api.CoordinateTransactionResponse, error) {
	s.logger.Debug("LEAD TRANSACTION MINING REQUEST - " + time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.CoordinateTransactionRequest{
		Transaction:             req.Transaction,
		MinimumValidations:      req.MinimumValidations,
		Timestamp:               req.Timestamp,
		ElectedCoordinatorNodes: req.ElectedCoordinatorNodes,
	})
	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	wHeaders, err := formatElectedNodeList(req.ElectedCoordinatorNodes)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "mining/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	coordVSym, err := p.Lookup("CoordinateTransactionProcessing")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	coordVF := coordVSym.(func(tx interface{}, nbValidations int, coordList interface{}, nodePv interface{}, nodePub interface{}, nodeReader interface{}, originPublicKeys []interface{}, poolReq interface{}) error)

	authKey, err := s.sharedKeyReader.AuthorizedNodesPublicKeys()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	authKeys := make([]interface{}, len(authKey))
	for i, k := range authKey {
		authKeys[i] = k
	}

	if err := coordVF(tx, int(req.MinimumValidations), wHeaders, s.nodePrivateKey, s.nodePublicKey, s.nodeReader, authKeys, s.poolR); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.CoordinateTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) CrossValidateTransaction(ctx context.Context, req *api.CrossValidateTransactionRequest) (*api.CrossValidateTransactionResponse, error) {
	s.logger.Debug("CONFIRM VALIDATION TRANSACTION REQUEST - " + time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.CrossValidateTransactionRequest{
		Transaction:      req.Transaction,
		CoordinatorStamp: req.CoordinatorStamp,
		Timestamp:        req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	lastPub, lastPv, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if ok, err := lastPub.Verify(reqBytes, req.SignatureRequest); !ok || err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	coordStamp, err := formatCoordinatorStamp(req.CoordinatorStamp)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	tPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	txSym, err := tPlugin.Lookup("NewTransaction")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	txF := txSym.(func(addr []byte, txType int, data map[string]interface{}, timestamp time.Time, pubK interface{}, sig []byte, originSig []byte, coordS interface{}, crossV []interface{}) (interface{}, error))
	minedTx, err := txF(tx.Address(), tx.Type(), tx.Data(), tx.Timestamp(), tx.PreviousPublicKey(), tx.Signature(), tx.OriginSignature(), coordStamp, nil)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "mining/plugin.so"))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	vSym, err := p.Lookup("CrossValidateTransaction")
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	vF := vSym.(func(tx interface{}, nodePub interface{}, nodePv interface{}) (interface{}, error))

	valid, err := vF(minedTx, s.nodePublicKey, s.nodePrivateKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	v, err := formatAPIValidation(valid.(validationStamp))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res := &api.CrossValidateTransactionResponse{
		Validation: v,
		Timestamp:  time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := lastPv.Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}
