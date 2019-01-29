package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/uniris/uniris-core/pkg/discovery"
)

type peerInfo struct {
	privateNetwork bool
	privateIface   string
}

func NewPeerInformer(privateNetwork bool, privateIface string) discovery.PeerInformer {
	return peerInfo{privateNetwork, privateIface}
}

type position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (i peerInfo) GeoPosition() (lon float64, lat float64, err error) {
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var pos position
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&pos); err != nil {
		return
	}
	return pos.Longitude, pos.Latitude, nil
}

func (i peerInfo) CPULoad() (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("cat", "/proc/loadavg")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("sysctl", "-n", "vm.loadavg")
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

func (i peerInfo) FreeDiskSpace() (float64, error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0.0, err
	}
	syscall.Statfs(wd, &stat)
	return float64((stat.Bavail * uint64(stat.Bsize)) / 1024), nil
}

func (i peerInfo) IP() (net.IP, error) {
	if i.privateNetwork {
		return privateIP(i.privateIface)
	}
	ip, err := ipifyIP()
	if err != nil {
		ip, err := myExternalIP()
		if err != nil {
			return nil, errors.New("Cannot get the peer IP. IP providers may failed")
		}
		return ip, nil
	}
	return ip, nil
}
