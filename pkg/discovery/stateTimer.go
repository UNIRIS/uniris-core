package discovery

import (
	"time"
)

type timerState struct {
	enabled bool
	start   time.Time
	ticker  *time.Ticker
}

var ts timerState

func isEnabledtimerState() bool {
	return ts.enabled
}

func disabletimerState() {
	ts.enabled = false
	ts.ticker.Stop()
}

func newtimerState() error {

	tic := time.NewTicker(1 * time.Second)

	ts.ticker = tic
	ts.start = time.Now()
	ts.enabled = true

	endtime := time.Now().Local().Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(0) + time.Second*time.Duration(BootstrapingMinTime))

	//Remove the timelock when the countdown is reached
	go func() {
		for range tic.C {
			if time.Now().Unix() == endtime.Unix() {
				disabletimerState()
			}
		}
	}()

	return nil

}
