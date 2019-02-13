package rpc

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"
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
			CpuLoad:               p.AppState().CPULoad(),
			DiscoveredPeersNumber: int32(p.AppState().DiscoveredPeersNumber()),
			FreeDiskSpace:         float32(p.AppState().FreeDiskSpace()),
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
		discovery.NewPeerAppState(p.AppState.Version, discovery.PeerStatus(p.AppState.Status), float64(p.AppState.GeoPosition.Longitude), float64(p.AppState.GeoPosition.Longitude), p.AppState.CpuLoad, float64(p.AppState.FreeDiskSpace), int(p.AppState.P2PFactor), int(p.AppState.DiscoveredPeersNumber)),
	)
}

func formatTransaction(tx *api.Transaction) (chain.Transaction, error) {

	propSharedKeys, err := shared.NewEmitterKeyPair(tx.SharedKeysEmitterProposal.EncryptedPrivateKey, tx.SharedKeysEmitterProposal.PublicKey)
	if err != nil {
		return chain.Transaction{}, err
	}

	data := make(map[string]string, 0)
	for k, v := range tx.Data {
		data[k] = v
	}

	return chain.NewTransaction(tx.Address, chain.TransactionType(tx.Type), data,
		time.Unix(tx.Timestamp, 0),
		tx.PublicKey,
		propSharedKeys,
		tx.Signature,
		tx.EmitterSignature,
		tx.TransactionHash)
}

func formatMinedTransaction(t *api.Transaction, mv *api.MasterValidation, valids []*api.MinerValidation) (chain.Transaction, error) {

	masterValid, err := formatMasterValidation(mv)
	if err != nil {
		return chain.Transaction{}, err
	}

	confValids := make([]chain.MinerValidation, 0)
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

	masterValidation, err := chain.NewMasterValidation(mv.PreviousTransactionMiners, mv.ProofOfWork, preValid)
	return masterValidation, err
}

func formatAPIValidation(v chain.MinerValidation) *api.MinerValidation {
	return &api.MinerValidation{
		PublicKey: v.MinerPublicKey(),
		Signature: v.MinerSignature(),
		Status:    api.MinerValidation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

func formatValidation(v *api.MinerValidation) (chain.MinerValidation, error) {
	return chain.NewMinerValidation(chain.ValidationStatus(v.Status), time.Unix(v.Timestamp, 0), v.PublicKey, v.Signature)
}

func formatAPITransaction(tx chain.Transaction) *api.Transaction {
	return &api.Transaction{
		Address:          tx.Address(),
		Data:             tx.Data(),
		Type:             api.TransactionType(tx.TransactionType()),
		PublicKey:        tx.PublicKey(),
		Signature:        tx.Signature(),
		EmitterSignature: tx.EmitterSignature(),
		Timestamp:        tx.Timestamp().Unix(),
		TransactionHash:  tx.TransactionHash(),
		SharedKeysEmitterProposal: &api.SharedKeyPair{
			EncryptedPrivateKey: tx.EmitterSharedKeyProposal().EncryptedPrivateKey(),
			PublicKey:           tx.EmitterSharedKeyProposal().PublicKey(),
		},
	}
}

func formatAPIMasterValidation(masterValid chain.MasterValidation) *api.MasterValidation {
	return &api.MasterValidation{
		ProofOfWork:               masterValid.ProofOfWork(),
		PreviousTransactionMiners: masterValid.PreviousTransactionMiners(),
		PreValidation:             formatAPIValidation(masterValid.Validation()),
	}
}
