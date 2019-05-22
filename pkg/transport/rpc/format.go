package rpc

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"plugin"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/uniris/uniris-core/pkg/discovery"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
)

func formatPeerDigest(p *api.PeerDigest) discovery.Peer {
	return discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP(p.Identity.Ip), int(p.Identity.Port), p.Identity.PublicKey),
		discovery.NewPeerHeartbeatState(time.Unix(p.HeartbeatState.GenerationTime, 0), p.HeartbeatState.ElapsedHeartbeats),
	)
}

func formatPeerIdentity(p *api.PeerIdentity) discovery.PeerIdentity {
	return discovery.NewPeerIdentity(net.ParseIP(p.Ip), int(p.Port), p.PublicKey)
}

func formatPeerDigestAPI(p discovery.Peer) *api.PeerDigest {
	return &api.PeerDigest{
		Identity: &api.PeerIdentity{
			Ip:        p.Identity().IP().String(),
			Port:      int32(p.Identity().Port()),
			PublicKey: p.Identity().PublicKey(),
		},
		HeartbeatState: &api.PeerHeartbeatState{
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
			GenerationTime:    p.HeartbeatState().GenerationTime().Unix(),
		},
	}
}

func formatPeerIdentityAPI(p discovery.Peer) *api.PeerIdentity {
	return &api.PeerIdentity{
		Ip:        p.Identity().IP().String(),
		Port:      int32(p.Identity().Port()),
		PublicKey: p.Identity().PublicKey(),
	}
}

func formatPeerDiscoveredAPI(p discovery.Peer) *api.PeerDiscovered {
	return &api.PeerDiscovered{
		Identity: &api.PeerIdentity{
			Ip:        p.Identity().IP().String(),
			Port:      int32(p.Identity().Port()),
			PublicKey: p.Identity().PublicKey(),
		},
		HeartbeatState: &api.PeerHeartbeatState{
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
			GenerationTime:    p.HeartbeatState().GenerationTime().Unix(),
		},
		AppState: &api.PeerAppState{
			CpuLoad:              p.AppState().CPULoad(),
			ReachablePeersNumber: int32(p.AppState().ReachablePeersNumber()),
			FreeDiskSpace:        float32(p.AppState().FreeDiskSpace()),
			GeoPosition: &api.PeerAppState_GeoCoordinates{
				Latitude:  float32(p.AppState().GeoPosition().Latitude()),
				Longitude: float32(p.AppState().GeoPosition().Longitude()),
			},
			P2PFactor: int32(p.AppState().P2PFactor()),
			Status:    api.PeerAppState_PeerStatus(p.AppState().Status()),
			Version:   p.AppState().Version(),
		},
	}
}

func formatPeerDiscovered(p *api.PeerDiscovered) discovery.Peer {
	return discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP(p.Identity.Ip), int(p.Identity.Port), p.Identity.PublicKey),
		discovery.NewPeerHeartbeatState(time.Unix(p.HeartbeatState.GenerationTime, 0), p.HeartbeatState.ElapsedHeartbeats),
		discovery.NewPeerAppState(p.AppState.Version, discovery.PeerStatus(p.AppState.Status), float64(p.AppState.GeoPosition.Longitude), float64(p.AppState.GeoPosition.Longitude), p.AppState.CpuLoad, float64(p.AppState.FreeDiskSpace), int(p.AppState.P2PFactor), int(p.AppState.ReachablePeersNumber)),
	)
}

func formatTransaction(tx *api.Transaction) (transaction, error) {

	data := make(map[string]interface{}, 0)
	for k, v := range tx.Data {
		data[k] = v
	}

	p, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return nil, err
	}

	sym, err := p.Lookup("NewTransaction")
	if err != nil {
		return nil, err
	}

	pc, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}

	pcSym, err := pc.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePubF := pcSym.(func(k []byte) (interface{}, error))
	txpub, err := parsePubF(tx.PublicKey)
	if err != nil {
		return nil, err
	}

	f := sym.(func(addr []byte, txType int, data map[string]interface{}, timestamp time.Time, pubK interface{}, sig []byte, originSig []byte, coordS interface{}, crossV []interface{}) (interface{}, error))

	t, err := f(tx.Address, int(tx.Type), data, time.Unix(tx.Timestamp, 0), txpub, tx.Signature, tx.OriginSignature, nil, nil)
	if err != nil {
		return nil, err
	}
	return t.(transaction), nil
}

func formatMinedTransaction(t *api.Transaction, mv *api.CoordinatorStamp, valids []*api.ValidationStamp) (transaction, error) {

	coordS, err := formatCoordinatorStamp(mv)
	if err != nil {
		return nil, err
	}

	crossV := make([]interface{}, 0)
	for _, v := range valids {
		txValid, err := formatValidationStamp(v)
		if err != nil {
			return nil, err
		}
		crossV = append(crossV, txValid)
	}

	pTx, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "transaction/plugin.so"))
	if err != nil {
		return nil, err
	}

	pTxSym, err := pTx.Lookup("NewTransaction")
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{}, 0)
	for k, v := range t.Data {
		data[k] = v
	}

	newTF := pTxSym.(func(addr []byte, txType int, data map[string]interface{}, timestamp time.Time, pubK interface{}, sig []byte, originSig []byte, coordS interface{}, crossV []interface{}) (interface{}, error))

	pK, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}

	pKSym, err := pK.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}

	parsePub := pKSym.(func([]byte) (interface{}, error))
	pubK, err := parsePub(t.PublicKey)
	if err != nil {
		return nil, err
	}

	tx, err := newTF(t.Address, int(t.Type), data, time.Unix(t.Timestamp, 0), pubK, t.Signature, t.OriginSignature, coordS, crossV)
	if err != nil {
		return nil, err
	}
	return tx.(transaction), nil
}

func formatCoordinatorStamp(cs *api.CoordinatorStamp) (coordinatorStamp, error) {
	preValid, err := formatValidationStamp(cs.ValidationStamp)
	if err != nil {
		return nil, err
	}

	coordN, err := formatElectedNodeList(cs.ElectedCoordinatorNodes)
	if err != nil {
		return nil, err
	}
	crossV, err := formatElectedNodeList(cs.ElectedCrossValidationNodes)
	if err != nil {
		return nil, err
	}

	storeN, err := formatElectedNodeList(cs.ElectedStorageNodes)
	if err != nil {
		return nil, err
	}

	pc, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}

	pcSym, err := pc.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePubF := pcSym.(func(k []byte) (interface{}, error))

	pow, err := parsePubF(cs.ProofOfWork)
	if err != nil {
		return nil, err
	}

	previousNodeKeys := make([][]byte, 0)
	for _, prevNodeKey := range cs.PreviousCrossValidators {
		previousNodeKeys = append(previousNodeKeys, prevNodeKey)
	}

	pC, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "coordinatorStamp/plugin.so"))
	if err != nil {
		return nil, err
	}
	pCSym, err := pC.Lookup("NewCoordinatorStamp")
	if err != nil {
		return nil, err
	}

	f := pCSym.(func(prevCrossV [][]byte, pow interface{}, validStamp interface{}, txHash []byte, elecCoordNodes interface{}, elecCrossVNodes interface{}, elecStorNodes interface{}) (interface{}, error))

	coordS, err := f(previousNodeKeys, pow, preValid, []byte("hash"), coordN, crossV, storeN)
	if err != nil {
		return nil, err
	}
	return coordS.(coordinatorStamp), nil
}

func formatAPIValidation(v validationStamp) (*api.ValidationStamp, error) {

	return &api.ValidationStamp{
		PublicKey: v.NodePublicKey().(publicKey).Marshal(),
		Signature: v.NodeSignature(),
		Status:    api.ValidationStamp_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}, nil
}

func formatValidationStamp(v *api.ValidationStamp) (validationStamp, error) {

	pc, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}

	pcSym, err := pc.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePubF := pcSym.(func(k []byte) (interface{}, error))

	nodeKey, err := parsePubF(v.PublicKey)
	if err != nil {
		return nil, err
	}

	pV, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "validationStamp/plugin.so"))
	if err != nil {
		return nil, err
	}

	pvSym, err := pV.Lookup("NewValidationStamp")
	if err != nil {
		return nil, err
	}

	f := pvSym.(func(status int, t time.Time, nodePubk interface{}, nodeSig []byte) (interface{}, error))

	vStamp, err := f(int(v.Status), time.Unix(v.Timestamp, 0), nodeKey, v.Signature)
	if err != nil {
		return nil, err
	}
	return vStamp.(validationStamp), nil
}

func formatAPITransaction(tx transaction) (*api.Transaction, error) {

	data := make(map[string]*any.Any, len(tx.Data()))
	for k, v := range tx.Data() {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data[k] = &any.Any{Value: b}
	}

	return &api.Transaction{
		Address:         tx.Address(),
		Timestamp:       tx.Timestamp().Unix(),
		Type:            api.TransactionType(tx.Type()),
		Data:            data,
		PublicKey:       tx.PreviousPublicKey().(publicKey).Marshal(),
		Signature:       tx.Signature(),
		OriginSignature: tx.OriginSignature(),
	}, nil
}

func formatAPICoordinatorStamp(coordStamp coordinatorStamp) (*api.CoordinatorStamp, error) {
	powKey := coordStamp.ProofOfWork().(publicKey).Marshal()
	prevNodeKeys := make([][]byte, 0)
	for _, k := range coordStamp.PreviousCrossValidators() {
		prevNodeKeys = append(prevNodeKeys, k)
	}

	v, err := formatAPIValidation(coordStamp.ValidationStamp().(validationStamp))
	if err != nil {
		return nil, err
	}

	coordNodes := formatElectedNodeListAPI(coordStamp.ElectedCoordinatorNodes().(electedNodeList))
	crossVNodes := formatElectedNodeListAPI(coordStamp.ElectedCrossValidationNodes().(electedNodeList))
	storNodes := formatElectedNodeListAPI(coordStamp.ElectedStorageNodes().(electedNodeList))

	return &api.CoordinatorStamp{
		ProofOfWork:                 powKey,
		PreviousCrossValidators:     prevNodeKeys,
		ValidationStamp:             v,
		TransactionHash:             coordStamp.TransactionHash(),
		ElectedCoordinatorNodes:     coordNodes,
		ElectedCrossValidationNodes: crossVNodes,
		ElectedStorageNodes:         storNodes,
	}, nil
}

func formatElectedNodeList(l *api.ElectedNodeList) (electedNodeList, error) {

	nodes := make([]interface{}, len(l.Nodes))
	for i, n := range l.Nodes {
		el, err := formatElectedNode(n)
		if err != nil {
			return nil, err
		}
		nodes[i] = el
	}

	pKey, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}
	pParsePubSym, err := pKey.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePub := pParsePubSym.(func([]byte) (interface{}, error))

	pubk, err := parsePub(l.PublicKey)
	if err != nil {
		return nil, err
	}

	pElec, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "poolElection/plugin.so"))
	if err != nil {
		return nil, err
	}
	pNewESym, err := pElec.Lookup("NewElectedNodeList")
	if err != nil {
		return nil, err
	}

	f := pNewESym.(func(nodes []interface{}, pubk interface{}, sig []byte) (interface{}, error))

	list, err := f(nodes, pubk, l.Signature)
	if err != nil {
		return nil, err
	}
	return list.(electedNodeList), nil
}

func formatElectedNodeListAPI(l electedNodeList) *api.ElectedNodeList {
	nodes := make([]*api.ElectedNode, len(l.Nodes()))
	for i, n := range l.Nodes() {
		nodes[i] = formatElectedNodeAPI(n.(electedNode))
	}

	return &api.ElectedNodeList{
		Nodes:     nodes,
		PublicKey: l.CreatorPublicKey().(publicKey).Marshal(),
		Signature: l.CreatorSignature(),
	}
}

func formatElectedNode(n *api.ElectedNode) (electedNode, error) {

	pElec, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "poolElection/plugin.so"))
	if err != nil {
		return nil, err
	}
	pNewESym, err := pElec.Lookup("NewElectedNode")
	if err != nil {
		return nil, err
	}
	newElectNodeF := pNewESym.(func(pb interface{}, isUnreach bool, isCoord bool, patchNb int, isOK bool) (interface{}, error))

	pKey, err := plugin.Open(filepath.Join(os.Getenv("PLUGINS_DIR"), "key/plugin.so"))
	if err != nil {
		return nil, err
	}
	pParsePubSym, err := pKey.Lookup("ParsePublicKey")
	if err != nil {
		return nil, err
	}
	parsePub := pParsePubSym.(func([]byte) (interface{}, error))

	pubk, err := parsePub(n.PublicKey)
	if err != nil {
		return nil, err
	}

	el, err := newElectNodeF(pubk, n.IsUnreachable, n.IsMaster, int(n.PatchNumber), n.IsOK)
	if err != nil {
		return nil, err
	}
	return el.(electedNode), nil
}

func formatElectedNodeAPI(n electedNode) *api.ElectedNode {
	return &api.ElectedNode{
		IsMaster:    n.IsCoordinator(),
		IsOK:        n.IsOK(),
		PatchNumber: int32(n.PatchNumber()),
		PublicKey:   n.PublicKey().(publicKey).Marshal(),
	}
}
