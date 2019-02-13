package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/uniris/uniris-core/pkg/discovery"
)

type sysReader struct {
	privateNetwork bool
	privateIface   string
}

//NewReader creates a new system reader
func NewReader(privateNetwork bool, privateIface string) discovery.SystemReader {
	return sysReader{privateNetwork, privateIface}
}

type position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (i sysReader) GeoPosition() (lon float64, lat float64, err error) {
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

func (i sysReader) CPULoad() (string, error) {
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

func (i sysReader) FreeDiskSpace() (float64, error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0.0, err
	}
	syscall.Statfs(wd, &stat)
	return float64((stat.Bavail * uint64(stat.Bsize)) / 1024), nil
}

func (i sysReader) IP() (net.IP, error) {
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

func ipifyIP() (net.IP, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(bytes)), nil
}

func myExternalIP() (net.IP, error) {
	resp, err := http.Get("http://www.myexternalip.com/raw")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(string(bytes)), nil
}

func privateIP(localIface string) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var iface net.Interface
	for _, i := range ifaces {
		if i.Name == localIface {
			iface = i
			break
		}
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	var ip net.IP
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		return ip, nil
	}
	return nil, fmt.Errorf("Cannot find a IP address from the interface %s", localIface)
}
