package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"github.com/uniris/uniris-core/autodiscovery/core/ports"
)

//RefreshPeer recalculate peer state
func RefreshPeer(peer *domain.Peer, repo ports.PeerRepository, conf ports.ConfigurationReader, analyser ports.MetricReader) error {

	ip, err := conf.GetIP()
	if err != nil {
		return err
	}

	geoPos, err := conf.GetGeoPosition()
	if err != nil {
		return err
	}

	ver, err := conf.GetVersion()
	if err != nil {
		return err
	}

	p2pFactor, err := conf.GetP2PFactor()
	if err != nil {
		return err
	}

	status, err := analyser.GetStatus()
	if err != nil {
		return err
	}

	cpuLoad, err := analyser.GetCPULoad()
	if err != nil {
		return err
	}

	freeSpace, err := analyser.GetFreeDiskSpace()
	if err != nil {
		return err
	}

	ioWaitRate, err := analyser.GetIOWaitRate()
	if err != nil {
		return err
	}

	peer.IP = ip
	peer.State.P2PFactor = p2pFactor
	peer.State.GeoPosition = geoPos
	peer.State.Version = ver
	peer.State.CPULoad = cpuLoad
	peer.State.Status = status
	peer.State.FreeDiskSpace = freeSpace
	peer.State.IOWaitRate = ioWaitRate
	if err := repo.UpdatePeer(*peer); err != nil {
		return err
	}

	return nil
}
