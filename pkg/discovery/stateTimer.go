package discovery

import (
	"time"
)

var timerState time.Time

func refreshtimerState() {
	timerState = time.Now()
}

func resettimerState() {
	layout := "0000-00-00 15:04:05 -0700 MST"
	timerState, _ = time.Parse(layout, "0000-01-01 00:00:00 +0000 UTC")
}

// checktimerState return true if the timerstate is expired
func checktimerState() bool {
	if timerState.IsZero() {
		return true
	}
	t := time.Now()
	endtime := timerState.Add(time.Hour*time.Duration(0) + time.Minute*time.Duration(0) + time.Second*time.Duration(BootstrapingMinTime))
	return endtime.Before(t)
}
