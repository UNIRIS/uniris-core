package rpc

import (
	"net"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	uniris "github.com/uniris/uniris-core/pkg"
)

func formatPeerDigest(p *api.PeerDigest) uniris.Peer {
	return uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP(p.Identity.Ip), int(p.Identity.Port), p.Identity.PublicKey),
		uniris.NewPeerHeartbeatState(time.Unix(p.HeartbeatState.GenerationTime, 0), p.HeartbeatState.ElapsedHeartbeats),
	)
}

func formatPeerDigestAPI(p uniris.Peer) *api.PeerDigest {
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

func formatPeerDiscoveredAPI(p uniris.Peer) *api.PeerDiscovered {
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

func formatPeerDiscovered(p *api.PeerDiscovered) uniris.Peer {
	return uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP(p.Identity.Ip), int(p.Identity.Port), p.Identity.PublicKey),
		uniris.NewPeerHeartbeatState(time.Unix(p.HeartbeatState.GenerationTime, 0), p.HeartbeatState.ElapsedHeartbeats),
		uniris.NewPeerAppState(p.AppState.Version, uniris.PeerStatus(p.AppState.Status), float64(p.AppState.GeoPosition.Longitude), float64(p.AppState.GeoPosition.Longitude), p.AppState.CpuLoad, float64(p.AppState.FreeDiskSpace), int(p.AppState.P2PFactor), int(p.AppState.DiscoveredPeersNumber)),
	)
}

func formatTransaction(tx *api.Transaction) uniris.Transaction {
	prop := uniris.NewTransactionProposal(
		uniris.NewSharedKeyPair(tx.Proposal.SharedEmitterKeys.EncryptedPrivateKey, tx.Proposal.SharedEmitterKeys.PublicKey),
	)

	return uniris.NewTransactionBase(tx.Address, uniris.TransactionType(tx.Type), tx.Data,
		time.Unix(tx.Timestamp, 0),
		tx.PublicKey,
		tx.Signature,
		tx.EmitterSignature,
		prop,
		tx.TransactionHash)
}

func formatMinedTransaction(tx *api.Transaction, mv *api.MasterValidation, valids []*api.MinerValidation) uniris.Transaction {

	prevMiners := make([]uniris.PeerIdentity, 0)
	for _, m := range mv.PreviousTransactionMiners {
		prevMiners = append(prevMiners, formatPeerIdentity(m))
	}

	masterValidation := uniris.NewMasterValidation(prevMiners, mv.ProofOfWork, formatValidation(mv.PreValidation))

	confValids := make([]uniris.MinerValidation, 0)
	for _, v := range valids {
		confValids = append(confValids, formatValidation(v))
	}

	return uniris.NewMinedTransaction(formatTransaction(tx), masterValidation, confValids)
}

func formatAPIValidation(v uniris.MinerValidation) *api.MinerValidation {
	return &api.MinerValidation{
		PublicKey: v.MinerPublicKey(),
		Signature: v.MinerSignature(),
		Status:    api.MinerValidation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

func formatValidation(v *api.MinerValidation) uniris.MinerValidation {
	return uniris.NewMinerValidation(uniris.ValidationStatus(v.Status), time.Unix(v.Timestamp, 0), v.PublicKey, v.Signature)
}

func formatAPITransaction(tx uniris.Transaction) *api.Transaction {
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

func formatAPIMasterValidationAPI(masterValid uniris.MasterValidation) *api.MasterValidation {

	prevMiners := make([]*api.PeerIdentity, 0)
	for _, m := range masterValid.PreviousTransactionMiners() {
		prevMiners = append(prevMiners, formatPeerIdentityAPI(m))
	}

	return &api.MasterValidation{
		ProofOfWork:               masterValid.ProofOfWork(),
		PreviousTransactionMiners: prevMiners,
		PreValidation: &api.MinerValidation{
			PublicKey: masterValid.Validation().MinerPublicKey(),
			Signature: masterValid.Validation().MinerSignature(),
			Status:    api.MinerValidation_ValidationStatus(masterValid.Validation().Status()),
			Timestamp: masterValid.Validation().Timestamp().Unix(),
		},
	}
}

func formatPeerIdentity(identity *api.PeerIdentity) uniris.PeerIdentity {
	return uniris.NewPeerIdentity(net.ParseIP(identity.Ip), int(identity.Port), identity.PublicKey)
}

func formatPeerIdentityAPI(identity uniris.PeerIdentity) *api.PeerIdentity {
	return &api.PeerIdentity{
		Ip:        identity.IP().String(),
		Port:      int32(identity.Port()),
		PublicKey: identity.PublicKey(),
	}
}
