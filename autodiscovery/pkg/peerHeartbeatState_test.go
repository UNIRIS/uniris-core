package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Refreshes elapsed hearbeats
	Given an heartbeat state
	When we want to refresh the elapsed heartbeats
	Then we get the new elapsed heartbeats based on the current time
*/
func TestRefreshElapsedHeartbeats(t *testing.T) {
	hb := heartbeatState{
		generationTime: time.Now(),
	}
	time.Sleep(2 * time.Second)
	hb.refreshElapsedHeartbeats()
	assert.Equal(t, int64(2), hb.ElapsedHeartbeats())
}

/*
Scenario: Gets the elapsed heartbeats when no previous refresh
	Given a fresh heartbeat state
	When we get the elaspsed hearbeats
	Then we refresh the elapsed hearbeats and returns it
*/
func TestGetElapsedBeatsWithoutPrevRefresh(t *testing.T) {
	hb := heartbeatState{
		generationTime: time.Now(),
	}
	time.Sleep(2 * time.Second)
	assert.Equal(t, int64(2), hb.ElapsedHeartbeats())
}

/*
Scenario: Checks if an heartbeat state is more recent based on the upper generation time
	Given an heartbeat state with a generation time set as (now + 2 seconds)
	When we compare with another state with generation time set as now
	Then the first heartbeat is more recent
*/
func TestMoreRecentUpperGenTime(t *testing.T) {
	hb := heartbeatState{generationTime: time.Now().Add(2 * time.Second)}
	hb2 := heartbeatState{generationTime: time.Now()}
	assert.True(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on the same generation time
	Given a heartbeat state with a generation time as now
	When we compare with another with the same generation time
	Then the first heartbeat is not more recent
*/
func TestMoreRecentSameGenTimeSameElapsedBeats(t *testing.T) {
	hb := heartbeatState{generationTime: time.Now()}
	hb2 := heartbeatState{generationTime: time.Now()}
	assert.False(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on  the same generation and upper elapsed beats
	Given an heartbeat state with a generation time set as now and 500 elapsed beats
	When we compare with another with the same generation time and 300 elapsed beats
	Then the first heartbeat is more recent
*/
func TestMoreRecentSameGenTimeUpperElapsedBeats(t *testing.T) {
	hb := heartbeatState{generationTime: time.Now(), elapsedHeartbeats: 500}
	hb2 := heartbeatState{generationTime: time.Now(), elapsedHeartbeats: 300}
	assert.True(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on the same generation time and lower elapsed beats
	Given an heartbeat state with a generation time set as now and 300 elapsed beats
	When we compare with another with the same generation time and 500 elapsed beats
	Then the first heartbeat is not more recent
*/
func TestMoreRecentSameGenTimeLowerElapsedBeats(t *testing.T) {
	hb := heartbeatState{generationTime: time.Now(), elapsedHeartbeats: 300}
	hb2 := heartbeatState{generationTime: time.Now(), elapsedHeartbeats: 500}
	assert.False(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on a lower generation time
	Given an hearbeat with a generation time set as now
	When we compare with another with generation set as now + 2 seconds
	Then the first heartbeat is not more recent
*/
func TestMoreRecentLowerGenTime(t *testing.T) {
	hb := heartbeatState{generationTime: time.Now()}
	hb2 := heartbeatState{generationTime: time.Now().Add(2 * time.Second)}
	assert.False(t, hb.MoreRecentThan(hb2))
}
