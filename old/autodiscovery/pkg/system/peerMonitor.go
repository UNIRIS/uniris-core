package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

type peerMonitor struct {
}

//CPULoad retrieves the load on the peer's CPU
func (m peerMonitor) CPULoad() (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("cat", "/proc/loadavg")
	} else {
		return "", errors.New("You platform is not supported")
	}

	outload, err := cmd.CombinedOutput()
	if err != nil {
		return "--", err
	}

	res := fmt.Sprintf("%d - %s", runtime.NumCPU(), string(outload))
	return res, nil
}

//FreeDiskSpace retrieves the available free disk (k bytes) space of the peer
func (m peerMonitor) FreeDiskSpace() (float64, error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0.0, err
	}
	syscall.Statfs(wd, &stat)
	return float64((stat.Bavail * uint64(stat.Bsize)) / 1024), nil
}

//P2PFactor request the update P2PFactor from the AI Daemon
func (m peerMonitor) P2PFactor() (int, error) {
	return 0, nil
}

//NewPeerMonitor creates an instance which implements monitoring.Watcher
func NewPeerMonitor() monitoring.PeerMonitor {
	return peerMonitor{}
}
