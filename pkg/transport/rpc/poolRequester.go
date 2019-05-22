package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync/atomic"
	"time"

	"github.com/uniris/uniris-core/pkg/logging"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"google.golang.org/grpc"
)

type PoolRequester struct {
	SharedKeyReader sharedKeyReader
	nodeReader      nodeReader
	Logger          logging.Logger
}

func (pr PoolRequester) RequestLastTransaction(pool electedNodeList, txAddr []byte, txType int) (transaction, error) {

	lastPub, lastPv, err := pr.SharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}

	ackChan := make(chan bool)
	var nbReqAcks int32
	var nbReqFailures int32
	var failed bool
	minReplies := 1 //TODO: specify minium required transaction replies

	req := &api.GetLastTransactionRequest{
		TransactionAddress: txAddr,
		Type:               api.TransactionType(txType),
		Timestamp:          time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sig, err := lastPv.Sign(reqBytes)
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig

	txRes := make([]transaction, 0)

	for _, p := range pool.Nodes() {
		go func(p electedNode) {

			n, err := pr.nodeReader.FindByPublicKey(p.PublicKey().(publicKey).(publicKey))
			if err != nil {
				ackChan <- false
				return
			}

			serverAddr := fmt.Sprintf("%s:%d", n.IP(), n.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				pr.Logger.Error("GET LAST TRANSACTION - ERROR: " + err.Error())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.GetLastTransaction(context.Background(), req)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					ackChan <- true
					return
				}
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("GET LAST TRANSACTION RESPONSE - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}

			pr.Logger.Debug("GET LAST TRANSACTION RESPONSE - " + time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.GetLastTransactionResponse{
				Timestamp:   res.Timestamp,
				Transaction: res.Transaction,
			})
			if err != nil {
				pr.Logger.Error("GET LAST TRANSACTION RESPONSE - ERROR: " + err.Error())
				ackChan <- false
				return
			}

			if ok, err := lastPub.Verify(resBytes, res.SignatureResponse); !ok || err != nil {
				pr.Logger.Error("GET LAST TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			if (res.Transaction != &api.MinedTransaction{}) {
				tx, err := formatMinedTransaction(res.Transaction.Transaction, res.Transaction.CoordinatorStamp, res.Transaction.CrossValidations)
				if err != nil {
					pr.Logger.Error("GET LAST TRANSACTION RESPONSE - ERROR: " + err.Error())
					ackChan <- false
					return
				}
				txRes = append(txRes, tx)
			}
			ackChan <- true
		}(p.(electedNode))
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbReqAcks, 1)
		} else {
			atomic.AddInt32(&nbReqFailures, 1)
		}
		if atomic.LoadInt32(&nbReqFailures) == int32(len(pool.Nodes())) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbReqAcks) == int32(minReplies) {
			break
		}
	}

	if failed {
		return nil, errors.New("transaction request failed")
	}

	if len(txRes) == 0 {
		return nil, nil
	}

	//TODO: consensus to implement to get the right result
	return txRes[0], nil
}

func (pr PoolRequester) RequestTransactionTimeLock(pool electedNodeList, txHash []byte, txAddress []byte, masterPublicKey publicKey) error {

	lastPub, lastPv, err := pr.SharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return err
	}

	ackChan := make(chan bool)
	var nbLockAcks int32
	var nbLockFailures int32
	var failed bool

	minLocks := 1 //TODO: specify minium required timelocks

	req := &api.TimeLockTransactionRequest{
		Address:             txAddress,
		TransactionHash:     txHash,
		MasterNodePublicKey: masterPublicKey.Marshal(),
		Timestamp:           time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := lastPv.Sign(reqBytes)
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool.Nodes() {
		go func(p electedNode) {

			n, err := pr.nodeReader.FindByPublicKey(p.PublicKey().(publicKey))
			if err != nil {
				ackChan <- false
				return
			}

			serverAddr := fmt.Sprintf("%s:%d", n.IP(), n.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("TIMELOCK TRANSACTION REQUEST - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.TimeLockTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("TIMELOCK TRANSACTION RESPONSE - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}

			pr.Logger.Debug("TIMELOCK TRANSACTION RESPONSE - " + time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.TimeLockTransactionResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				pr.Logger.Error("TIMELOCK TRANSACTION RESPONSE - ERROR: " + err.Error())
				ackChan <- false
				return
			}
			if ok, err := lastPub.Verify(resBytes, res.SignatureResponse); !ok || err != nil {
				pr.Logger.Error("LOCK TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			ackChan <- true
		}(p.(electedNode))
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbLockAcks, 1)
		} else {
			atomic.AddInt32(&nbLockFailures, 1)
		}
		if atomic.LoadInt32(&nbLockFailures) == int32(len(pool.Nodes())) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbLockAcks) == int32(minLocks) {
			break
		}
	}

	if failed {
		return errors.New("transaction timelock failed")
	}

	return nil
}

func (pr PoolRequester) RequestTransactionValidations(pool electedNodeList, tx transaction, minValids int, masterValid coordinatorStamp) ([]validationStamp, error) {

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}
	sym, err := p.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePub := sym.(func([]byte) (interface{}, error))

	vPlugin, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "validationStamp/plugin.so"))
	if err != nil {
		return nil, err
	}
	vPlugSym, err := vPlugin.Lookup("NewValidationStamp")
	if err != nil {
		return nil, err
	}
	newValidF := vPlugSym.(func(status int, t time.Time, nodePubk interface{}, nodeSig []byte) (interface{}, error))

	lastPub, lastPv, err := pr.SharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}

	txf, err := formatAPITransaction(tx)
	if err != nil {
		return nil, err
	}
	cs, err := formatAPICoordinatorStamp(masterValid)
	if err != nil {
		return nil, err
	}

	req := &api.CrossValidateTransactionRequest{
		CoordinatorStamp: cs,
		Transaction:      txf,
		Timestamp:        time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	sig, err := lastPv.Sign(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	req.SignatureRequest = sig

	validations := make([]validationStamp, 0)

	ackChan := make(chan bool)
	var nbValidationAck int32
	var nbValidationFailures int32
	var failed bool

	for _, p := range pool.Nodes() {

		go func(p electedNode) {

			n, err := pr.nodeReader.FindByPublicKey(p.PublicKey().(publicKey))
			if err != nil {
				ackChan <- false
				return
			}

			serverAddr := fmt.Sprintf("%s:%d", n.IP().String(), n.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.CrossValidateTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}

			pr.Logger.Debug("CONFIRM VALIDATION TRANSACTION RESPONSE - " + time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.CrossValidateTransactionResponse{
				Timestamp:  res.Timestamp,
				Validation: res.Validation,
			})
			if err != nil {
				pr.Logger.Error("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: " + err.Error())
				ackChan <- false
				return
			}
			if ok, err := lastPub.Verify(resBytes, res.SignatureResponse); !ok || err != nil {
				pr.Logger.Error("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			vKey, err := parsePub(res.Validation.PublicKey)
			if err != nil {
				pr.Logger.Error("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: invalid validator public key")
				ackChan <- false
				return
			}

			v, err := newValidF(int(res.Validation.Status), time.Unix(res.Timestamp, 0), vKey, res.Validation.Signature)
			if err != nil {
				ackChan <- false
				return
			}
			validations = append(validations, v.(validationStamp))
			ackChan <- true
		}(p.(electedNode))
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbValidationAck, 1)
		} else {
			atomic.AddInt32(&nbValidationFailures, 1)
		}
		if atomic.LoadInt32(&nbValidationFailures) == int32(len(pool.Nodes())) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbValidationAck) == int32(minValids) {
			break
		}
	}

	if failed {
		return nil, errors.New("transaction storage failed")
	}

	return validations, nil
}

func (pr PoolRequester) RequestTransactionStorage(pool electedNodeList, minStorage int, tx transaction) error {

	lastPub, lastPv, err := pr.SharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	confValids := make([]*api.ValidationStamp, 0)
	for _, v := range tx.CrossValidations() {
		vf, err := formatAPIValidation(v.(validationStamp))
		if err != nil {
			return err
		}
		confValids = append(confValids, vf)
	}

	txf, err := formatAPITransaction(tx)
	if err != nil {
		return err
	}

	mvf, err := formatAPICoordinatorStamp(tx.CoordinatorStamp().(coordinatorStamp))
	if err != nil {
		return err
	}

	req := &api.StoreTransactionRequest{
		MinedTransaction: &api.MinedTransaction{
			CoordinatorStamp: mvf,
			CrossValidations: confValids,
			Transaction:      txf,
		},
		Timestamp: time.Now().Unix(),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}
	sig, err := lastPv.Sign(reqBytes)
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	req.SignatureRequest = sig

	ackChan := make(chan bool)
	var nbStorageAck int32
	var nbStorageFailures int32
	var failed bool

	for _, p := range pool.Nodes() {
		go func(p electedNode) {

			n, err := pr.nodeReader.FindByPublicKey(p.PublicKey().(publicKey))
			if err != nil {
				ackChan <- false
				return
			}

			serverAddr := fmt.Sprintf("%s:%d", n.IP().String(), n.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("STORE TRANSACTION REQUEST - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.StoreTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				pr.Logger.Error("STORE TRANSACTION RESPONSE - ERROR: " + grpcErr.Message())
				ackChan <- false
				return
			}

			pr.Logger.Debug("STORE TRANSACTION RESPONSE - " + time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.StoreTransactionResponse{
				Timestamp: res.Timestamp,
			})
			if err != nil {
				pr.Logger.Error("STORE TRANSACTION RESPONSE - ERROR: " + err.Error())
				ackChan <- false
				return
			}
			if ok, err := lastPub.Verify(resBytes, res.SignatureResponse); !ok || err != nil {
				pr.Logger.Error("STORE TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			ackChan <- true
		}(p.(electedNode))
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbStorageAck, 1)
		} else {
			atomic.AddInt32(&nbStorageFailures, 1)
		}
		if atomic.LoadInt32(&nbStorageFailures) == int32(len(pool.Nodes())) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbStorageAck) == int32(minStorage) {
			break
		}
	}

	if failed {
		return errors.New("transaction storage failed")
	}

	return nil
}
