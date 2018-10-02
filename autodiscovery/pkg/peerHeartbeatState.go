package discovery

import (
	"fmt"
	"time"
)

//PeerHeartbeatState describes the living state of a peer
type PeerHeartbeatState interface {
	GenerationTime() time.Time
	ElapsedHeartbeats() int64
	MoreRecentThan(PeerHeartbeatState) bool
	String() string
}

type heartbeatState struct {
	generationTime    time.Time
	elapsedHeartbeats int64
}

//GenerationTime returns the peer's generation time
func (hb heartbeatState) GenerationTime() time.Time {
	return hb.generationTime
}

//ElapsedHeartbeats returns the peer's elapsed living seconds from the latest refresh
func (hb heartbeatState) ElapsedHeartbeats() int64 {
	if hb.elapsedHeartbeats == 0 {
		hb.refreshElapsedHeartbeats()
	}
	return hb.elapsedHeartbeats
}

func (hb *heartbeatState) refreshElapsedHeartbeats() {
	hb.elapsedHeartbeats = time.Now().Unix() - hb.generationTime.Unix()
}

//RecentThan check if the current heartbeat state is more recent than the another heartbeat state
func (hb heartbeatState) MoreRecentThan(hbS PeerHeartbeatState) bool {

	//more recent generation time
	if hb.generationTime.Unix() > hbS.GenerationTime().Unix() {
		return true
	} else if hb.generationTime.Unix() == hbS.GenerationTime().Unix() {
		if hb.elapsedHeartbeats == hbS.ElapsedHeartbeats() {
			return false
		} else if hb.elapsedHeartbeats > hbS.ElapsedHeartbeats() {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (h heartbeatState) String() string {
	return fmt.Sprintf("Generation time: %s, Elapsed heartbeats: %d",
		h.GenerationTime().String(),
		h.ElapsedHeartbeats())
}

//NewPeerHeartbeatState creates a new peer's heartbeat state
func NewPeerHeartbeatState(genTime time.Time, elapsedHb int64) PeerHeartbeatState {
	return heartbeatState{
		generationTime:    genTime,
		elapsedHeartbeats: elapsedHb,
	}
}
