package discovery

import (
	"errors"
	"fmt"
	"log"
	"net"
)

//ErrNTPShift is returned when the NTP clock drift to much
var ErrNTPShift = errors.New("system Clock have a big Offset check the ntp configuration of the system")

//ErrNTPFailure is returned when the NTP server cannot be reached
var ErrNTPFailure = errors.New("could not get reply from ntp servers")

//ErrGRPCServer is returned when the GRPC servers cannot be reached
var ErrGRPCServer = errors.New("GRPC servers are not running")

//ErrGeoPosition is returned when the Geoposition cannot be reached
var ErrGeoPosition = errors.New("geographic position cannot be found")

//BootstrapingMinTime is the necessary minimum time on seconds to finish learning about the network
const BootstrapingMinTime = 1800

//NetworkChecker is the interface that provides methods to get the peer network monrmation
type NetworkChecker interface {
	CheckNtpState() error
	CheckInternetState() error
	CheckGRPCServer() error
}

//SystemReader retrieve local system information
type SystemReader interface {

	//GeoPosition retrieves the peer's geographic position
	GeoPosition() (lon float64, lat float64, err error)

	//CPULoad retrieves the load on the peer's CPU
	CPULoad() (string, error)

	//FreeDiskSpace retrieves the available free disk space of the peer
	FreeDiskSpace() (float64, error)

	IP() (net.IP, error)
}

func updateSelf(self Peer, reachables []Peer, seeds []PeerIdentity, db dbWriter, netCheck NetworkChecker, sysR SystemReader) (Peer, error) {
	status, err := localStatus(self, seedReachableAverage(seeds, reachables), netCheck)
	if err != nil {
		return self, err
	}

	_, _, _, cpu, space, err := systemInfo(sysR)
	if err != nil {
		if err == ErrGeoPosition {
			status = FaultyPeer
			log.Println(ErrGeoPosition)
		} else {
			return self, err
		}
	}

	self.SelfRefresh(status, space, cpu, p2pFactor(reachables), len(reachables))
	return self, nil
}

//systemInfo retrieves system information such geo position, IP, CPU load and free disk space
func systemInfo(sr SystemReader) (lon float64, lat float64, ip net.IP, cpu string, space float64, err error) {
	lon, lat, err = sr.GeoPosition()
	if err != nil {
		err = ErrGeoPosition
		return
	}

	ip, err = sr.IP()
	if err != nil {
		return
	}

	cpu, err = sr.CPULoad()
	if err != nil {
		return
	}

	space, err = sr.FreeDiskSpace()
	if err != nil {
		return
	}

	return
}

//localStatus retrieves the status of the local peer.
func localStatus(p Peer, seedAvgDiscovery int, nv NetworkChecker) (PeerStatus, error) {
	if err := nv.CheckInternetState(); err != nil {
		fmt.Printf("networking error: %s\n", err.Error())
		return FaultyPeer, nil
	}

	if err := nv.CheckNtpState(); err != nil {
		if err == ErrNTPShift || err == ErrNTPFailure {
			fmt.Printf("networking error: %s\n", err.Error())
			return StorageOnlyPeer, nil
		}
		return FaultyPeer, err
	}

	if err := nv.CheckGRPCServer(); err != nil {
		if err == ErrGRPCServer {
			fmt.Printf("networking error: %s\n", err.Error())
			return FaultyPeer, nil
		}
		return FaultyPeer, err
	}

	if seedAvgDiscovery == 0 {
		return BootstrapingPeer, nil
	}

	if t := p.HeartbeatState().ElapsedHeartbeats(); t < BootstrapingMinTime && seedAvgDiscovery > p.AppState().ReachablePeersNumber() {
		return BootstrapingPeer, nil
	} else if t < BootstrapingMinTime && seedAvgDiscovery <= p.AppState().ReachablePeersNumber() {
		return OkPeerStatus, nil
	} else {
		return OkPeerStatus, nil
	}
}

func p2pFactor(peers []Peer) int {
	return 1
}

//seedReachableAverage computes an avergage of the reachables peers retrieved by the seeds
func seedReachableAverage(seeds []PeerIdentity, reachablePeers []Peer) int {
	avg := 0
	for i := 0; i < len(seeds); i++ {
		ipseed := seeds[i].IP()

		var foundPeer *Peer
		for _, p := range reachablePeers {
			if p.Identity().IP().Equal(ipseed) {
				foundPeer = &p
				break
			}
		}
		if foundPeer == nil {
			continue
		}
		avg += foundPeer.AppState().ReachablePeersNumber()
	}
	avg = avg / len(seeds)
	return avg
}
