package discovery

import "time"

//PeerHeartbeatState describes the living state of a peer
type PeerHeartbeatState interface {
	GenerationTime() time.Time
	ElapsedHeartbeats() int64
}

type heartbeatState struct {
	generationTime    time.Time
	elapsedHeartbeats int64
}

func (hb heartbeatState) GenerationTime() time.Time {
	return hb.generationTime
}

func (hb heartbeatState) ElapsedHeartbeats() int64 {
	return hb.elapsedHeartbeats
}

//NewPeerHeartbeatState creates a new peer's heartbeat state
func NewPeerHeartbeatState(genTime time.Time, elapsedHb int64) PeerHeartbeatState {
	return heartbeatState{
		generationTime:    genTime,
		elapsedHeartbeats: elapsedHb,
	}
}
