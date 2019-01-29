package discovery

import (
	"errors"
	"net"
)

var ErrNTPShift = errors.New("System Clock have a big Offset check the ntp configuration of the system")
var ErrNTPFailure = errors.New("Could not get reply from ntp servers")

//BootstrapingMinTime is the necessary minimum time on seconds to finish learning about the network
const BootstrapingMinTime = 1800

//PeerNetworker is the interface that provides methods to get the peer network information
type PeerNetworker interface {
	CheckNtpState() error
	CheckInternetState() error
}

type PeerInformer interface {

	//GeoPosition retrieves the peer's geographic position
	GeoPosition() (lon float64, lat float64, err error)

	//CPULoad retrieves the load on the peer's CPU
	CPULoad() (string, error)

	//FreeDiskSpace retrieves the available free disk space of the peer
	FreeDiskSpace() (float64, error)

	IP() (net.IP, error)
}

func getPeerSystemInfo(info PeerInformer) (lon float64, lat float64, ip net.IP, cpu string, space float64, err error) {
	lon, lat, err = info.GeoPosition()
	if err != nil {
		return
	}

	ip, err = info.IP()
	if err != nil {
		return
	}

	cpu, err = info.CPULoad()
	if err != nil {
		return
	}

	space, err = info.FreeDiskSpace()
	if err != nil {
		return
	}

	return
}

func getPeerStatus(p Peer, seedAvgDiscovery int, pn PeerNetworker) (PeerStatus, error) {
	if err := pn.CheckInternetState(); err != nil {
		return FaultyPeer, err
	}

	if err := pn.CheckNtpState(); err != nil {
		if err == ErrNTPShift || err == ErrNTPFailure {
			return StorageOnlyPeer, nil
		}
		return FaultyPeer, err
	}

	if seedAvgDiscovery == 0 {
		return BootstrapingPeer, nil
	}

	if t := p.HeartbeatState().ElapsedHeartbeats(); t < BootstrapingMinTime && seedAvgDiscovery > p.AppState().DiscoveredPeersNumber() {
		return BootstrapingPeer, nil
	} else if t < BootstrapingMinTime && seedAvgDiscovery <= p.AppState().DiscoveredPeersNumber() {
		return OkPeerStatus, nil
	} else {
		return OkPeerStatus, nil
	}
}

func getP2PFactor(peers []Peer) int {
	return 1
}

func getSeedDiscoveryAverage(seeds []Seed, knownPeers []Peer) int {
	avg := 0
	for i := 0; i < len(seeds); i++ {
		ipseed := seeds[i].IP()

		var foundPeer *Peer
		for _, p := range knownPeers {
			if p.Identity().IP().Equal(ipseed) {
				foundPeer = &p
				break
			}
		}
		if foundPeer == nil {
			continue
		}
		avg += foundPeer.AppState().DiscoveredPeersNumber()
	}
	avg = avg / len(seeds)
	return avg
}
