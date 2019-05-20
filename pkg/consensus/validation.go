package consensus

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/uniris/uniris-core/pkg/logging"

// 	"github.com/uniris/uniris-core/pkg/shared"

// 	"github.com/uniris/uniris-core/pkg/chain"
// 	"github.com/uniris/uniris-core/pkg/crypto"
// )

// //LeadMining lead the mining workflow
// //
// //The workflow includes:
// // - TimeLock the transaction
// // - Pre-validate (master validation)
// // - Executes the proof of work
// // - Requests validation confirmations
// // - Requests storage
// func LeadMining(tx chain.Transaction, nbValidations int, wHeaders chain.WelcomeNodeHeader, poolR PoolRequester, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, sharedKeyReader shared.KeyReader, nodeReader NodeReader, l logging.Logger) error {
// 	l.Info("transaction " + string(tx.TransactionHash()) + "is in progress")

// 	if !tx.Address().IsValid() {
// 		return errors.New("invalid transaction address")
// 	}

// 	lastVPool, err := findLastValidationPool(tx.Address(), tx.TransactionType(), poolR, nodeReader)
// 	if err != nil {
// 		return err
// 	}

// 	sPool, err := FindStoragePool(tx.Address(), nodeReader)
// 	if err != nil {
// 		return err
// 	}

// 	if err := poolR.RequestTransactionTimeLock(sPool, tx.TransactionHash(), tx.Address(), nodePub); err != nil {
// 		return fmt.Errorf("transaction lock failed: %s", err.Error())
// 	}

// 	l.Info("transaction " + string(tx.TransactionHash()) + "is locked")

// 	go func() {
// 		vPool, err := FindValidationPool(tx.Address(), nbValidations, nodePub, nodeReader, sharedKeyReader)
// 		if err != nil {
// 			l.Error("transaction find validation pool failed: " + err.Error())
// 			return
// 		}

// 		masterValid, err := preValidateTransaction(tx, wHeaders, sPool, vPool, lastVPool, nodePub, nodePv, sharedKeyReader, nodeReader)
// 		if err != nil {
// 			l.Error("transaction pre-validation failed: " + err.Error())
// 			return
// 		}
// 		confirmValids, err := requestValidations(tx, masterValid, vPool, nbValidations, poolR)
// 		if err != nil {
// 			l.Error("transaction validation confirmations failed: " + err.Error())
// 		}
// 		if err := tx.Mined(masterValid, confirmValids); err != nil {
// 			l.Error("transaction mining is invalid: " + err.Error())
// 		}
// 		l.Info("transaction " + string(tx.TransactionHash()) + "is validated")
// 		if err := storeTransaction(tx, sPool, poolR, l); err != nil {
// 			l.Error("transaction storage failed: " + err.Error())
// 			return
// 		}
// 	}()

// 	return nil
// }

// func storeTransaction(tx chain.Transaction, sPool Pool, poolR PoolRequester, l logging.Logger) error {
// 	minReplicas := GetMinimumReplicas(tx.TransactionHash().Digest())
// 	if err := poolR.RequestTransactionStorage(sPool, minReplicas, tx); err != nil {
// 		return fmt.Errorf("transaction storage failed: %s", err.Error())
// 	}

// 	l.Info("transaction " + string(tx.TransactionHash()) + "is stored")
// 	return nil
// }

// //preValidateTransaction checks the incoming transaction as master node by ensure the transaction integrity and perform the proof of work. A valiation will result from this action
// func preValidateTransaction(tx chain.Transaction, wHeaders chain.WelcomeNodeHeader, sPool Pool, vPool Pool, lastVPool Pool, nodePub crypto.PublicKey, nodePv crypto.PrivateKey, sharedKeyReader shared.KeyReader, nodeReader NodeReader) (chain.MasterValidation, error) {
// 	if _, err := tx.IsValid(); err != nil {
// 		return chain.MasterValidation{}, err
// 	}

// 	pow, err := proofOfWork(tx, sharedKeyReader)
// 	if err != nil {
// 		return chain.MasterValidation{}, err
// 	}
// 	validStatus := chain.ValidationKO
// 	if pow != nil {
// 		validStatus = chain.ValidationOK
// 	}
// 	preValid, err := buildValidation(validStatus, nodePub, nodePv)
// 	if err != nil {
// 		return chain.MasterValidation{}, err
// 	}

// 	lastsKeys := make([]crypto.PublicKey, 0)
// 	for _, pm := range lastVPool {
// 		lastsKeys = append(lastsKeys, pm.PublicKey())
// 	}

// 	vHeaders, sHeaders, err := buildHeaders(nodeReader, nodePub, vPool, sPool)
// 	if err != nil {
// 		return chain.MasterValidation{}, err
// 	}
// 	masterValid, err := chain.NewMasterValidation(lastsKeys, pow, preValid, wHeaders, vHeaders, sHeaders)
// 	if err != nil {
// 		return chain.MasterValidation{}, err
// 	}

// 	return masterValid, nil
// }

// func buildHeaders(r NodeReader, nodePubK crypto.PublicKey, vPool Pool, sPool Pool) (vHeaders []chain.NodeHeader, sHeaders []chain.NodeHeader, err error) {

// 	//Add the master node in the header of the validation node list (for traceability and reward)
// 	masterNode, err := r.FindByPublicKey(nodePubK)
// 	if err != nil {
// 		return
// 	}
// 	vHeaders = append(vHeaders, chain.NewNodeHeader(nodePubK, false, true, masterNode.patch.patchid, masterNode.status == NodeOK))

// 	//Fill the validation headers with the elected validation pool
// 	for _, n := range vPool {
// 		vHeaders = append(vHeaders, chain.NewNodeHeader(n.PublicKey(), !n.isReachable, false, n.patch.patchid, n.status == NodeOK))
// 	}

// 	//Fill the storage headers with the elected storage pool
// 	for _, n := range sPool {
// 		sHeaders = append(sHeaders, chain.NewNodeHeader(n.PublicKey(), !n.isReachable, false, n.patch.patchid, n.status == NodeOK))
// 	}

// 	return
// }

// func proofOfWork(tx chain.Transaction, sharedKeyReader shared.KeyReader) (pow crypto.PublicKey, err error) {
// 	emKeys, err := sharedKeyReader.EmitterCrossKeypairs()
// 	if err != nil {
// 		return
// 	}

// 	txBytes, err := tx.MarshalBeforeEmitterSignature()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, kp := range emKeys {
// 		if ok := kp.PublicKey().Verify(txBytes, tx.EmitterSignature()); ok {
// 			return kp.PublicKey(), nil
// 		}
// 	}

// 	return nil, nil
// }

// func findLastValidationPool(txAddr crypto.VersionnedHash, txType chain.TransactionType, pr PoolRequester, r NodeReader) (prevP Pool, err error) {

// 	sPool, err := FindStoragePool(txAddr, r)
// 	if err != nil {
// 		return
// 	}

// 	tx, err := pr.RequestLastTransaction(sPool, txAddr, txType)
// 	if err != nil {
// 		return
// 	}
// 	if tx == nil {
// 		return
// 	}

// 	for _, key := range tx.MasterValidation().PreviousValidationNodes() {
// 		node, err := r.FindByPublicKey(key)
// 		if err != nil {
// 			return Pool{}, err
// 		}
// 		prevP = append(prevP, node)
// 	}

// 	return
// }

// func requestValidations(tx chain.Transaction, masterValid chain.MasterValidation, vPool Pool, minValids int, poolR PoolRequester) ([]chain.Validation, error) {
// 	validations, err := poolR.RequestTransactionValidations(vPool, tx, minValids, masterValid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !IsValidationConsensusReach(validations) {
// 		return nil, errors.New("invalid transaction")
// 	}

// 	return validations, nil
// }

// //RequiredValidationNumber returns the need number of validations for a transaction either based on the network topology or the transaction fees
// func RequiredValidationNumber(txType chain.TransactionType, txFees float64, nodeReader NodeReader, keyReader shared.KeyReader) (int, error) {

// 	nodeKeys, err := keyReader.AuthorizedNodesPublicKeys()
// 	if err != nil {
// 		return 0, err
// 	}

// 	nbReachables, err := nodeReader.CountReachables()
// 	if err != nil {
// 		return 0, nil
// 	}

// 	if txType != chain.SystemTransactionType && len(nodeKeys) <= 3 {
// 		return 0, errors.New("no enough nodes in the network to validate this transaction")
// 	}

// 	if txType == chain.SystemTransactionType {
// 		return requiredValidationNumberForSysTX(len(nodeKeys), nbReachables)
// 	}

// 	return requiredValidationNumberWithFees(txFees, nbReachables), nil
// }

// //requiredValidationNumberForSysTX returns the number of validations needed for a validation based on the network topology
// func requiredValidationNumberForSysTX(nbNodes int, nbReachableNodes int) (int, error) {
// 	if nbNodes <= 2 && nbReachableNodes == 1 {
// 		return 1, nil
// 	}
// 	if nbNodes <= 5 && nbReachableNodes >= 1 {
// 		return nbReachableNodes, nil
// 	}
// 	if nbNodes > 5 && nbReachableNodes >= 5 {
// 		return 5, nil
// 	}
// 	return 0, errors.New("no enough nodes in the network to validate this transaction")
// }

// //requiredValidationNumberWithFees returns the number of validations needed for a validation based on the transaction fees
// func requiredValidationNumberWithFees(txFees float64, nbReachablesNodes int) (validationNumber int) {
// 	fees := feesMatrix()

// 	//3,5,7,9,11,13,15,17,19,21,23,.....
// 	validationsRange := make([]int, 0)
// 	for i := 3; i <= 100; i += 2 {
// 		validationsRange = append(validationsRange, i)
// 	}

// 	for i := range validationsRange {
// 		if txFees <= fees[i] {
// 			validationNumber = validationsRange[i]
// 			if validationNumber > nbReachablesNodes {
// 				validationNumber = nbReachablesNodes
// 			}
// 			break
// 		}
// 	}

// 	return
// }

// //TransactionFees compute the fees earned for a given transaction
// func TransactionFees(txType chain.TransactionType, txData map[string][]byte) float64 {
// 	if txType == chain.SystemTransactionType {
// 		return 0
// 	}
// 	//TODO: compute fees base on the data sent and type of the transaction
// 	return 0.001
// }

// //feesMatrix returns the fees matrix
// func feesMatrix() (fees []float64) {
// 	//0.001,0.01,0.1,1,10,100,1000,10000,0.1M,1M
// 	for i := 0.001; i < 1000000; i *= 10 {
// 		fees = append(fees, i)
// 	}
// 	return
// }

// // IsValidationConsensusReach determinates if for the node validations the consensus is reached
// func IsValidationConsensusReach(valids []chain.Validation) bool {
// 	//TODO: maybe to improve
// 	for _, v := range valids {
// 		if v.Status() == chain.ValidationKO {
// 			return false
// 		}
// 	}
// 	return true
// }
