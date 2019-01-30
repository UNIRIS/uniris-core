package rpc

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/shared"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/transaction"
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

func formatTransaction(tx *api.Transaction) (transaction.Transaction, error) {

	sharedKeys, err := shared.NewKeyPair(tx.Proposal.SharedEmitterKeys.EncryptedPrivateKey, tx.Proposal.SharedEmitterKeys.PublicKey)
	if err != nil {
		return transaction.Transaction{}, err
	}

	prop, err := transaction.NewProposal(sharedKeys)
	if err != nil {
		return transaction.Transaction{}, err
	}

	data := make(map[string]string, 0)
	for k, v := range tx.Data {
		data[k] = v
	}

	return transaction.New(tx.Address, transaction.Type(tx.Type), data,
		time.Unix(tx.Timestamp, 0),
		tx.PublicKey,
		tx.Signature,
		tx.EmitterSignature,
		prop,
		tx.TransactionHash)
}

func formatMinedTransaction(t *api.Transaction, mv *api.MasterValidation, valids []*api.MinerValidation) (transaction.Transaction, error) {

	masterValid, err := formatMasterValidation(mv)
	if err != nil {
		return transaction.Transaction{}, err
	}

	confValids := make([]transaction.MinerValidation, 0)
	for _, v := range valids {
		txValid, err := formatValidation(v)
		if err != nil {
			return transaction.Transaction{}, err
		}
		confValids = append(confValids, txValid)
	}

	txRoot, err := formatTransaction(t)
	if err != nil {
		return transaction.Transaction{}, err
	}
	tx, err := transaction.New(txRoot.Address(), txRoot.Type(), txRoot.Data(), txRoot.Timestamp(), txRoot.PublicKey(), txRoot.Signature(), txRoot.EmitterSignature(), txRoot.Proposal(), txRoot.TransactionHash())
	if err != nil {
		return transaction.Transaction{}, err
	}

	if err := tx.AddMining(masterValid, confValids); err != nil {
		return transaction.Transaction{}, err
	}

	return tx, nil
}

func formatMasterValidation(mv *api.MasterValidation) (transaction.MasterValidation, error) {
	prevMiners := make([]transaction.PoolMember, 0)
	for _, m := range mv.PreviousTransactionMiners {
		prevMiners = append(prevMiners, transaction.NewPoolMember(net.ParseIP(m.Ip), int(m.Port)))
	}

	preValid, err := formatValidation(mv.PreValidation)
	if err != nil {
		return transaction.MasterValidation{}, err
	}

	masterValidation, err := transaction.NewMasterValidation(prevMiners, mv.ProofOfWork, preValid)
	return masterValidation, err
}

func formatAPIValidation(v transaction.MinerValidation) *api.MinerValidation {
	return &api.MinerValidation{
		PublicKey: v.MinerPublicKey(),
		Signature: v.MinerSignature(),
		Status:    api.MinerValidation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

func formatValidation(v *api.MinerValidation) (transaction.MinerValidation, error) {
	return transaction.NewMinerValidation(transaction.ValidationStatus(v.Status), time.Unix(v.Timestamp, 0), v.PublicKey, v.Signature)
}

func formatAPITransaction(tx transaction.Transaction) *api.Transaction {
	return &api.Transaction{
		Address:          tx.Address(),
		Data:             tx.Data(),
		Type:             api.TransactionType(tx.Type()),
		PublicKey:        tx.PublicKey(),
		Signature:        tx.Signature(),
		EmitterSignature: tx.EmitterSignature(),
		Timestamp:        tx.Timestamp().Unix(),
		TransactionHash:  tx.TransactionHash(),
		Proposal: &api.TransactionProposal{
			SharedEmitterKeys: &api.SharedKeys{
				EncryptedPrivateKey: tx.Proposal().SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           tx.Proposal().SharedEmitterKeyPair().PublicKey(),
			},
		},
	}
}

func formatAPIMasterValidation(masterValid transaction.MasterValidation) *api.MasterValidation {

	prevMiners := make([]*api.PoolMember, 0)
	for _, m := range masterValid.PreviousTransactionMiners() {
		prevMiners = append(prevMiners, &api.PoolMember{
			Ip:   m.IP().String(),
			Port: int32(m.Port()),
		})
	}

	return &api.MasterValidation{
		ProofOfWork:               masterValid.ProofOfWork(),
		PreviousTransactionMiners: prevMiners,
		PreValidation:             formatAPIValidation(masterValid.Validation()),
	}
}
