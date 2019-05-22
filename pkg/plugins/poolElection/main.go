package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

type publicKey interface {
	Marshal() []byte
	Verify(data []byte, sig []byte) (bool, error)
}

type privateKey interface {
	Sign(data []byte) ([]byte, error)
}

//ElectedNode represents a elected node in a pool
type ElectedNode interface {
	PublicKey() interface{}
	IsUnreachable() bool
	IsCoordinator() bool
	PatchNumber() int
	IsOK() bool
	MarshalJSON() ([]byte, error)
}

type electedNode struct {
	publicKey     publicKey
	isUnreachable bool
	isCoord       bool
	patchNumber   int
	isOK          bool
}

//NewElectedNode create a new elected node
func NewElectedNode(pb interface{}, isUnreach bool, isCoord bool, patchNb int, isOK bool) (interface{}, error) {
	nPub, ok := pb.(publicKey)
	if !ok {
		return nil, errors.New("elected node: public key is not valid")
	}

	return electedNode{
		publicKey:     nPub,
		isUnreachable: isUnreach,
		isCoord:       isCoord,
		patchNumber:   patchNb,
		isOK:          isOK,
	}, nil
}

func (e electedNode) PublicKey() interface{} {
	return e.publicKey
}

func (e electedNode) IsUnreachable() bool {
	return e.isUnreachable
}

func (e electedNode) IsCoordinator() bool {
	return e.isCoord
}

func (e electedNode) PatchNumber() int {
	return e.patchNumber
}

func (e electedNode) IsOK() bool {
	return e.isOK
}

func (e electedNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"publicKey":     e.publicKey.Marshal(),
		"isUnreachable": e.isUnreachable,
		"isCoordinator": e.isCoord,
		"patchNumber":   e.patchNumber,
		"isOk":          e.isOK,
	})
}

//ElectedNodeList represents a list of elected node within signature and public key
type ElectedNodeList interface {
	Nodes() []interface{}
	CreatorPublicKey() interface{}
	CreatorSignature() []byte
}

type electedNodeList struct {
	nodes            []interface{}
	creatorPublicKey interface{}
	creatorSignature []byte

	ElectedNodeList
}

//NewElectedNodeList create elected node list with its signature and public key
func NewElectedNodeList(nodes []interface{}, pubk interface{}, sig []byte) (interface{}, error) {

	if len(nodes) == 0 {
		return nil, errors.New("elected node list: missing elected nodes")
	}

	en := make([]interface{}, 0)
	for _, n := range nodes {
		if _, ok := n.(ElectedNode); !ok {
			return nil, errors.New("elected node list: invalid node type")
		}
		en = append(en, n.(ElectedNode))
	}

	if pubk == nil {
		return nil, errors.New("elected node list: missing creator's public key")
	}
	cPubk, ok := pubk.(publicKey)
	if !ok {
		return nil, errors.New("elected node list: invalid creator's public key")
	}

	if len(sig) == 0 {
		return nil, errors.New("elected node list: missing creator's signature")
	}

	nJSON, err := json.Marshal(en)
	if err != nil {
		return nil, err
	}

	if ok, err := cPubk.Verify(nJSON, sig); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("elected node list: signature is invalid")
	}

	return electedNodeList{
		nodes:            en,
		creatorPublicKey: cPubk,
		creatorSignature: sig,
	}, nil
}

func (e electedNodeList) Nodes() []interface{} {
	return e.nodes
}

func (e electedNodeList) CreatorPublicKey() interface{} {
	return e.creatorPublicKey
}

func (e electedNodeList) CreatorSignature() []byte {
	return e.creatorSignature
}

type geoPatch interface {
	ID() int
}

type node interface {
	PublicKey() interface{}
	Patch() geoPatch
	IsReachable() bool
	Status() int
}

type nodeReader interface {
	Reachables() ([]node, error)
	CountReachables() (int, error)
	FindByPublicKey(key []byte) (node, error)
}

//FindCoordinatorPool finds a list of coordinator nodes by using an entropy sort based on the transaction and minimum number of coordinator
func FindCoordinatorPool(txHash []byte, authorizedNodesKeys [][]byte, firstCrossPvKey []byte, nodePvKey interface{}, nodePubKey interface{}, r interface{}) (interface{}, error) {

	nPv, ok := nodePvKey.(privateKey)
	if !ok {
		return nil, errors.New("find coordinator pool: invalid node private key")
	}

	nPb, ok := nodePubKey.(publicKey)
	if !ok {
		return nil, errors.New("find coordinator pool: invalid node public key")
	}

	nReader, ok := r.(nodeReader)
	if !ok {
		return nil, errors.New("find coordinator pool: invalid node reader type")
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "entropySort/plugin.so"))
	if err != nil {
		return nil, err
	}

	entSortSym, err := p.Lookup("EntropySort")
	if err != nil {
		return nil, err
	}

	entropySort := entSortSym.(func([]byte, [][]byte, []byte) ([][]byte, error))
	sortedKeys, err := entropySort(txHash, authorizedNodesKeys, firstCrossPvKey)
	if err != nil {
		return nil, err
	}

	nbReachables, err := nReader.CountReachables()
	if err != nil {
		return nil, err
	}
	nbCoordinators := RequiredNumberOfCoordinators(len(authorizedNodesKeys), nbReachables)
	var nbReachableCoords int

	pool := make([]interface{}, 0)

	for i := 0; nbReachableCoords < nbCoordinators && i < len(sortedKeys); i++ {
		n, err := nReader.FindByPublicKey(sortedKeys[i])
		if err != nil {
			return nil, err
		}

		//check if the node exists, happens only when there is some networking issues
		//or if the node has not been discovered yet by the gossip service
		if n == nil {
			continue
		}

		//Add the node to the pool
		electNode, err := NewElectedNode(n.PublicKey(), !n.IsReachable(), true, n.Patch().ID(), n.Status() == 0)
		if err != nil {
			return nil, err
		}
		pool = append(pool, electNode.(ElectedNode))
		if n.IsReachable() {
			nbReachableCoords++
		}
	}
	if nbReachableCoords != nbCoordinators {
		return nil, fmt.Errorf("cannot proceed transaction with an invalid number of reachables coordinator nodes (%d)", nbReachableCoords)
	}

	//Sign the elected node list
	poolJSON, err := json.Marshal(pool)
	if err != nil {
		return nil, err
	}

	sig, err := nPv.Sign(poolJSON)
	if err != nil {
		return nil, err
	}

	return NewElectedNodeList(pool, nPb, sig)
}

//FindValidationPool lookups a validation pool from a transaction hash and a required number using the entropy sort
func FindValidationPool(txHash []byte, minValidations int, masterNodeKey []byte, authorizedNodesKeys [][]byte, firstCrossPvKey []byte, nodePvKey interface{}, nodePubKey interface{}, r interface{}) (interface{}, error) {

	nPv, ok := nodePvKey.(privateKey)
	if !ok {
		return nil, errors.New("find validation pool: invalid node private key")
	}

	nPb, ok := nodePubKey.(publicKey)
	if !ok {
		return nil, errors.New("find validation pool: invalid node public key")
	}

	nReader, ok := r.(nodeReader)
	if !ok {
		return nil, errors.New("find validation pool: invalid node reader type")
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "entropySort/plugin.so"))
	if err != nil {
		return nil, err
	}

	entSortSym, err := p.Lookup("EntropySort")
	if err != nil {
		return nil, err
	}

	entropySort := entSortSym.(func([]byte, [][]byte, []byte) ([][]byte, error))
	sortedKeys, err := entropySort(txHash, authorizedNodesKeys, firstCrossPvKey)
	if err != nil {
		return nil, err
	}

	p, err = plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "patch/plugin.so"))
	if err != nil {
		return nil, err
	}

	reqValidNbSym, err := p.Lookup("ValidationRequiredPatchNumber")
	if err != nil {
		return nil, err
	}

	validationRequiredPatchNumber := reqValidNbSym.(func(int, []interface{}) (int, error))

	reachabesNodes, err := nReader.Reachables()
	reachablesNodePatches := make([]interface{}, 0)
	for _, n := range reachabesNodes {
		reachablesNodePatches = append(reachablesNodePatches, n.Patch())
	}

	requiredPatchNb, err := validationRequiredPatchNumber(minValidations, reachablesNodePatches)
	if err != nil {
		return nil, err
	}

	//challenge the validations nodes by providing more nodes validations
	nbReachables, err := nReader.CountReachables()
	if err != nil {
		return nil, err
	}
	maxNbValidations := minValidations
	if nbReachables >= minValidations+(minValidations/2) {
		maxNbValidations = minValidations + (minValidations / 2)
	}

	var sortedReachables int
	var sortedPatchIDs []int
	pool := make([]interface{}, 0)

	for i := 0; (sortedReachables < maxNbValidations || len(sortedPatchIDs) < requiredPatchNb) && i < len(sortedKeys); i++ {
		n, err := nReader.FindByPublicKey(sortedKeys[i])
		if err != nil {
			return nil, err
		}

		//Add the node to the pool
		electedNode, err := NewElectedNode(n.PublicKey(), !n.IsReachable(), false, n.Patch().ID(), n.Status() == 0)
		if err != nil {
			return nil, err
		}
		pool = append(pool, electedNode)

		//Need a view of the reachable and unreachables for a better validation
		if n.IsReachable() {

			//Reference the patch of the node if it's not already insert by helping to determinate
			//the number of distinct patches retrieved for the check of the required number of patches
			var existingPatch bool
			for _, id := range sortedPatchIDs {
				if id == n.Patch().ID() {
					existingPatch = true
					break
				}
			}
			if !existingPatch {
				sortedPatchIDs = append(sortedPatchIDs, n.Patch().ID())
			}

			sortedReachables++
		}

	}

	if sortedReachables < maxNbValidations {
		return nil, errors.New("cannot proceed transaction with an invalid number of reachabled validation nodes")
	}

	if len(sortedPatchIDs) < requiredPatchNb {
		return nil, errors.New("cannot proceed transaction with missing patches validation nodes")
	}

	//Sign the elected node list
	poolJSON, err := json.Marshal(pool)
	if err != nil {
		return nil, err
	}

	sig, err := nPv.Sign(poolJSON)
	if err != nil {
		return nil, err
	}

	return NewElectedNodeList(pool, nPb, sig)
}

//FindStoragePool searches a storage pool for the given address
func FindStoragePool(address []byte, nodePv interface{}, nodePub interface{}, r interface{}) (interface{}, error) {

	nPv, ok := nodePv.(privateKey)
	if !ok {
		return nil, errors.New("find storage pool: invalid node private key")
	}

	nPb, ok := nodePub.(publicKey)
	if !ok {
		return nil, errors.New("find storage pool: invalid node public key")
	}

	nReader, ok := r.(nodeReader)
	if !ok {
		return nil, errors.New("find storage pool: invalid node reader type")
	}

	nodes, err := nReader.Reachables()
	if err != nil {
		return nil, fmt.Errorf("find storage pool: %s", err.Error())
	}

	pool := make([]interface{}, 0)

	//TODO: implement storage pool election
	for _, n := range nodes {

		//Add the node to the pool
		electedNode, err := NewElectedNode(n.PublicKey(), !n.IsReachable(), false, n.Patch().ID(), n.Status() == 0)
		if err != nil {
			return nil, err
		}
		pool = append(pool, electedNode)
	}

	//Sign the elected node list
	poolJSON, err := json.Marshal(pool)
	if err != nil {
		return nil, err
	}

	sig, err := nPv.Sign(poolJSON)
	if err != nil {
		return nil, fmt.Errorf("find storage pool: %s", err.Error())
	}

	l, err := NewElectedNodeList(pool, nPb, sig)
	if err != nil {
		return nil, fmt.Errorf("find storage pool: %s", err.Error())
	}
	return l, nil
}

//RequiredNumberOfCoordinators returns the number of coordinator based on the network capacity
func RequiredNumberOfCoordinators(nbNodes int, nbReachables int) int {
	if nbNodes < 5 && nbReachables >= 1 {
		return 1
	} else if nbNodes >= 5 && nbReachables <= 5 {
		return 1
	}
	return 5
}

//RequiredValidationNumber returns the need number of validations for a transaction either based on the network topology or the transaction fees
func RequiredValidationNumber(txType int, txFees float64, nbReachables int, authorizedNodesKeys [][]byte) (int, error) {

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return 0, err
	}

	sysTxSym, err := p.Lookup("SystemTransactionType")
	if err != nil {
		return 0, err
	}

	if txType != *sysTxSym.(*int) && len(authorizedNodesKeys) <= 3 {
		return 0, errors.New("no enough nodes in the network to validate this transaction")
	}

	if txType == *sysTxSym.(*int) {
		return requiredValidationNumberForSysTX(len(authorizedNodesKeys), nbReachables)
	}

	return requiredValidationNumberWithFees(txFees, nbReachables)
}

//requiredValidationNumberForSysTX returns the number of validations needed for a validation based on the network topology
func requiredValidationNumberForSysTX(nbNodes int, nbReachableNodes int) (int, error) {
	if nbNodes <= 2 && nbReachableNodes == 1 {
		return 1, nil
	}
	if nbNodes <= 5 && nbReachableNodes >= 1 {
		return nbReachableNodes, nil
	}
	if nbNodes > 5 && nbReachableNodes >= 5 {
		return 5, nil
	}
	return 0, errors.New("no enough nodes in the network to validate this transaction")
}

//requiredValidationNumberWithFees returns the number of validations needed for a validation based on the transaction fees
func requiredValidationNumberWithFees(txFees float64, nbReachablesNodes int) (validationNumber int, err error) {

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "fee/plugin.so"))
	if err != nil {
		return 0, err
	}

	feeMatrixSym, err := p.Lookup("FeeMatrix")
	if err != nil {
		return 0, err
	}

	feeMatrix := feeMatrixSym.(func() []float64)
	fees := feeMatrix()

	//3,5,7,9,11,13,15,17,19,21,23,.....
	validationsRange := make([]int, 0)
	for i := 3; i <= 100; i += 2 {
		validationsRange = append(validationsRange, i)
	}

	for i := range validationsRange {
		if txFees <= fees[i] {
			validationNumber = validationsRange[i]
			if validationNumber > nbReachablesNodes {
				validationNumber = nbReachablesNodes
			}
			break
		}
	}

	return
}
