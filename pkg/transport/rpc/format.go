package rpc

import (
	"errors"
	"net"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"

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

func formatTransaction(tx *api.Transaction) (chain.Transaction, error) {

	propPubKey, err := crypto.ParsePublicKey(tx.SharedKeysEmitterProposal.PublicKey)
	if err != nil {
		return chain.Transaction{}, nil
	}
	propSharedKeys, err := shared.NewEmitterCrossKeyPair(tx.SharedKeysEmitterProposal.EncryptedPrivateKey, propPubKey)
	if err != nil {
		return chain.Transaction{}, err
	}

	data := make(map[string][]byte, 0)
	for k, v := range tx.Data {
		data[k] = v
	}

	txPubKey, err := crypto.ParsePublicKey(tx.PublicKey)
	if err != nil {
		return chain.Transaction{}, errors.New("invalid public key")
	}

	return chain.NewTransaction(tx.Address, chain.TransactionType(tx.Type), data,
		time.Unix(tx.Timestamp, 0),
		txPubKey,
		propSharedKeys,
		tx.Signature,
		tx.EmitterSignature,
		tx.TransactionHash)
}

func formatMinedTransaction(t *api.Transaction, mv *api.MasterValidation, valids []*api.Validation) (chain.Transaction, error) {

	masterValid, err := formatMasterValidation(mv)
	if err != nil {
		return chain.Transaction{}, err
	}

	confValids := make([]chain.Validation, 0)
	for _, v := range valids {
		txValid, err := formatValidation(v)
		if err != nil {
			return chain.Transaction{}, err
		}
		confValids = append(confValids, txValid)
	}

	txRoot, err := formatTransaction(t)
	if err != nil {
		return chain.Transaction{}, err
	}
	tx, err := chain.NewTransaction(txRoot.Address(), txRoot.TransactionType(), txRoot.Data(), txRoot.Timestamp(), txRoot.PublicKey(), txRoot.EmitterSharedKeyProposal(), txRoot.Signature(), txRoot.EmitterSignature(), txRoot.TransactionHash())
	if err != nil {
		return chain.Transaction{}, err
	}

	if err := tx.Mined(masterValid, confValids); err != nil {
		return chain.Transaction{}, err
	}

	return tx, nil
}

func formatMasterValidation(mv *api.MasterValidation) (chain.MasterValidation, error) {
	preValid, err := formatValidation(mv.PreValidation)
	if err != nil {
		return chain.MasterValidation{}, err
	}

	wHeaders, err := formatWelcomeNodeHeaders(mv.WelcomeHeaders)
	if err != nil {
		return chain.MasterValidation{}, err
	}
	vHeaders, err := formatNodeHeaders(mv.ValidationHeaders)
	if err != nil {
		return chain.MasterValidation{}, err
	}
	sHeaders, err := formatNodeHeaders(mv.StorageHeaders)
	if err != nil {
		return chain.MasterValidation{}, err
	}
	powKey, err := crypto.ParsePublicKey(mv.ProofOfWork)
	if err != nil {
		return chain.MasterValidation{}, errors.New("invalid proof of work public key")
	}

	previousNodeKeys := make([]crypto.PublicKey, 0)
	for _, prevNodeKey := range mv.PreviousValidationNodes {
		nodePubKey, err := crypto.ParsePublicKey(prevNodeKey)
		if err != nil {
			return chain.MasterValidation{}, errors.New("invalid previous transaction node public key")
		}
		previousNodeKeys = append(previousNodeKeys, nodePubKey)
	}

	masterValidation, err := chain.NewMasterValidation(previousNodeKeys, powKey, preValid, wHeaders, vHeaders, sHeaders)
	return masterValidation, err
}

func formatAPIValidation(v chain.Validation) (*api.Validation, error) {

	nodeKey, err := v.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	return &api.Validation{
		PublicKey: nodeKey,
		Signature: v.Signature(),
		Status:    api.Validation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}, nil
}

func formatValidation(v *api.Validation) (chain.Validation, error) {
	nodeKey, err := crypto.ParsePublicKey(v.PublicKey)
	if err != nil {
		return chain.Validation{}, errors.New("validation public key is invalid")
	}

	return chain.NewValidation(chain.ValidationStatus(v.Status), time.Unix(v.Timestamp, 0), nodeKey, v.Signature)
}

func formatAPITransaction(tx chain.Transaction) (*api.Transaction, error) {

	txPub, err := tx.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	propPub, err := tx.EmitterSharedKeyProposal().PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	return &api.Transaction{
		Address:          tx.Address(),
		Data:             tx.Data(),
		Type:             api.TransactionType(tx.TransactionType()),
		PublicKey:        txPub,
		Signature:        tx.Signature(),
		EmitterSignature: tx.EmitterSignature(),
		Timestamp:        tx.Timestamp().Unix(),
		TransactionHash:  tx.TransactionHash(),
		SharedKeysEmitterProposal: &api.SharedKeyPair{
			EncryptedPrivateKey: tx.EmitterSharedKeyProposal().EncryptedPrivateKey(),
			PublicKey:           propPub,
		},
	}, nil
}

func formatAPIMasterValidation(masterValid chain.MasterValidation) (*api.MasterValidation, error) {

	powKey, err := masterValid.ProofOfWork().Marshal()
	if err != nil {
		return nil, err
	}

	prevNodeKeys := make([][]byte, 0)
	for _, k := range masterValid.PreviousValidationNodes() {
		nodeKey, err := k.Marshal()
		if err != nil {
			return nil, err
		}
		prevNodeKeys = append(prevNodeKeys, nodeKey)
	}

	v, err := formatAPIValidation(masterValid.Validation())
	if err != nil {
		return nil, err
	}

	wHeaders, err := formatWelcomeNodeHeadersAPI(masterValid.WelcomeHeaders())
	if err != nil {
		return nil, err
	}
	vHeaders, err := formatNodeHeadersAPI(masterValid.ValidationHeaders())
	if err != nil {
		return nil, err
	}
	sHeaders, err := formatNodeHeadersAPI(masterValid.StorageHeaders())
	if err != nil {
		return nil, err
	}

	return &api.MasterValidation{
		ProofOfWork:             powKey,
		PreviousValidationNodes: prevNodeKeys,
		PreValidation:           v,
		WelcomeHeaders:          wHeaders,
		ValidationHeaders:       vHeaders,
		StorageHeaders:          sHeaders,
	}, nil
}

func formatNodeHeadersAPI(headers []chain.NodeHeader) (apiHeaders []*api.NodeHeader, err error) {
	for _, h := range headers {
		pubKey, err := h.PublicKey().Marshal()
		if err != nil {
			return nil, err
		}
		apiHeaders = append(apiHeaders, &api.NodeHeader{
			IsMaster:      h.IsMaster(),
			IsUnreachable: h.IsUnreachable(),
			PublicKey:     pubKey,
			PatchNumber:   int32(h.PatchNumber()),
			IsOK:          h.IsOk(),
		})
	}
	return
}

func formatWelcomeNodeHeadersAPI(wheaders chain.WelcomeNodeHeader) (*api.WelcomeNodeHeader, error) {

	masterlist := make([]*api.NodeHeader, 0)

	wnpubk, err := wheaders.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	masterlist, err = formatNodeHeadersAPI(wheaders.NodeHeaders())
	if err != nil {
		return nil, err
	}

	return &api.WelcomeNodeHeader{
		PublicKey:   wnpubk,
		MastersList: masterlist,
		Signature:   wheaders.Sig(),
	}, nil
}

func formatNodeHeaders(apiHeaders []*api.NodeHeader) (headers []chain.NodeHeader, err error) {
	for _, h := range apiHeaders {
		pubKey, err := crypto.ParsePublicKey(h.PublicKey)
		if err != nil {
			return nil, err
		}
		headers = append(headers, chain.NewNodeHeader(
			pubKey,
			h.IsUnreachable,
			h.IsMaster,
			int(h.PatchNumber),
			h.IsOK,
		))
	}
	return
}

func formatWelcomeNodeHeaders(apiwHeaders *api.WelcomeNodeHeader) (chain.WelcomeNodeHeader, error) {

	masterlist := make([]chain.NodeHeader, 0)
	wnh := chain.WelcomeNodeHeader{}

	wpubk, err := crypto.ParsePublicKey(apiwHeaders.PublicKey)
	if err != nil {
		return wnh, err
	}

	masterlist, err = formatNodeHeaders(apiwHeaders.MastersList)
	if err != nil {
		return wnh, err
	}

	return chain.NewWelcomeNodeHeader(wpubk, masterlist, apiwHeaders.Signature), nil
}
