package entities

import "time"

//PeerHeartbeat represents how fresh is the peer information
type PeerHeartbeat struct {
	GenerationTime time.Time
	ElapsedBeats   int64
}
